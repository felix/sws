package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var funcMap = template.FuncMap{
	"sparkline": func(id int) string {
		// This will enable "caching" for an hour
		now := time.Now().Truncate(time.Hour)
		//then := now.Add(-720 * time.Hour)
		then := now.Add(-24 * time.Hour)
		return fmt.Sprintf("/sites/%d/sparklines/%d-%d.svg", id, then.Unix(), now.Unix())
	},
	"timeShort": func(t time.Time) string {
		return t.Format("2006-01-02 15:04")
	},
	"timeLong": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"timeHour": func(t time.Time) string {
		return t.Format("15:04")
	},
	"percent": func(a, b int) float64 {
		return (float64(a) / float64(b)) * 100
	},
	"percentInv": func(a, b int) float64 {
		return 100.0 - ((float64(a) / float64(b)) * 100)
	},
}

func extractTimeRange(r *http.Request) (*time.Time, *time.Time) {
	begin := timePtr(time.Now().Add(-24 * time.Hour))
	end := timePtr(time.Now())
	q := r.URL.Query()
	if b := q.Get("begin"); b != "" {
		if bs, err := strconv.ParseInt(b, 10, 64); err == nil {
			begin = timePtr(time.Unix(bs, 0))
		}
	}
	if e := q.Get("end"); e != "" {
		if es, err := strconv.ParseInt(e, 10, 64); err == nil {
			end = timePtr(time.Unix(es, 0))
		}
	}
	return begin, end
}

func timePtr(t time.Time) *time.Time {
	return &t
}
