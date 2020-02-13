package main

//go:generate go run gen.go

import (
	"os"

	"src.userspace.com.au/templates"
)

func main() {
	tmpl := templates.Must(templates.New())
	if _, err := tmpl.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}
