package main

//go:generate go run gen.go

import (
	"os"

	"src.userspace.com.au/templates"
)

func main() {
	tmpl, err := templates.New()
	if err != nil {
		panic(err)
	}
	if _, err = tmpl.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}
