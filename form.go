package formailer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Form contains all the setting to send an email
type Form struct {
	Name     string
	To       string
	From     string
	Subject  string
	Redirect string
	Template string
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
