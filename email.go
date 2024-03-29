package formailer

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"

	// embed is used to embed the default template file
	_ "embed"

	"github.com/aymerick/douceur/inliner"
	mail "github.com/xhit/go-simple-mail/v2"
)

//go:embed template.html
var defaultTemplate string

// Email contains all the setting to send an email
type Email struct {
	// ID is used when looking up SMTP settings.
	// It is case-insensitive but will be matched as UPPERCASE. ex: SMTP_FORM-ID_HOST.
	ID string

	To      string
	From    string
	Cc      []string
	Bcc     []string
	ReplyTo string
	Subject string

	// Template is a go html template to be used when generating the email.
	Template string
}

func or(a, b string) string {
	if len(a) < 1 {
		return b
	}
	return a
}

// template allows for fallback on the default template when no template has beeen provided
func (e *Email) template() string {
	return or(e.Template, defaultTemplate)
}

// server returns a sever using the ENV for auth falling back on the default for each missing param
func (e *Email) server() (*mail.SMTPServer, error) {
	prefix := fmt.Sprintf("SMTP_%s_", strings.ToUpper(e.ID))
	host := or(os.Getenv(prefix+"HOST"), os.Getenv("SMTP_HOST"))
	user := or(os.Getenv(prefix+"USER"), os.Getenv("SMTP_USER"))
	pass := or(os.Getenv(prefix+"PASS"), os.Getenv("SMTP_PASS"))
	defaultPort := os.Getenv("SMTP_PORT")
	emailPort := os.Getenv(prefix + "PORT")
	stringPort := or(emailPort, defaultPort)

	if len(host) < 1 {
		return nil, fmt.Errorf("incomplete SMTP configuration missing %sHOST or SMTP_HOST for %s", prefix, e.ID)
	}
	if len(stringPort) < 1 {
		return nil, fmt.Errorf("incomplete SMTP configuration missing %sPORT or SMTP_PORT for %s", prefix, e.ID)
	}
	if len(user) < 1 {
		return nil, fmt.Errorf("incomplete SMTP configuration missing %sUSER or SMTP_USER for %s", prefix, e.ID)
	}
	if len(pass) < 1 {
		return nil, fmt.Errorf("incomplete SMTP configuration missing %sPASS or SMTP_PASS for %s", prefix, e.ID)
	}

	port, err := strconv.Atoi(stringPort)
	if err != nil {
		if len(emailPort) < 1 {
			prefix = "SMTP_"
		}
		return nil, fmt.Errorf("could not parse %sPORT: %w", prefix, err)
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
	server.SendTimeout = 10 * time.Minute

	return server, nil
}

func (e *Email) generate(s *Submission) (string, error) {
	t := template.New("email").Funcs(templateFuncMap)
	_, err := t.Parse(e.template())
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

// Email returns a *mail.Email generating the message with the provided submission
func (e *Email) Email(submission *Submission) (*mail.Email, error) {
	message, err := e.generate(submission)
	if err != nil {
		return nil, fmt.Errorf("failed to generate message: %w", err)
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return nil, fmt.Errorf("failed to generate message-id: %w", err)
	}

	email := mail.NewMSG()
	email.AddTo(e.To)
	email.SetFrom(e.From)
	email.SetSubject(e.Subject)
	email.SetBody(mail.TextHTML, message)
	email.AddHeader("Message-Id", base32.StdEncoding.EncodeToString(token))

	if len(e.ReplyTo) > 0 {
		email.SetReplyTo(e.ReplyTo)
	}
	for _, a := range e.Cc {
		email.AddCc(a)
	}
	for _, a := range e.Bcc {
		email.AddBcc(a)
	}
	for _, attachment := range submission.Attachments {
		email.AddAttachmentData(attachment.Data, attachment.Filename, attachment.MimeType)
	}

	return email, nil
}

// Send sends the provided email
func (e *Email) Send(email *mail.Email) error {
	server, err := e.server()
	if err != nil {
		return err
	}

	client, err := server.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	return email.Send(client)
}
