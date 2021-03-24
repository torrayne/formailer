package formailer

import (
	"errors"
	"fmt"
	"mime"
)

type Config map[string]*Form
type Form struct {
	Name      string
	Emails    []Email
	Redirect  string
	ReCAPTCHA bool

	ignore map[string]bool
}

// DefaultConfig is the default Config
var DefaultConfig = make(Config)

func New(id string) *Form {
	f := &Form{Name: id, ignore: make(map[string]bool)}
	f.Ignore("_form_name", "g-recaptcha-response")
	Add(f)
	return f
}

// Add adds forms to the default config
func Add(forms ...*Form) {
	DefaultConfig.Add(forms...)
}

// Parse creates a submission using the default config
func Parse(contentType, body string) (*Submission, error) {
	return DefaultConfig.Parse(contentType, body)
}

// Add adds forms to the config
func (c Config) Add(forms ...*Form) {
	for _, form := range forms {
		c[form.Name] = form
	}
}

// Parse creates a submission
func (c Config) Parse(contentType string, body string) (*Submission, error) {
	submission := new(Submission)
	submission.Values = make(map[string]interface{})

	contentType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	switch contentType {
	case "application/json":
		err = submission.parseJSON(body)
	case "application/x-www-form-urlencoded":
		err = submission.parseURLEncoded(body)
	case "multipart/form-data":
		err = submission.parseMultipartForm(params["boundary"], body)
	default:
		err = errors.New("invalid content type")
	}

	form, ok := submission.Values["_form_name"].(string)
	if !ok {
		return nil, fmt.Errorf("field _form_name not of type string or not set")
	}

	submission.Form, ok = c[form]
	if !ok {
		return nil, fmt.Errorf("missing emails for form %s", form)
	}

	submission.removeIgnored()

	return submission, err
}

// AddEmail adds emails to the form
func (f *Form) AddEmail(emails ...Email) {
	f.Emails = append(f.Emails, emails...)
}

func (f *Form) Ignore(fields ...string) {
	for _, field := range fields {
		f.ignore[field] = true
	}
}
