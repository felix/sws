package main

import (
	"net/http"
)

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "counter/test.html")
	}
}

func handleExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html><html><head><meta charset="utf-8"><script>var _sws = { title: "test title" }</script>
<script async src="http://localhost:5000/sws.js" data-sws="http://localhost:5000/sws.gif"></script>
    <title>This is the title</title>
    <noscript><img src="http://localhost:5000/sws.gif" /></noscript></head><body><a href="?referred">test</a></body></html>`))
	}
}
