package main

import (
	"html/template"
	"net/http"
	"time"

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

func handleSite(db sws.SiteStore, tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			log("no site in context")
			http.Error(w, http.StatusText(422), 422)
			return
		}
		begin, end := extractTimeRange(r)
		if begin == nil || end == nil {
			log("invalid time range")
			http.Error(w, http.StatusText(406), 406)
			return
		}

		pages, err := db.GetPages(*site, map[string]interface{}{"begin": *begin, "end": *end})
		if err != nil {
			log(err)
		}

		hits, err := db.GetHits(*site, map[string]interface{}{
			"begin": *begin,
			"end":   *end,
		})
		if err != nil {
			log(err)
		}
		buckets := sws.HitsToTimeBuckets(hits, time.Hour)
		buckets.Fill(begin, end)

		payload := struct {
			Site  *sws.Site
			Pages []*sws.Page
			Hits  sws.TimeBuckets
		}{
			Site:  site,
			Pages: pages,
			Hits:  buckets,
		}
		tmpls.ExecuteTemplate(w, "site", payload)
	}
}
