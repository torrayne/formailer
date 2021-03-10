package formailer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		if !cmp.Equal(forms[form][i], emails[i]) {
			t.Error("Unexpected result getting form from config")
		}
	}
}
