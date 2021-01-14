package formailer

import (
	"fmt"
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

	if contact != getForm(contact.name) {
		t.Error("Failed to read form from env")
	}
}

func TestParseData(t *testing.T) {
	// var out bytes.Buffer
	// w := multipart.NewWriter(&out)

	// formData := map[string]string{
	// 	"name":    "Daniel",
	// 	"subject": "Free Consultation",
	// }

	// for name, value := range formData {
	// 	fw, err := w.CreateFormField(name)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	fw.Write([]byte(value))
	// }

	// w.Close()

	// _, err := parseData(w.FormDataContentType(), out.String())
	// if err != nil {
	// 	t.Errorf("Failed to parse data: %v", err)
	// }

	body := "LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0zNTk3MDk5NzAzNjY4MTM2NDQzMjUyMzE0NjAyDQpDb250ZW50LURpc3Bvc2l0aW9uOiBmb3JtLWRhdGE7IG5hbWU9Ik5hbWUiDQoNCkRhbmllbA0KLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0zNTk3MDk5NzAzNjY4MTM2NDQzMjUyMzE0NjAyDQpDb250ZW50LURpc3Bvc2l0aW9uOiBmb3JtLWRhdGE7IG5hbWU9IlN1YmplY3QiDQoNCm5ldyBzdWJqZWN0DQotLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLTM1OTcwOTk3MDM2NjgxMzY0NDMyNTIzMTQ2MDINCkNvbnRlbnQtRGlzcG9zaXRpb246IGZvcm0tZGF0YTsgbmFtZT0iTWVzc2FnZSINCg0KdGhpcyBpcyB0aGUgbWVzc2FnZQ0KLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0zNTk3MDk5NzAzNjY4MTM2NDQzMjUyMzE0NjAyDQpDb250ZW50LURpc3Bvc2l0aW9uOiBmb3JtLWRhdGE7IG5hbWU9IlBob3RvIjsgZmlsZW5hbWU9IkZGNEQwMC0wLjgucG5nIg0KQ29udGVudC1UeXBlOiBpbWFnZS9wbmcNCg0KiVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/TQBcNTh/AAAAAXRSTlPM0jRW/QAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYIINCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tMzU5NzA5OTcwMzY2ODEzNjQ0MzI1MjMxNDYwMi0t"
	contentType := "multipart/form-data; boundary=---------------------------3597099703668136443252314602"
	data, attachments, err := parseData(contentType, body)
	if err != nil {
		t.Errorf("Failed to parse data: %v", err)
	}

	fmt.Println(data, attachments)
}

func TestFormatData(t *testing.T) {
	form := form{
		name:    "contact",
		subject: "New Contact Form Submission",
	}

	expected := "<h1>New Contact Submission</h1><table><tbody><tr><th>Name</th><td>Daniel</td></tr></tbody></table>"

	data := map[string]string{"Name": "Daniel"}
	output := formatData(form, data)

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
