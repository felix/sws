package main

import (
	"net/http"
	"strings"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			site := &sws.Site{
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Aliases:     r.FormValue("aliases"),
			}
			if errs := site.Validate(); len(errs) > 0 {
				log("invalid site:", errs)
				r = flashSet(r, flashError, strings.Join(errs, "<br>"))
			} else if err := db.SaveSite(site); err != nil {
				httpError(w, http.StatusInternalServerError, err.Error())
				return
			}
			r = flashSet(r, flashSuccess, "site created")
		}

		sites, err := db.GetSites()
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		payload := newTemplateData(r)
		payload.Sites = sites

		if err := rndr.Render(w, "sites", payload); err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
