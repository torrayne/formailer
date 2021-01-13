package main

import (
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
	setup()
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

func TestFormatData(t *testing.T) {
	form := form{
		name:    "contact",
		subject: "New Contact Form Submission",
	}

	expected := "<h1>New Contact Submission</h1><table><tbody><tr><th>Name</th><td>Daniel</td></tr><tr><th>Message</th><td>Hello, World!</td></tr></tbody></table>"

	data := map[string]string{
		"Name":    "Daniel",
		"Message": "Hello, World!",
	}
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
