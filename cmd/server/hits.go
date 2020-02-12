package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"src.userspace.com.au/sws"
)

const gif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

func handleHits(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleHitCounter(db sws.CounterStore) http.HandlerFunc {
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

		domain, err := db.GetDomainByName(*hit.Host)
		if err != nil {
			log("failed to get domain", err)
			http.Error(w, "invalid domain", http.StatusNotFound)
			return
		}
		hit.DomainID = domain.ID
		hit.Addr = &r.RemoteAddr

		if err := db.SaveHit(hit); err != nil {
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
