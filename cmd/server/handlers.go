package main

import (
	"encoding/base64"
	"net/http"
	"strings"

	"src.userspace.com.au/sws/models"
)

const transGif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "client/test.html")
	}
}

func handleSnippet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		//w.Header().Set("Expires", "600")
		w.Write([]byte(snippet))
		//http.ServeFile(w, r, "client/sws.min.js")
	}
}

func handlePageView(db models.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pv, err := models.PageViewFromURL(*r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		domain, err := models.GetDomainByName(db, *pv.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		pv.DomainID = domain.ID
		pv.Address = &(strings.Split(r.RemoteAddr, ":")[0])

		if err := pv.Save(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, _ := base64.StdEncoding.DecodeString(transGif)
		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Write(b)
		return
	}
}
