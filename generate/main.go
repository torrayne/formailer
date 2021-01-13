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
	out, err := os.Create("template.go")
	if err != nil {
		panic(err)
	}
	out.Write([]byte("package main\n\nconst emailTemplate = `"))
	io.Copy(out, template)
	out.Write([]byte("`\n"))
}
