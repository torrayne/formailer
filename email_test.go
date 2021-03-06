package formailer

import (
	"os"
	"testing"
)

type testGetTemplate struct {
	email    Email
	expected string
}

func TestGetTemplate(t *testing.T) {
	tests := []testGetTemplate{
		{
			email:    Email{},
			expected: defaultTemplate,
		},
		{
			email:    Email{Template: "THIS IS MY TEMPLATE"},
			expected: "THIS IS MY TEMPLATE",
		},
	}

	for _, test := range tests {
		if test.expected != test.email.template() {
			t.Errorf("Unexpected result from Form.GetTemplate")
		}
	}
}

func TestSMTPSetup(t *testing.T) {
	os.Setenv("SMTP_HOST", "mail.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_CONTACT_USER", "username@example.com")
	os.Setenv("SMTP_CONTACT_PASS", "mysupersecretpassword")

	f := Email{ID: "Contact"}
	_, err := f.server()
	if err != nil {
		t.Error(err)
	}
}

func TestGenerate(t *testing.T) {
	form := "contact"
	email := Email{
		ID:      form,
		Subject: "New Contact Form Submission",
	}

	forms := make(Forms)
	forms.Add(form, email)

	submission := Submission{
		Emails: []Email{email},
		Values: map[string]interface{}{
			"Name":       []string{"Daniel", "Atwood"},
			"Message":    "Hello, World!",
			"_form_name": "contact",
		},
	}

	_, err := email.generate(&submission)
	if err != nil {
		t.Error(err)
	}
}
