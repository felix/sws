package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"

	"src.userspace.com.au/sws"
)

var funcMap = template.FuncMap{
	"piechart": func(siteID int, dataType string, begin, end time.Time) string {
		return fmt.Sprintf("/sites/%d/charts/p-%s-%d-%d.svg", siteID, dataType, begin.Unix(), end.Unix())
	},
	"sparkline": func(id int) string {
		now := time.Now().Truncate(time.Hour)
		then := now.Add(-168 * time.Hour) // 7 days
		return fmt.Sprintf("/sites/%d/charts/s-h-%d-%d.svg", id, then.Unix(), now.Unix())
	},
	"tz": func(s string, t time.Time) time.Time {
		tz, _ := time.LoadLocation(s)
		// TODO error
		return t.In(tz)
	},
	"countryName": func(code string) string {
		if n, ok := sws.CountryCodes[code]; ok {
			return n
		}
		return "Unknown"
	},
	/*
		"seq": func(start, stop, step int) []int {
			count := (stop - start) / step
			out := make([]int, count)
			c := start
			for i := 0; i < count; i++ {
				out[i] = c
				c += step
			}
			return out
		},
		"div": func(a, b int) int {
			return a / b
		},
	*/
	"datetimeShort": func(t time.Time) string {
		return t.Format("2006-01-02 15:04")
	},
	"datetimeLong": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"datetimeHour": func(t time.Time) string {
		return t.Format("15:04 Jan 2")
	},
	"dateRFC": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"timeRFC": func(t time.Time) string {
		return t.Format("15:04")
	},
	"datetimeRelative": func(d string) int64 {
		dur, _ := time.ParseDuration(d)
		return time.Now().Add(dur).Unix()
	},
	"percent": func(a, b int) float64 {
		return (float64(a) / float64(b)) * 100
	},
	"percentInv": func(a, b int) float64 {
		return 100.0 - ((float64(a) / float64(b)) * 100)
	},
	"round": func(n int, a float64) float64 {
		n = n * 10
		return math.Round(a*float64(n)) / float64(n)
	},
}

func httpError(w http.ResponseWriter, code int, msg string) {
	log(msg)
	http.Error(w, http.StatusText(code), code)
}

func extractTimeRange(r *http.Request) (*time.Time, *time.Time) {
	// Default to 1 week ago
	begin := timePtr(time.Now().Truncate(time.Hour).Add(-168 * time.Hour))
	end := timePtr(time.Now().Truncate(30 * time.Minute))
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

func expandPayload(hs *sws.HitSet, pl *templateData) error {
	pl.Hits = hs

	pageSet, err := sws.NewPageSet(hs)
	if err != nil {
		return err
	}

	if pageSet != nil {
		pageSet.SortByHits()
		pl.PageSet = pageSet
	}
	pl.Browsers = sws.NewBrowserSet(hs)
	pl.CountrySet = sws.NewCountrySet(hs)

	refSet := sws.NewReferrerSet(hs)
	if refSet != nil {
		refSet.SortByHits()
		pl.ReferrerSet = refSet
	}
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
