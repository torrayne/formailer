package formailer

//go:generate go run generate/main.go

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aymerick/douceur/inliner"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Config is a map of form configs
type Config map[string]*Form

// Form contains all the setting to send an email
type Form struct {
	Name     string
	To       string
	From     string
	Subject  string
	Redirect string
	Template string
}

// Submission is parsed from the body
type Submission struct {
	Form        *Form
	Values      map[string]string
	Attachments []Attachment
}

// Attachment is an array of files to be attached to the email
type Attachment struct {
	Filename string
	MimeType string
	Data     []byte
}

type smtpAuth struct {
	host, port, user, pass string
}

// Set is a case insentive way to set form configs
func (c *Config) Set(forms ...*Form) {
	for _, form := range forms {
		(*c)[strings.ToLower(form.Name)] = form
	}
}

// Get is a case insentive way to get form configs
func (c *Config) Get(name string) *Form {
	return (*c)[strings.ToLower(name)]
}

// Parse parses the body string based on the provided content type
func (c *Config) Parse(contentType string, body string) (*Submission, error) {
	submission := new(Submission)
	submission.Values = make(map[string]string)

	var err error
	if strings.Contains(contentType, "application/json") {
		err = submission.parseJSON(body)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err = submission.parseURLEncoded(body)
	} else if strings.Contains(contentType, "multipart/form-data") {
		err = submission.parseMultipartForm(contentType, body)
	} else {
		err = errors.New("invalid content type")
	}

	submission.Form = c.Get(submission.Values["_form_name"])

	return submission, err
}

// GetTemplate allows for fallback on the default template when no template has beeen provided
func (f *Form) GetTemplate() string {
	if len(f.Template) < 1 {
		return defaultTemplate
	}
	return f.Template
}

// SMTPServer returns a sever using the ENV for auth falling back on the default for each missing param
func (f *Form) SMTPServer() (*mail.SMTPServer, error) {
	prefix := fmt.Sprintf("SMTP_%s_", strings.ToUpper(f.Name))

	def := defaultSMTP()
	host := or(os.Getenv(prefix+"HOST"), def.host)
	port := or(os.Getenv(prefix+"PORT"), def.port)
	user := or(os.Getenv(prefix+"USER"), def.user)
	pass := or(os.Getenv(prefix+"PASS"), def.pass)

	if len(host) < 1 || len(port) < 1 || len(user) < 1 || len(pass) < 1 {
		return nil, fmt.Errorf("form %s missing SMTP configuration ", f.Name)
	}

	{
		port, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}

		server := mail.NewSMTPClient()
		server.Host = host
		server.Port = port
		server.Username = user
		server.Password = pass
		server.Encryption = mail.EncryptionTLS
		server.Authentication = mail.AuthLogin
		server.KeepAlive = false
		server.ConnectTimeout = 10 * time.Second
		server.SendTimeout = 10 * time.Second

		return server, nil
	}
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
		s.Values[k] = vals.Get(k)
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
		return err
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
	t, err := template.New("email").Parse(s.Form.GetTemplate())
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

func defaultSMTP() smtpAuth {
	return smtpAuth{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
	}
}

func or(a, b string) string {
	if len(a) < 1 {
		return b
	}
	return a
}
