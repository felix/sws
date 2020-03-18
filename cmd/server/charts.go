package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"src.userspace.com.au/sws"
)

func chartHandler(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			log("no site in context")
			http.Error(w, http.StatusText(422), 422)
			return
		}

		chartType := chi.URLParam(r, "type")
		dataType := chi.URLParam(r, "data")
		beginSecs, err := strconv.ParseInt(chi.URLParam(r, "begin"), 10, 64)
		if err != nil {
			httpError(w, http.StatusNotFound, err.Error())
			return
		}
		endSecs, err := strconv.ParseInt(chi.URLParam(r, "end"), 10, 64)
		if err != nil {
			httpError(w, http.StatusNotFound, err.Error())
			return
		}
		begin := time.Unix(beginSecs, 0)
		end := time.Unix(endSecs, 0)

		var b strings.Builder
		b.WriteString(r.URL.Path)
		// FIXME
		b.WriteString(time.Now().Truncate(30 * time.Minute).String())
		etag := fmt.Sprintf(`"%x"`, sha1.Sum([]byte(b.String())))

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		filter := map[string]interface{}{
			"begin": begin,
			"end":   end,
		}

		hits, err := db.GetHits(*site, filter)
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		hitSet, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if hitSet == nil {
			httpError(w, http.StatusInternalServerError, "missing hitset")
			return
		}

		hitSet.Fill(&begin, &end)
		hitSet.SortByDate()

		w.Header().Set("Etag", etag)
		w.Header().Set("Cache-Control", "no-cache")

		switch dataType {
		case "h":
			w.Header().Set("Content-Type", "image/svg+xml")
			switch chartType {
			case "b":
				sws.HitChartSVG(w, hitSet, time.Minute)
			case "s":
				sws.SparklineSVG(w, hitSet, time.Hour)
			}
		case "p":
			pages := sws.NewBrowserSet(hitSet)
			pages.SortByHits()
			w.Header().Set("Content-Type", "image/svg+xml")
			switch chartType {
			case "p":
				sws.PieChartSVG(w, pages)
			}
		case "b":
			browsers := sws.NewBrowserSet(hitSet)
			browsers.SortByHits()
			w.Header().Set("Content-Type", "image/svg+xml")
			switch chartType {
			case "p":
				sws.PieChartSVG(w, browsers)
			}
		case "c":
		default:
			log("invalid chart data type:", dataType)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}
