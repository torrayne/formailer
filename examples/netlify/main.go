package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/handlers"
)

func main() {
	contact := formailer.New("Contact")
	contact.AddEmail(formailer.Email{
		ID:      "contact",
		To:      "info@domain.com",
		From:    `"Company" <noreply@domain.com>`,
		Subject: "New Contact Submission",
	})

	lambda.Start(handlers.Netlify(formailer.DefaultConfig))
}
