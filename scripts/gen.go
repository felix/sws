package main

//go:generate go run gen.go

import (
	//"fmt"
	"os"

	"src.userspace.com.au/templates"
)

func main() {
	tmpl := templates.Must(templates.New(
		templates.EnableHTMLTemplates(),
		templates.Extensions([]string{".tmpl"}),
		templates.Map([]templates.Mapping{
			{Base: "../", Source: "tmpl"},
			{Base: "../", Source: "static", Extensions: []string{".css", ".js"}},
		}),
		//templates.Debug(func(a ...interface{}) { fmt.Fprintln(os.Stderr, a...) }),
	))
	if _, err := tmpl.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}
