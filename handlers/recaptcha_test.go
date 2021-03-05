package handlers

import (
	"testing"
)

func TestVerifyRecaptcha(t *testing.T) {
	_, err := verifyRecaptcha("")
	if err != nil && err != ErrRecaptchaBadRequest {
		t.Error(err)
	}
}
