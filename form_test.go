package formailer

import (
	"os"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	tests := []map[string]interface{}{
		{
			"form":     &Form{Name: "testing"},
			"expected": defaultTemplate,
		},
		{
			"form":     &Form{Name: "testing", Template: "THIS IS MY TEMPLATE"},
			"expected": "THIS IS MY TEMPLATE",
		},
	}

	for _, test := range tests {
		if test["expected"] != test["form"].(*Form).GetTemplate() {
			t.Errorf("Unexpected result from Form.GetTemplate")
		}
	}
}

func TestSMTPSetup(t *testing.T) {
	os.Setenv("SMTP_HOST", "mail.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_CONTACT_USER", "username@example.com")
	os.Setenv("SMTP_CONTACT_PASS", "mysupersecretpassword")

	f := Form{Name: "Contact"}
	_, err := f.SMTPServer()
	if err != nil {
		t.Error(err)
	}
}
