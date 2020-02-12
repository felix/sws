package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"src.userspace.com.au/sws"
)

const gif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "counter/test.html")
	}
}

func handleDomains(db sws.DomainStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleDomain(db sws.DomainStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleHits(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleHitCounter(db sws.CounterStore) http.HandlerFunc {
	gifBytes, err := base64.StdEncoding.DecodeString(gif)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		hit, err := sws.HitFromRequest(r)
		if err != nil {
			log("failed to create hit", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		domain, err := db.GetDomainByName(*hit.Host)
		if err != nil {
			log("failed to get domain", err)
			http.Error(w, "invalid domain", http.StatusNotFound)
			return
		}
		hit.DomainID = domain.ID
		hit.Addr = &r.RemoteAddr

		if err := db.SaveHit(hit); err != nil {
			log("failed to save hit", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifBytes)
		log("hit", hit)
		return
	}
}

func handleCounter(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		reader := strings.NewReader(counter)
		if _, err := io.Copy(w, reader); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleChart(db sws.HitStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		start, err := time.Parse(time.RFC3339, q.Get("start"))
		if err != nil {
			panic(err)
		}
		end, err := time.Parse(time.RFC3339, q.Get("end"))
		if err != nil {
			panic(err)
		}
		ctx := r.Context()
		domain, ok := ctx.Value("domain").(*sws.Domain)
		if !ok {
			http.Error(w, http.StatusText(422), 422)
			return
		}

		width := 400
		height := 200
		if w := q.Get("width"); w != "" {
			width, _ = strconv.Atoi(w)
		}
		if h := q.Get("height"); h != "" {
			height, _ = strconv.Atoi(h)
		}

		hits, err := db.GetHits(*domain, start, end, nil)
		if err != nil {
			panic(err)
		}
		debug("retrieved", len(hits), "hits")

		// Buckets
		buckets := make(map[int64]float64)
		for _, h := range hits {
			k := h.CreatedAt.Round(time.Minute).UnixNano()
			buckets[k] = buckets[k] + 1
		}

		totalHits := chart.TimeSeries{
			Name: "Hits",
			Style: chart.Style{
				Show:        true,
				StrokeWidth: 4.3,
				StrokeColor: drawing.Color{21, 198, 148, 100},
				DotColorProvider: func(_, _ chart.Range, _ int, _, _ float64) drawing.Color {
					return drawing.Color{21, 198, 148, 100}
				},
				DotWidthProvider: func(_, _ chart.Range, _ int, _, _ float64) float64 {
					return 5
				},
			},
			XValues: make([]time.Time, len(buckets)),
			YValues: make([]float64, len(buckets)),
			//YAxis:   chart.YAxisSecondary,
		}
		xticks := make([]chart.Tick, len(buckets))

		i := 0
		for t, c := range buckets {
			totalHits.XValues[i] = time.Unix(t, 0)
			totalHits.YValues[i] = c
			xticks[i] = chart.Tick{
				Value: float64(t),
				Label: time.Unix(t, 0).Format("Jan 02"),
			}
			i++
		}
		graph := chart.Chart{
			Width:  width,
			Height: height,
			Series: []chart.Series{totalHits},
			Background: chart.Style{
				Padding: chart.Box{Top: 10, Right: 10, Bottom: 10, Left: 27},
			},
			XAxis: chart.XAxis{
				Style: chart.Style{Show: true},
				ValueFormatter: func(v interface{}) string {
					return time.Unix(0, int64(v.(float64))).Format("Jan 02")
				},
				//Ticks: xticks,
			},
			YAxis: chart.YAxis{
				Name:      "Hits",
				NameStyle: chart.StyleShow(),
				Style:     chart.Style{Show: true},
				ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%.0f", v.(float64))
				},
			},
			// YAxis: chart.YAxis{
			// 	Name:      "Unique visitors",
			// 	NameStyle: chart.StyleShow(),
			// 	Style:     chart.Style{Show: true},
			// 	ValueFormatter: func(v interface{}) string {
			// 		return fmt.Sprintf("%.1f", v.(float64))
			// 	},
			// },
		}
		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		graph.Render(chart.SVG, w)
	}
}

func handleExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html><html><head><meta charset="utf-8"><script>var _sws = { title: "test title" }</script>
<script async src="http://localhost:5000/sws.js" data-sws="http://localhost:5000/sws.gif"></script>
    <title>This is the title</title>
    <noscript><img src="http://localhost:5000/sws.gif" /></noscript></head><body><a href="?referred">test</a></body></html>`))
	}
}
