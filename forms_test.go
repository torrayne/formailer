package formailer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSetAndGet(t *testing.T) {
	form := Form{
		Name:   "contact",
		Ignore: []string{"_form_name"},
	}
	email := Email{
		To:      "daniel@atwood.io",
		From:    "daniel@atwood.io",
		Subject: "New Contact Form Submission",
	}
	form.AddEmail(email)
	Add(form)

	for _, set := range DefaultConfig {
		if !cmp.Equal(set, form, cmpopts.IgnoreUnexported(Form{})) {
			t.Error("Unexpected result getting form from config")
		}
	}
}
