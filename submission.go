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

// Submission is the unmarshaled version on the form submission.
// It contains the submitted values and the form settings needed for sending emails.
type Submission struct {
	// Form is the form this submission submitted as.
	Form *Form

	// Order is a list of the fields in the original order of submisson.
	// Maps in go are sorted alphabetically which causes readability issues in the generated emails. Note only first level elements are ordered.
	// The fields set using Form.Ignore will be removed from this list.
	Order []string

	// Values contains the submitted form data.
	Values map[string]interface{}

	// Attachments is a list of files to be attached to the email
	Attachments []Attachment
}

// Attachment contains file data for an email attachment
type Attachment struct {
	Filename string
	MimeType string
	Data     []byte
}

var forceStringFields = []string{
	"_form_name", "g-recaptcha-response",
}

func (s *Submission) forceString(vals url.Values) {
	for _, key := range forceStringFields {
		s.Values[key] = vals.Get(key)
	}
}

func (s *Submission) removeIgnored() {
	for i := 0; i < len(s.Order); i++ {
		if s.Form.ignore[s.Order[i]] {
			s.Order = append(s.Order[:i], s.Order[i+1:]...)
			i--
		}
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
		s.Order = append(s.Order, key)
		esc, _ := json.Marshal(key)
		index[key] = bytes.Index(b, append(esc, ':'))
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
		s.Order = append(s.Order, key)
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

		if _, ok := s.Values[key]; !ok {
			s.Order = append(s.Order, key)
		}

		s.Values[key] = values[key]
	}

	s.forceString(values)

	return nil
}

// Send sends all the emails for this form
func (s *Submission) Send() error {
	for _, e := range s.Form.Emails {
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
