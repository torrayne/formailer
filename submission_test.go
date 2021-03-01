package formailer

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	form := Form{
		Name:    "Contact",
		Subject: "New Contact Form Submission",
	}

	cfg := make(Config)
	cfg.Set(&form)

	submission := Submission{
		Form: &form,
		Values: map[string]interface{}{
			"Name":       []string{"Daniel", "Atwood"},
			"Message":    "Hello, World!",
			"_form_name": "contact",
		},
	}

	_, err := submission.generate()
	if err != nil {
		t.Error(err)
	}
}
