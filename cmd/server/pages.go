package main

import (
	"net/http"

	"src.userspace.com.au/sws"
)

func handlePages(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			httpError(w, http.StatusBadRequest, "no site in context")
			return
		}

		payload := newTemplateData(r)
		payload.Site = site

		begin, end := extractTimeRange(r)
		if begin == nil || end == nil {
			httpError(w, http.StatusBadRequest, "invalid time range")
			return
		}
		payload.Begin = *begin
		payload.End = *end
		debug("begin", *begin)
		debug("end", *end)

		filter := map[string]interface{}{
			"begin": *begin,
			"end":   *end,
		}

		q := r.URL.Query()

		path := q.Get("path")
		if path != "" {
			filter["path"] = path
		}

		hits, err := db.GetHits(*site, filter)
		hitSet, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		hitSet.Fill(begin, end)
		hitSet.SortByDate()

		payload.Page = sws.NewPage(hitSet)
		payload.Hits = hitSet
		payload.Browsers = sws.NewBrowserSet(hitSet)
		payload.CountrySet = sws.NewCountrySet(hitSet)
		payload.ReferrerSet = sws.NewReferrerSet(hitSet)

		// Single or multiple paths
		if path == "" {
			pageSet, err := sws.NewPageSet(hitSet)
			if err != nil {
				httpError(w, http.StatusInternalServerError, err.Error())
				return
			}

			pageSet.SortByHits()
			payload.PageSet = pageSet
		}

		if err := rndr.Render(w, "pages", payload); err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
