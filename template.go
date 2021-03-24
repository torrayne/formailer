package formailer

//go:generate go run generate/main.go

import (
	"html/template"
	"reflect"
)

var templateFuncMap = template.FuncMap{
	"isSlice": isSlice,
}

func isSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind().String() == "slice"
}
