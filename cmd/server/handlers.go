package main

import (
	"net/http"
	"strings"

	"src.userspace.com.au/sws/models"
)

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "client/test.html")
	}
}

func handleDomains(db models.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handlePageViews(db models.Queryer) http.HandlerFunc {
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
		return
	}
}
