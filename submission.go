package formailer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"sort"
	"strings"
)

// Submission is parsed from the body
type Submission struct {
	Emails      []Email
	Order       []string
	Values      map[string]interface{}
	Attachments []Attachment
}

// Attachment is an array of files to be attached to the email
type Attachment struct {
	Filename string
	MimeType string
	Data     []byte
}

var forceStringFields = []string{
	"_form_name", "_redirect",
	"g-recaptcha-response",
}

var ignoreFields = map[string]bool{
	"g-recaptcha-response": true,
}

func (s *Submission) forceString(vals url.Values) {
	for _, key := range forceStringFields {
		s.Values[key] = vals.Get(key)
	}
}

func (s *Submission) parseJSON(body string) error {
	b := []byte(body)
	err := json.Unmarshal(b, &s.Values)
	if err != nil {
		return err
	}

	index := make(map[string]int)
	for key := range s.Values {
		if !ignoreFields[key] {
			s.Order = append(s.Order, key)
			esc, _ := json.Marshal(key)
			index[key] = bytes.Index(b, append(esc, ':'))
		}
	}

	sort.Slice(s.Order, func(i, j int) bool {
		return index[s.Order[i]] < index[s.Order[j]]
	})

	return nil
}

func (s *Submission) parseURLEncoded(body string) error {
	vals, err := url.ParseQuery(body)
	if err != nil {
		return err
	}

	index := make(map[string]int)
	for key := range vals {
		s.Values[key] = vals[key]
		index[key] = strings.Index(body, key+"=")

		if !ignoreFields[key] {
			s.Order = append(s.Order, key)
		}
	}

	sort.Slice(s.Order, func(i, j int) bool {
		return index[s.Order[i]] < index[s.Order[j]]
	})

	s.forceString(vals)

	return nil
}

func (s *Submission) parseMultipartForm(boundary, body string) error {
	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		decodedBody = []byte(body)
	}
	decodedBody = append(decodedBody, '\n')

	values := make(url.Values)
	reader := multipart.NewReader(bytes.NewReader(decodedBody), boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		key := part.FormName()
		value := new(bytes.Buffer)
		value.ReadFrom(part)

		filename := part.FileName()
		if len(filename) > 0 {
			attachment := Attachment{
				Filename: part.FileName(),
				MimeType: part.Header.Get("Content-Type"),
				Data:     value.Bytes(),
			}
			s.Attachments = append(s.Attachments, attachment)
			values[key] = append(values[key], part.FileName())
		} else {
			values[key] = append(values[key], value.String())
		}

		if _, ok := s.Values[key]; !ok && !ignoreFields[key] {
			s.Order = append(s.Order, key)
		}

		s.Values[key] = values[key]
	}

	s.forceString(values)

	return nil
}

// Send sends all the emails for this form
func (s *Submission) Send() error {
	for _, e := range s.Emails {
		email, err := e.Email(s)
		if err != nil {
			return err
		}

		if err := e.Send(email); err != nil {
			return err
		}
	}

	return nil
}
