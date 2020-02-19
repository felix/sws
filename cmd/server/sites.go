package main

import (
	"net/http"
	"time"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
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
		if err := rndr.Render(w, "sites", payload); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func handleSite(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
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

		hits, err := db.GetHits(*site, map[string]interface{}{
			"begin": *begin,
			"end":   *end,
		})
		if err != nil {
			log(err)
		}

		pages := sws.PagesFromHits(hits)
		userAgents := sws.UserAgentsFromHits(hits)

		buckets := sws.HitsToTimeBuckets(hits, time.Hour)
		buckets.Fill(begin, end)

		payload := struct {
			Site       *sws.Site
			Pages      map[string]*sws.Page
			UserAgents map[string]*sws.UserAgent
			Hits       sws.TimeBuckets
		}{
			Site:       site,
			Pages:      pages,
			UserAgents: userAgents,
			Hits:       buckets,
		}
		if err := rndr.Render(w, "site", payload); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
