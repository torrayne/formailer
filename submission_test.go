package formailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testSubmissionData = url.Values{
	"_form_name": {"contact"},
	"name":       {"Daniel", "Atwood"},
	"message":    {"This is my message"},
	"urlencoded": {`message=`},
	"multipart":  {`name="message"`},
	"json":       {`"message":`},
}

var expectedSubmissionOrder = []string{
	"_form_name",
	"name", "message",
	"urlencoded", "multipart", "json",
}

func TestParseJSON(t *testing.T) {
	var body string
	for _, key := range expectedSubmissionOrder {
		k, err := json.Marshal(key)
		if err != nil {
			t.Error(err)
		}
		v, err := json.Marshal(testSubmissionData[key])
		if err != nil {
			t.Error(err)
		}
		body += fmt.Sprintf(`,%s:%s`, k, v)
	}
	body = "{" + body[1:] + "}"

	submission := new(Submission)
	submission.Values = make(map[string]interface{})
	err := submission.parseJSON(string(body))
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(submission.Order, expectedSubmissionOrder) {
		t.Errorf("Submission has incorrect order\n%v", submission.Order)
	}
}
func TestParseURLEncoded(t *testing.T) {
	var body string
	for _, key := range expectedSubmissionOrder {
		for _, value := range testSubmissionData[key] {
			k := url.QueryEscape(key)
			v := url.QueryEscape(value)
			body += fmt.Sprintf("&%s=%s", k, v)
		}
	}

	submission := new(Submission)
	submission.Values = make(map[string]interface{})
	err := submission.parseURLEncoded(body[1:])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(submission.Order, expectedSubmissionOrder) {
		t.Errorf("Submission has incorrect order\n%v", submission.Order)
	}
}
func TestParseMultipartForm(t *testing.T) {
	boundary := "--myspecificboundary"
	b := bytes.NewBuffer([]byte{})
	w := multipart.NewWriter(b)
	w.SetBoundary(boundary)
	for _, k := range expectedSubmissionOrder {
		for _, v := range testSubmissionData[k] {
			w.WriteField(k, v)
		}
	}
	w.Close()

	submission := new(Submission)
	submission.Values = make(map[string]interface{})
	err := submission.parseMultipartForm(boundary, b.String())
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(submission.Order, expectedSubmissionOrder) {
		t.Errorf("Submission has incorrect order\n%v\n%v", submission.Order, expectedSubmissionOrder)
	}
}
