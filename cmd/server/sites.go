package main

import (
	"net/http"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites, err := db.GetSites()
		if err != nil {
			httpError(w, 500, err.Error())
			return
		}

		payload := newTemplateData(r)
		payload.Sites = sites

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

		hitSet, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			httpError(w, 406, err.Error())
			return
		}
		hitSet.Fill(begin, end)
		hitSet.SortByDate()

		pageSet, err := sws.NewPageSet(hitSet)
		if err != nil {
			httpError(w, 406, err.Error())
			return
		}
		pageSet.SortByHits()

		browserSet := sws.NewBrowserSet(hitSet)

		payload := newTemplateData(r)
		payload.Site = site
		payload.PageSet = &pageSet
		payload.Browsers = &browserSet
		payload.Hits = hitSet

		if err := rndr.Render(w, "site", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
	}
}
