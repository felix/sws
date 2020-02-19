package main

import (
	"net/http"
)

func handleIndex(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := rndr.Render(w, "home", nil); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func handleExample(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := rndr.Render(w, "example", nil); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
