package api

import (
	"net/http"

	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/handlers"
)

// Formailer handles all form submissions
func Formailer(w http.ResponseWriter, r *http.Request) {
	forms := make(formailer.Forms)
	forms.Add("Contact", formailer.Email{
		ID:      "contact",
		To:      "info@domain.com",
		From:    `"Company" <noreply@domain.com>`,
		Subject: "New Contact Submission",
	})

	handlers.Vercel(forms, w, r)
}
