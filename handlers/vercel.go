package handlers

import (
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
func Vercel(c *formailer.Config, w http.ResponseWriter, r *http.Request) {
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

	submission, err := c.Parse(r.Header.Get("Content-Type"), body.String())
	if err != nil && err != io.EOF {
		vercelResponse(w, http.StatusBadRequest, err)
		return
	}

	if v, ok := submission.Values["g-recaptcha-response"]; ok {
		ok, err := verifyRecaptcha(v.(string))
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

	server, err := submission.Form.SMTPServer()
	if err != nil {
		vercelResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = submission.Send(server)
	if err != nil {
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
