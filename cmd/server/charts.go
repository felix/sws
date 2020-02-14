package main

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"src.userspace.com.au/sws"
)

func renderHTMLChart(w http.ResponseWriter, hits []*sws.Hit) error {
	buckets := sws.HitsToTimeBuckets(hits, time.Minute)

	t := template.New("chart")
	t, _ = t.Parse(barChart)
	t.Execute(w, buckets)
	return nil
}

func sparklineHandler(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			log("no site in context")
			http.Error(w, http.StatusText(422), 422)
			return
		}

		startSecs, err := strconv.ParseInt(chi.URLParam(r, "s"), 10, 64)
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
		start := time.Unix(startSecs, 0)
		end := time.Unix(endSecs, 0)

		hits, err := db.GetHits(*site, start, end, nil)
		if err != nil {
			log(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		debug("retrieved", len(hits), "hits")
		data := sws.HitsToTimeBuckets(hits, time.Minute)

		w.Header().Set("Content-Type", "image/svg+xml")
		sws.SparklineSVG(w, data, time.Minute)
	}
}

func svgChartHandler(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startSecs, err := strconv.ParseInt(chi.URLParam(r, "s"), 10, 64)
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
		start := time.Unix(startSecs, 0)
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

		hits, err := db.GetHits(*site, start, end, nil)
		if err != nil {
			panic(err)
		}
		debug("retrieved", len(hits), "hits")

		data := sws.HitsToTimeBuckets(hits, time.Minute)

		w.Header().Set("Content-Type", "image/svg+xml")
		sws.HitChartSVG(w, data, time.Minute)
	}
}

const (
	barChart = `
<dl class="chart">
{{ range .Buckets }}
  <dt class="date">{{ .Time }}</dt>
  <dd class="bar">{{ .Count }}</dd>
{{ end }}
</dl>`
)
