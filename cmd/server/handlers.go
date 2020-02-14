package main

import (
	"html/template"
	"net/http"
)

func handleIndex(tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpls.ExecuteTemplate(w, "home", nil)
	}
}

func handleExample(tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		tmpls.ExecuteTemplate(w, "example", nil)
	}
}
