package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djatwood/formailer"
)

func main() {
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

	lambda.Start(cfg.Netlify)
}
