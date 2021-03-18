package formailer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/url"
)

// Submission is parsed from the body
type Submission struct {
	Emails      []Email
	Values      map[string]interface{}
	Attachments []Attachment
}

// Attachment is an array of files to be attached to the email
type Attachment struct {
	Filename string
	MimeType string
	Data     []byte
}

func (s *Submission) parseJSON(body string) error {
	return json.Unmarshal([]byte(body), &s.Values)
}

func (s *Submission) parseURLEncoded(body string) error {
	vals, err := url.ParseQuery(body)
	if err != nil {
		return err
	}

	for k := range vals {
		switch k {
		case "_form_name", "_redirect", "g-recaptcha-response":
			s.Values[k] = vals.Get(k)
		default:
			s.Values[k] = vals[k]
		}
	}
	return nil
}

func (s *Submission) parseMultipartForm(contentType, body string) error {
	_, headerParams, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		decodedBody = []byte(body)
	}
	decodedBody = append(decodedBody, '\n')

	reader := multipart.NewReader(bytes.NewReader(decodedBody), headerParams["boundary"])
	for {
		part, err := reader.NextPart()
		if err != nil {
			return err
		}

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
			s.Values[part.FormName()] = part.FileName()
		} else {
			s.Values[part.FormName()] = value.String()
		}
	}
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
