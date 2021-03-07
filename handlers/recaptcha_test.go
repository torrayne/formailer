package handlers

import (
	"testing"
)

func TestVerifyRecaptcha(t *testing.T) {
	_, err := VerifyRecaptcha("")
	if err != nil && err != errRecaptchaBadRequest {
		t.Error(err)
	}
}
