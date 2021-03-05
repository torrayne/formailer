package api

import (
	"net/http"

	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/handlers"
)

// Formailer handles all form submissions
func Formailer(w http.ResponseWriter, r *http.Request) {
	cfg := make(formailer.Config)
	cfg.Set(
		&formailer.Form{
			To:       "support@domain.com",
			From:     `"Company" <noreply@domain.com>`,
			Subject:  "New Submission",
			Redirect: "/success",
		}, &formailer.Form{
			Name:     "Contact",
			To:       "info@domain.com",
			From:     `"Company" <noreply@domain.com>`,
			Subject:  "New Contact Submission",
			Redirect: "https://domin.com/thankyou",
		},
	)

	handlers.Vercel(&cfg, w, r)
}
