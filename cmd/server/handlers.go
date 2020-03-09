package main

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"src.userspace.com.au/sws"
)

type templateData struct {
	Payload     string
	Endpoint    string
	User        *sws.User
	Flash       template.HTML
	Begin       time.Time
	End         time.Time
	Site        *sws.Site
	Sites       []*sws.Site
	PageSet     sws.PageSet
	Browsers    sws.BrowserSet
	ReferrerSet sws.ReferrerSet
	Hits        *sws.HitSet
}

func newTemplateData(r *http.Request) *templateData {
	out := &templateData{
		Payload:  "//" + domain + "/sws.js",
		Endpoint: "//" + domain + "/sws.gif",
	}
	if r != nil {
		flashes := flashGet(r)
		var flash strings.Builder
		for _, f := range flashes {
			flash.WriteString(`<span class="`)
			flash.WriteString(string(f.Level))
			flash.WriteString(`">`)
			flash.WriteString(f.Message)
			flash.WriteString("</span>")
		}
		if len(flashes) > 0 {
			out.Flash = template.HTML(flash.String())
		}

		if user := r.Context().Value("user"); user != nil {
			out.User = user.(*sws.User)
		}
	}
	return out
}

func handleIndex(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := newTemplateData(r)
		if err := rndr.Render(w, "home", payload); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func handleExample(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := rndr.Render(w, "example", nil); err != nil {
			log(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
