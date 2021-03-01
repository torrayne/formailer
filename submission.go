package formailer

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/url"
	"strings"
	"text/template"

	"github.com/aymerick/douceur/inliner"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Submission is parsed from the body
type Submission struct {
	Form        *Form
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
		if k == "_form_name" {
			s.Values[k] = vals.Get(k)
		} else {
			s.Values[k] = vals[k]
		}
	}
	return nil
}

func (s *Submission) parseMultipartForm(contentType, body string) error {
	header := strings.Split(contentType, ";")
	var boundary string
	for _, h := range header {
		index := strings.Index(h, "boundary")
		if index > -1 {
			boundary = strings.TrimSpace(h[index+9:])
			break
		}
	}

	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		decodedBody = []byte(body)
	}
	decodedBody = append(decodedBody, '\n')

	reader := multipart.NewReader(bytes.NewReader(decodedBody), boundary)
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

func (s *Submission) generate() (string, error) {
	t := template.New("email").Funcs(templateFuncMap)
	_, err := t.Parse(s.Form.GetTemplate())
	if err != nil {
		return "", err
	}

	var email bytes.Buffer
	err = t.Execute(&email, s)
	if err != nil {
		return "", err
	}

	return inliner.Inline(email.String())
}

// Send sends the email
func (s *Submission) Send(server *mail.SMTPServer) error {
	message, err := s.generate()
	if err != nil {
		return err
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return errors.New("failed to generate message-id")
	}

	email := mail.NewMSG()
	email.AddTo(s.Form.To)
	email.SetFrom(s.Form.From)
	email.SetSubject(s.Form.Subject)
	email.SetBody(mail.TextHTML, message)
	email.AddHeader("Message-Id", base32.StdEncoding.EncodeToString(token))

	for _, attachment := range s.Attachments {
		email.AddAttachmentData(attachment.Data, attachment.Filename, attachment.MimeType)
	}

	client, err := server.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	return email.Send(client)
}
