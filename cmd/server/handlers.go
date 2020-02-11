package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"src.userspace.com.au/sws"
	//"src.userspace.com.au/sws/counter"
)

const gif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "counter/test.html")
	}
}

func handleDomains(db sws.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleHits(db sws.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleHitCounter(db sws.Queryer) http.HandlerFunc {
	gifBytes, err := base64.StdEncoding.DecodeString(gif)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		hit, err := sws.HitFromRequest(r)
		if err != nil {
			log("failed to create hit", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		domain, err := sws.GetDomainByName(db, *hit.Host)
		if err != nil {
			log("failed to get domain", err)
			http.Error(w, "invalid domain", http.StatusNotFound)
			return
		}
		hit.DomainID = domain.ID
		hit.Addr = &r.RemoteAddr

		if err := hit.Save(db); err != nil {
			log("failed to save hit", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifBytes)
		log("hit", hit)
		return
	}
}

func handleCounter(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		reader := strings.NewReader(counter)
		if _, err := io.Copy(w, reader); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
