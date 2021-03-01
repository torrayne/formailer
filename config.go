package formailer

import (
	"errors"
	"fmt"
	"strings"
)

// Config is a map of form configs
type Config map[string]*Form

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

	formName, ok := submission.Values["_form_name"].(string)
	if !ok {
		return nil, fmt.Errorf("field _form_name not of type string or not set")
	}

	submission.Form = c.Get(formName)
	return submission, err
}
