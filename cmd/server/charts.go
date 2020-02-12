package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"src.userspace.com.au/sws"
)

type bucket struct {
	key  time.Time
	hits []*sws.Hit
}

func (b bucket) String() string {
	return fmt.Sprintf("%s => %d", b.key, len(b.hits))
}

func hitsToTimeBuckets(hits []sws.Hit, d time.Duration) []bucket {
	out := make([]bucket, 0)
	for _, h := range hits {
		k := h.CreatedAt.Round(d)
		var found bool
		for i, tb := range out {
			if tb.key.Equal(k) {
				out[i].hits = append(out[i].hits, &h)
				found = true
			}
		}
		if !found {
			out = append(out, bucket{key: k, hits: []*sws.Hit{&h}})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].key.Before(out[j].key)
	})
	return out
}

//func renderHTMLChart(w http.ResponseWriter, hits []sws.Hit) error {
//}

func svgChartHandler(db sws.HitStore) http.HandlerFunc {
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

		buckets := hitsToTimeBuckets(hits, time.Minute)
		debug(buckets)

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
		max := float64(0)
		for _, b := range buckets {
			v := float64(len(b.hits))
			totalHits.XValues[i] = b.key
			totalHits.YValues[i] = v
			xticks[i] = chart.Tick{
				Value: float64(b.key.UnixNano()),
				Label: b.key.Format("04"),
			}
			i++
			if v > max {
				max = v
			}
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
				Ticks: xticks,
			},
			YAxis: chart.YAxis{
				Name:      "Hits",
				NameStyle: chart.StyleShow(),
				Style:     chart.Style{Show: true},
				ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%.0f", v.(float64))
				},
				Range: &chart.ContinuousRange{
					Min: 0.0,
					Max: max,
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
		graph.YAxis.Range.SetMin(0)
		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		graph.Render(chart.SVG, w)
	}
}
