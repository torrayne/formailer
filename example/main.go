package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djatwood/formailer"
)

func main() {
	lambda.Start(formailer.Handler)
}
