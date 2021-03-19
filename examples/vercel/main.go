package api

import (
	"net/http"

	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/handlers"
)

// Formailer handles all form submissions
func Formailer(w http.ResponseWriter, r *http.Request) {
	contact := formailer.Form{Name: "Contact"}
	contact.AddEmail(formailer.Email{
		ID:      "contact",
		To:      "info@domain.com",
		From:    `"Company" <noreply@domain.com>`,
		Subject: "New Contact Submission",
	})

	formailer.Add(contact)
	handlers.Vercel(formailer.Forms, w, r)
}
