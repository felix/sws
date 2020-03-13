package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"src.userspace.com.au/sws"
)

const (
	gif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
)

// func handleHits(db sws.HitStore) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		return
// 	}
// }

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
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifBytes)
		return
	}
}

func handleCounter(addr string) http.HandlerFunc {
	counter := getCounter()
	tmpl, err := template.New("counter").Parse(counter)
	if err != nil || tmpl == nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, newTemplateData(nil)); err != nil {
		panic(err)
	}
	etag := fmt.Sprintf(`"%x"`, sha1.Sum(buf.Bytes()))

	return func(w http.ResponseWriter, r *http.Request) {
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		// TODO restrict to site sites
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Etag", etag)
		w.Header().Set("Content-Type", "application/javascript")

		if _, err := io.Copy(w, &buf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
