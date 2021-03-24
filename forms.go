package formailer

import (
	"errors"
	"fmt"
	"mime"
	"strings"
)

// Config is a map of Forms used when parsing a submission to load the correct form settings and emails.
type Config map[string]*Form

// Form is for settings that should be set per submission but not per email. Such as redirects and ReCAPTCHA.
type Form struct {
	// ID is a case-insensitive string used to look up forms while parsing a submission. It will be matched to a submission's _form_name field.
	// If ID is not set, a case-insensitive version of Name will be used for matching instead.
	ID string

	// Name is a way to store a "Pretty" version of the form ID.
	Name string

	// Emails is a list of emails. Generally you want to use the AddEmail method instead of adding emails directly.
	Emails []Email

	// Redirect is used when with the default handlers to return 303 See Other and points the browser to the set value.
	Redirect string

	// When ReCAPTCHA is set to true the default handlers with verify the g-recaptcha-response field.
	ReCAPTCHA bool

	ignore map[string]bool
}

// DefaultConfig is the config used when using functions New, Add, and Parse.
// This helps keep boilerplate code to a minimum.
var DefaultConfig = make(Config)

// New creates a new Form and adds it to the default config.
// It also automatically sets the name to the ID and adds ignores the form name and recaptcha fields.
func New(id string) *Form {
	f := &Form{ID: id, ignore: make(map[string]bool)}
	f.Ignore("_form_name", "g-recaptcha-response")
	Add(f)
	return f
}

// Add adds forms to the default config. It allows you to add forms without the default settings provided by New.
func Add(forms ...*Form) {
	DefaultConfig.Add(forms...)
}

// Parse creates a submission using the default config.
func Parse(contentType, body string) (*Submission, error) {
	return DefaultConfig.Parse(contentType, body)
}

// Add adds forms to the config falling back on Name if ID is not set.
func (c Config) Add(forms ...*Form) {
	for _, form := range forms {
		id := strings.ToLower(or(form.ID, form.Name))
		c[id] = form
	}
}

// Parse creates a submission parsing the data based on the Content-Type header.
// Setting Submission.Form based on the _form_name field and removing any ignored fields from Submisson.Order.
func (c Config) Parse(contentType string, body string) (*Submission, error) {
	submission := new(Submission)
	submission.Values = make(map[string]interface{})

	contentType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content-type: %w", err)
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
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}

	form, ok := submission.Values["_form_name"].(string)
	if !ok || len(form) < 1 {
		return nil, errors.New("missing _form_name field in submitted form data")
	}

	form = strings.ToLower(form)
	submission.Form, ok = c[form]
	if !ok {
		return nil, fmt.Errorf("missing form config for form %s", form)
	}

	submission.removeIgnored()

	return submission, err
}

// AddEmail adds emails to the form.
func (f *Form) AddEmail(emails ...Email) {
	f.Emails = append(f.Emails, emails...)
}

// Ignore updates the Form.ignore map
func (f *Form) Ignore(fields ...string) {
	for _, field := range fields {
		f.ignore[field] = true
	}
}
