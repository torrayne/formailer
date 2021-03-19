package formailer

import (
	"errors"
	"fmt"
	"mime"
)

// Forms is a map of form configs
type Forms map[string][]Email

// Add is a case insentive way to set form configs
func (f Forms) Add(form string, emails ...Email) {
	f[form] = append(f[form], emails...)
}

// Parse parses the body string based on the provided content type
func (f Forms) Parse(contentType string, body string) (*Submission, error) {
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

	submission.Emails, ok = f[form]
	if !ok {
		return nil, fmt.Errorf("missing emails for form %s", form)
	}
	return submission, err
}
