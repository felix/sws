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
	"tz": func(s string, t time.Time) time.Time {
		tz, _ := time.LoadLocation(s)
		// TODO error
		return t.In(tz)
	},
	"timeShort": func(t time.Time) string {
		return t.Format("2006-01-02 15:04")
	},
	"timeLong": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"timeHour": func(t time.Time) string {
		return t.Format("15:04 Jan 2")
	},
	"percent": func(a, b int) float64 {
		return (float64(a) / float64(b)) * 100
	},
	"percentInv": func(a, b int) float64 {
		return 100.0 - ((float64(a) / float64(b)) * 100)
	},
}

func httpError(w http.ResponseWriter, code int, msg string) {
	log(msg)
	http.Error(w, http.StatusText(code), code)
}

func extractTimeRange(r *http.Request) (*time.Time, *time.Time) {
	begin := timePtr(time.Now().Truncate(time.Hour).Add(-168 * time.Hour))
	end := timePtr(time.Now())
	q := r.URL.Query()
	if b := q.Get("begin"); b != "" {
		if bs, err := strconv.ParseInt(b, 10, 64); err == nil {
			begin = timePtr(time.Unix(bs, 0).Truncate(time.Hour))
		}
	}
	if e := q.Get("end"); e != "" {
		if es, err := strconv.ParseInt(e, 10, 64); err == nil {
			end = timePtr(time.Unix(es, 0).Truncate(time.Hour))
		}
	}
	return begin, end
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
