package formailer

import (
	"testing"
)

var testForm = New("test")

func TestSetAndGet(t *testing.T) {
	if testForm != DefaultConfig["test"] {
		t.Error("Unexpected result getting form from config")
	}
}
