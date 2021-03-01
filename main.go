package formailer

import (
	"os"
	"reflect"
	"text/template"

	// embed is used to emed the default template
	_ "embed"
)

//go:embed template.html
var defaultTemplate string

var templateFuncMap = template.FuncMap{
	"isSlice": isSlice,
}

func defaultSMTP() smtpAuth {
	return smtpAuth{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
	}
}

func or(a, b string) string {
	if len(a) < 1 {
		return b
	}
	return a
}

func isSlice(v interface{}) bool {
	return "slice" == reflect.TypeOf(v).Kind().String()
}
