package main

import (
	"net/http"
	"time"

	"src.userspace.com.au/sws"
)

type templateData struct {
	User     *sws.User
	Flashes  []flashMsg
	Begin    *time.Time
	End      *time.Time
	Site     *sws.Site
	Sites    []*sws.Site
	Pages    *sws.PageSet
	Browsers *sws.BrowserSet
	Hits     *sws.HitSet
}

func newTemplateData(r *http.Request) *templateData {
	out := &templateData{Flashes: flashGet(r)}
	if user := r.Context().Value("user"); user != nil {
		out.User = user.(*sws.User)
	}
	return out
}

func handleIndex(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		payload := newTemplateData(r)
		if err := rndr.Render(w, "home", payload); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func handleExample(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := rndr.Render(w, "example", nil); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
