package formailer

import (
	"errors"
	"fmt"
	"strings"
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
