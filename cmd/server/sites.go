package main

import (
	"html/template"
	"net/http"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites, err := db.GetSites()
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
		payload := struct {
			Sites []*sws.Site
		}{
			Sites: sites,
		}
		tmpls.ExecuteTemplate(w, "sites", payload)
	}
}

func handleSite(db sws.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
