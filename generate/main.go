package main

import (
	"io"
	"os"
)

func main() {
	template, err := os.Open("template.html")
	if err != nil {
		panic(err)
	}
	out, err := os.Create("template_default.go")
	if err != nil {
		panic(err)
	}
	out.Write([]byte("package formailer\n\nconst defaultTemplate = `"))
	io.Copy(out, template)
	out.Write([]byte("`\n"))
}
