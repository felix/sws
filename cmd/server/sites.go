package main

import (
	"fmt"
	"net/http"
	"time"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites, err := db.GetSites()
		if err != nil {
			httpError(w, 500, err.Error())
			return
		}
		payload := struct {
			Sites []*sws.Site
		}{
			Sites: sites,
		}
		if err := rndr.Render(w, "sites", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
	}
}

func handleSite(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			httpError(w, 422, "no site in context")
			return
		}
		begin, end := extractTimeRange(r)
		if begin == nil || end == nil {
			httpError(w, 406, "invalid time range")
			return
		}

		hits, err := db.GetHits(*site, map[string]interface{}{
			"begin": *begin,
			"end":   *end,
		})
		if err != nil {
			httpError(w, 500, err.Error())
			return
		}

		hitSet := sws.NewHitSet(hits, *begin, *end, time.Hour)
		fmt.Printf("site begin: %s end: %s\n", *begin, *end)
		hitSet.Fill(begin, end)
		fmt.Printf("hitset begin: %s end: %s\n", hitSet.Begin(), hitSet.End())
		pageSet := sws.NewPageSet(hitSet)
		uaSet := sws.NewUserAgentSet(hitSet)

		payload := struct {
			Site       *sws.Site
			Pages      sws.PageSet
			UserAgents sws.UserAgentSet
			Hits       *sws.HitSet
		}{
			Site:       site,
			Pages:      pageSet,
			UserAgents: uaSet,
			Hits:       hitSet,
		}
		if err := rndr.Render(w, "site", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
	}
}
