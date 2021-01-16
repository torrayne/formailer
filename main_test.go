package formailer

import (
	"io"
	"os"
	"testing"
)

func testsetup() {
	os.Setenv("FORM_CONTACT_TO", "djatwood01@gmail.com")
	os.Setenv("FORM_CONTACT_FROM", "daniel@atwood.io")
	os.Setenv("FORM_CONTACT_SUBJECT", "New Contact Form Submission")
}

func shutdown() {
	os.Clearenv()
}

func TestMain(m *testing.M) {
	testsetup()
	code := m.Run()
	os.Exit(code)
}

func TestGetForm(t *testing.T) {
	contact := form{
		name:    "contact",
		to:      "djatwood01@gmail.com",
		from:    "daniel@atwood.io",
		subject: "New Contact Form Submission",
	}

	res, err := getForm(contact.name)
	if err != nil {
		t.Errorf("Error getting form: %w", err)
	}
	if contact != res {
		t.Error("Unexpected result reading form from env")
	}
}

func TestParseData(t *testing.T) {
	tests := map[string]string{
		"application/json":                  `{"Name":"Daniel", "message": "This is my message"}`,
		"application/x-www-form-urlencoded": "name=Daniel&message=This is my message",
		"multipart/form-data  boundary=---------------------------382742568133097519731421599912": "LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0zODI3NDI1NjgxMzMwOTc1MTk3MzE0MjE1OTk5MTINCkNvbnRlbnQtRGlzcG9zaXRpb246IGZvcm0tZGF0YTsgbmFtZT0iTmFtZSINCg0KRGFuaWVsDQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLTM4Mjc0MjU2ODEzMzA5NzUxOTczMTQyMTU5OTkxMg0KQ29udGVudC1EaXNwb3NpdGlvbjogZm9ybS1kYXRhOyBuYW1lPSJTdWJqZWN0Ig0KDQpRdW90ZSByZXF1ZXN0DQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLTM4Mjc0MjU2ODEzMzA5NzUxOTczMTQyMTU5OTkxMg0KQ29udGVudC1EaXNwb3NpdGlvbjogZm9ybS1kYXRhOyBuYW1lPSJNZXNzYWdlIg0KDQpUaGlzIGlzIG15IG1lc3NhZ2UNCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tMzgyNzQyNTY4MTMzMDk3NTE5NzMxNDIxNTk5OTEyDQpDb250ZW50LURpc3Bvc2l0aW9uOiBmb3JtLWRhdGE7IG5hbWU9IlBob3RvIjsgZmlsZW5hbWU9IkZGNEQwMC0wLjgucG5nIg0KQ29udGVudC1UeXBlOiBpbWFnZS9wbmcNCg0KiVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/TQBcNTh/AAAAAXRSTlPM0jRW/QAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYIINCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tMzgyNzQyNTY4MTMzMDk3NTE5NzMxNDIxNTk5OTEyLS0N",
	}

	for contentType, body := range tests {
		_, err := parseData(contentType, body)
		if err != nil && err != io.EOF {
			t.Errorf("Failed to parse data: %v", err)
		}
	}
}

func TestFormatData(t *testing.T) {
	form := form{
		name:    "contact",
		subject: "New Contact Form Submission",
	}

	expected := "<h1>New Contact Submission</h1><table><tbody><tr><th>Name</th><td>Daniel</td></tr></tbody></table>"

	data := formData{values: map[string]string{"Name": "Daniel"}}
	output := data.format(form)

	if output != expected {
		t.Error("Failed to format data")
	}
}

func TestGenerateMessage(t *testing.T) {
	form := form{
		name:    "contact",
		subject: "New Contact Form Submission",
	}
	message := "<h1>New Contact Submission</h1><table><tbody><tr><th>Name</th><td>Daniel</td></tr><tr><th>Message</th><td>Hello, World!</td></tr></tbody></table>"

	_, err := generateMessage(form, message)
	if err != nil {
		t.Error(err)
	}
}
