package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"src.userspace.com.au/sws"
)

func sparklineHandler(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			log("no site in context")
			http.Error(w, http.StatusText(422), 422)
			return
		}

		beginSecs, err := strconv.ParseInt(chi.URLParam(r, "b"), 10, 64)
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		endSecs, err := strconv.ParseInt(chi.URLParam(r, "e"), 10, 64)
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		begin := time.Unix(beginSecs, 0)
		end := time.Unix(endSecs, 0)
		filter := map[string]interface{}{
			"begin": begin,
			"end":   end,
		}

		hits, err := db.GetHits(*site, filter)
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		debug("retrieved", len(hits), "hits")
		data, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		// Ensure the buckets start and end at the right time
		data.Fill(&begin, &end)

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public")
		sws.SparklineSVG(w, data, time.Hour)
	}
}

func svgChartHandler(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		beginSecs, err := strconv.ParseInt(chi.URLParam(r, "b"), 10, 64)
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		endSecs, err := strconv.ParseInt(chi.URLParam(r, "e"), 10, 64)
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		begin := time.Unix(beginSecs, 0)
		end := time.Unix(endSecs, 0)

		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		// width := 400
		// height := 200
		// if w := q.Get("width"); w != "" {
		// 	width, _ = strconv.Atoi(w)
		// }
		// if h := q.Get("height"); h != "" {
		// 	height, _ = strconv.Atoi(h)
		// }

		hits, err := db.GetHits(*site, map[string]interface{}{
			"begin": begin, "end": end,
		})
		if err != nil {
			panic(err)
		}
		debug("retrieved", len(hits), "hits")

		data, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		data.Fill(&begin, &end)

		w.Header().Set("Content-Type", "image/svg+xml")
		sws.HitChartSVG(w, data, time.Minute)
	}
}
