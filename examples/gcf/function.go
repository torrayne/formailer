package gcf

import (
	"net/http"
	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/handlers"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Formailer handles all form submissions
func Formailer(w http.ResponseWriter, r *http.Request) {
	contact := formailer.New("Contact")
	contact.AddEmail(formailer.Email{
		ID:      "contact",
		To:      "info@domain.com",
		From:    `"Company" <noreply@domain.com>`,
		Subject: "New Contact Submission",
	})

	handlers.Vercel(formailer.DefaultConfig, w, r)
}

// Google Cloud Function entry point defined as "main":
func init() {
	functions.HTTP("main", Formailer)
}
