package main

import (
	"net/http"
)

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		b, err := loadTemplate("example")
		if err != nil {
			panic(err)
		}
		w.Write(b)
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
