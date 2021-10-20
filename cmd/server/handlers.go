package main

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"src.userspace.com.au/render"

	"src.userspace.com.au/sws"
)

type templateData struct {
	Domain   string
	Payload  string
	Endpoint string
	User     *sws.User
	Flash    template.HTML
	Begin    *time.Time
	End      *time.Time

	Query url.Values

	Site        *sws.Site
	Sites       []*sws.Site
	Hits        *sws.HitSet
	PageSet     *sws.PageSet
	BrowserSet  *sws.BrowserSet
	ReferrerSet *sws.ReferrerSet
	CountrySet  *sws.CountrySet
}

func (td templateData) QuerySetEncode(k, v string) template.URL {
	qs, _ := url.ParseQuery(td.Query.Encode())
	if v == "" {
		qs.Del(k)
	} else {
		qs.Set(k, v)
	}
	return template.URL(qs.Encode())
}

func (td templateData) QuerySetContains(s string) bool {
	for k, _ := range td.Query {
		if k == s {
			return true
		}
	}
	return false
}

func newTemplateData(r *http.Request) *templateData {
	out := &templateData{
		Domain:   domain,
		Payload:  "//" + domain + "/sws.js",
		Endpoint: "//" + domain + "/sws.gif",
	}
	if r == nil {
		return out
	}

	flashes := flashGet(r)
	var flash strings.Builder
	for _, f := range flashes {
		flash.WriteString(`<span class="notification is-`)
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

	out.Query = r.URL.Query()

	return out
}

func handleIndex(rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := newTemplateData(r)
		rndr.HTML(
			w, 200, payload,
			render.Template("home.tmpl"),
			render.Layout("layout.tmpl"),
		)
	}
}
