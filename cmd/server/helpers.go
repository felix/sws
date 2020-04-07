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
	"piechart": func(dataType string, pl templateData) string {
		pl.Query.Add("begin", strconv.Itoa(int(pl.Begin.Unix())))
		pl.Query.Add("end", strconv.Itoa(int(pl.End.Unix())))
		return fmt.Sprintf("/sites/%d/charts/p-%s.svg?%s", *pl.Site.ID, dataType, pl.Query.Encode())
	},
	"sparkline": func(id int) string {
		now := time.Now().Truncate(time.Hour)
		then := now.Add(-168 * time.Hour) // 7 days
		return fmt.Sprintf("/sites/%d/charts/s-h.svg?begin=%d&end=%d", id, then.Unix(), now.Unix())
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
	"datetimeRelative": func(d string) string {
		dur, _ := time.ParseDuration(d)
		return strconv.Itoa(int(time.Now().Add(dur).Unix()))
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

func createHitFilter(r *http.Request) map[string]interface{} {
	filter := make(map[string]interface{})

	q := r.URL.Query()

	if path := q.Get("path"); path != "" {
		filter["path"] = path
	}
	if country := q.Get("country"); country != "" {
		filter["country_code"] = country
	}
	if browser := q.Get("browser"); browser != "" {
		filter["ua.browser"] = browser
	}
	if referrer := q.Get("referrer"); referrer != "" {
		filter["referrer"] = referrer
	}
	if bots := q.Get("bots"); bots != "" {
		filter["ua.bot"] = (bots == "1")
	}
	if mobile := q.Get("mobile"); mobile != "" {
		filter["ua.mobile"] = (mobile == "1")
	}
	return filter
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

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
