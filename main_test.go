package formailer

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	contact := Form{
		Name:    "contact",
		To:      "djatwood01@gmail.com",
		From:    "daniel@atwood.io",
		Subject: "New Contact Form Submission",
	}

	cfg := make(Config)
	cfg.Set(&contact)

	res := cfg.Get(contact.Name)
	if res == nil {
		t.Errorf("Could not get config for form: %s", contact.Name)
	}
	if &contact != res {
		t.Error("Unexpected result getting form from config")
	}
}

func TestParse(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	w := multipart.NewWriter(b)
	w.SetBoundary("--myspecificboundary")
	w.WriteField("_form_name", "contact")
	w.WriteField("Name", "Daniel")
	w.WriteField("message", "This is my message")
	w.Close()

	tests := map[string]string{
		"application/json":                                   `{"_form_name": "contact", "Name":"Daniel", "message": "This is my message"}`,
		"application/x-www-form-urlencoded":                  "_form_name=contact&name=Daniel&message=This is my message",
		"multipart/form-data  boundary=--myspecificboundary": b.String(),
	}

	cfg := make(Config)
	for contentType, body := range tests {
		_, err := cfg.Parse(contentType, body)
		if err != nil && err != io.EOF {
			t.Errorf("Failed to parse data: %v", err)
		}
	}
}

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
