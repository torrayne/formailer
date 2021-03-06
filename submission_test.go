package formailer

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"
)

func TestParseJSON(t *testing.T) {
	body := `{"_form_name": "contact", "Name":"Daniel", "message": "This is my message"}`
	submission := new(Submission)
	submission.Values = make(map[string]interface{})
	err := submission.parseJSON(body)
	if err != nil {
		t.Error(err)
	}
}
func TestParseURLEncoded(t *testing.T) {
	body := "_form_name=contact&name=Daniel&message=This is my message"
	submission := new(Submission)
	submission.Values = make(map[string]interface{})
	err := submission.parseURLEncoded(body)
	if err != nil {
		t.Error(err)
	}
}
func TestParseMultipartForm(t *testing.T) {
	boundary := "--myspecificboundary"
	b := bytes.NewBuffer([]byte{})
	w := multipart.NewWriter(b)
	w.SetBoundary(boundary)
	w.WriteField("_form_name", "contact")
	w.WriteField("Name", "Daniel")
	w.WriteField("message", "This is my message")
	w.Close()

	submission := new(Submission)
	submission.Values = make(map[string]interface{})
	err := submission.parseMultipartForm("multipart/form-data  boundary="+boundary, b.String())
	if err != nil && err != io.EOF {
		t.Error(err)
	}
}
