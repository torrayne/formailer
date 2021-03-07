package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

var errRecaptchaBadRequest = errors.New(http.StatusText(http.StatusInternalServerError))

type recaptchaResponse struct {
	Success    bool
	ErrorCodes []string `json:"error-codes"`
}

// VerifyRecaptcha verifies the recaptcha response
func VerifyRecaptcha(response string) (bool, error) {
	data := url.Values{}
	data.Set("secret", os.Getenv("RECAPTCHA_SECRET"))
	data.Set("response", response)
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", data)
	if err != nil {
		return false, err
	}

	var body recaptchaResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return false, err
	}

	for _, code := range body.ErrorCodes {
		if code == "timeout-or-duplicate" {
			return body.Success, nil
		}
	}

	if !body.Success {
		return body.Success, errRecaptchaBadRequest
	}

	return body.Success, nil
}
