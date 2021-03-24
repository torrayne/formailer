package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/djatwood/formailer"
)

func vercelResponse(w http.ResponseWriter, code int, err error) {
	body := http.StatusText(code)
	if err != nil {
		body = err.Error()
		w.Header().Set("location", w.Header().Get("location")+"?error="+err.Error())
	}

	w.WriteHeader(code)
	w.Write([]byte(body))
}

// Vercel just needs a normal http handler
func Vercel(c formailer.Config, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		vercelResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	body := new(strings.Builder)
	_, err := io.Copy(body, r.Body)
	if err != nil {
		vercelResponse(w, http.StatusInternalServerError, err)
		return
	}

	submission, err := formailer.Parse(r.Header.Get("Content-Type"), body.String())
	if err != nil {
		vercelResponse(w, http.StatusBadRequest, err)
		return
	}

	if submission.Form.ReCAPTCHA {
		v, exists := submission.Values["g-recaptcha-response"].(string)
		if !exists || len(v) < 1 {
			vercelResponse(w, http.StatusBadRequest, nil)
			return
		}

		ok, err := VerifyRecaptcha(v)
		if err != nil {
			vercelResponse(w, http.StatusInternalServerError, err)
			return
		}
		if !ok {
			vercelResponse(w, http.StatusBadRequest, nil)
			return
		}

		delete(submission.Values, "g-recaptcha-response")
	}

	err = submission.Send()
	if err != nil {
		err = fmt.Errorf("failed to send email: %w", err)
		vercelResponse(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := http.StatusOK
	if len(submission.Form.Redirect) > 0 {
		statusCode = http.StatusSeeOther
		w.Header().Add("Location", submission.Form.Redirect)
	}

	vercelResponse(w, statusCode, nil)
}
