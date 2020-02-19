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
			log("failed to extract hit", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		site, err := db.GetSiteByName(hit.Host)
		if err != nil {
			log("failed to get site", err)
			http.Error(w, "invalid site", http.StatusNotFound)
			return
		}
		hit.SiteID = site.ID
		hit.Addr = r.RemoteAddr

		if err := db.SaveHit(hit); err != nil {
			log("failed to save hit", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}
		// TODO restrict to site sites
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifBytes)
		log("hit", hit)
		return
	}
}

func handleCounter(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO restrict to site sites
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/javascript")
		reader := strings.NewReader(counter)
		if _, err := io.Copy(w, reader); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
