package main

import (
	"html/template"
	"net/http"
)

func handleIndex() http.HandlerFunc {
	tmplHome := loadTemplateMust("home")
	tmplNav := loadTemplateMust("partials/navMain")
	tmplLayout := loadTemplateMust("layout")
	tmpl := template.Must(template.New("layout").Parse(string(tmplLayout)))
	_ = template.Must(tmpl.Parse(string(tmplHome)))
	_ = template.Must(tmpl.Parse(string(tmplNav)))
	debug(tmpl)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, nil)
	}
}

func handleExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		b, err := loadTemplate("example")
		if err != nil {
			panic(err)
		}
		w.Write(b)
	}
}
