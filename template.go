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

func isSlice(v interface{}) bool {
	return "slice" == reflect.TypeOf(v).Kind().String()
}
