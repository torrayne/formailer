package formailer

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	form := "contact"
	emails := []Email{
		{
			To:      "daniel@atwood.io",
			From:    "daniel@atwood.io",
			Subject: "New Contact Form Submission",
		},
	}

	forms := make(Forms)
	forms.Add(form, emails...)

	for i := range forms[form] {
		if forms[form][i] != emails[i] {
			t.Error("Unexpected result getting form from config")
		}
	}
}
