package main

import (
	"net/http"

	"src.userspace.com.au/sws"
)

func handleCountries(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
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

		cc := q.Get("country")
		if cc != "" {
			filter["countryCode"] = cc
		}

		hits, err := db.GetHits(*site, filter)
		hitSet, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		hitSet.Fill(begin, end)
		hitSet.SortByDate()
		if err := expandPayload(hitSet, payload); err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Single or multiple paths
		if cc == "" {
			countrySet := sws.NewCountrySet(hitSet)
			countrySet.SortByHits()
			payload.CountrySet = countrySet
		}

		if err := rndr.Render(w, "site", payload); err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
