package formailer

import (
	"io"
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
	cfg := make(Config)
	tests := map[string]string{
		"application/json":                  `{"Name":"Daniel", "message": "This is my message"}`,
		"application/x-www-form-urlencoded": "name=Daniel&message=This is my message",
		"multipart/form-data  boundary=---------------------------382742568133097519731421599912": "LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0zODI3NDI1NjgxMzMwOTc1MTk3MzE0MjE1OTk5MTINCkNvbnRlbnQtRGlzcG9zaXRpb246IGZvcm0tZGF0YTsgbmFtZT0iTmFtZSINCg0KRGFuaWVsDQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLTM4Mjc0MjU2ODEzMzA5NzUxOTczMTQyMTU5OTkxMg0KQ29udGVudC1EaXNwb3NpdGlvbjogZm9ybS1kYXRhOyBuYW1lPSJTdWJqZWN0Ig0KDQpRdW90ZSByZXF1ZXN0DQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLTM4Mjc0MjU2ODEzMzA5NzUxOTczMTQyMTU5OTkxMg0KQ29udGVudC1EaXNwb3NpdGlvbjogZm9ybS1kYXRhOyBuYW1lPSJNZXNzYWdlIg0KDQpUaGlzIGlzIG15IG1lc3NhZ2UNCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tMzgyNzQyNTY4MTMzMDk3NTE5NzMxNDIxNTk5OTEyDQpDb250ZW50LURpc3Bvc2l0aW9uOiBmb3JtLWRhdGE7IG5hbWU9IlBob3RvIjsgZmlsZW5hbWU9IkZGNEQwMC0wLjgucG5nIg0KQ29udGVudC1UeXBlOiBpbWFnZS9wbmcNCg0KiVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/TQBcNTh/AAAAAXRSTlPM0jRW/QAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYIINCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tMzgyNzQyNTY4MTMzMDk3NTE5NzMxNDIxNTk5OTEyLS0N",
	}

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
		Values: map[string]string{
			"Name":       "Daniel",
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
