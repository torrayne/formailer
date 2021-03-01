package formailer

import (
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

func or(a, b string) string {
	if len(a) < 1 {
		return b
	}
	return a
}

func isSlice(v interface{}) bool {
	return "slice" == reflect.TypeOf(v).Kind().String()
}
