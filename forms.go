package formailer

import (
	"errors"
	"fmt"
	"mime"
)

// Forms is a map of forms
var Forms = make(map[string]Form)

type Form struct {
	Name   string
	Emails []Email

	Ignore []string
	ignore map[string]bool
}

func sliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, value := range slice {
		m[value] = true
	}
	return m
}

// Add adds forms to the config
func Add(forms ...Form) {
	for _, form := range forms {
		if len(form.Ignore) < 1 {
			form.Ignore = []string{
				"_form_name", "_redirect",
				"g-recaptcha-response",
			}
		}

		form.ignore = sliceToMap(form.Ignore)
		Forms[form.Name] = form
	}
}

// AddEmail adds emails to the form
func (f *Form) AddEmail(emails ...Email) {
	f.Emails = append(f.Emails, emails...)
}

// Parse parses the body string based on the provided content type
func Parse(contentType string, body string) (*Submission, error) {
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

	submission.Form, ok = Forms[form]
	if !ok {
		return nil, fmt.Errorf("missing emails for form %s", form)
	}
	return submission, err
}
