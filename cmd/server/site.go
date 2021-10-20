package main

import (
	"net/http"
	"strings"

	"src.userspace.com.au/render"

	"src.userspace.com.au/sws"
)

func handleSite(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			httpError(w, 422, "no site in context")
			return
		}

		payload := newTemplateData(r)
		payload.Site = site

		filter := createHitFilter(r)
		payload.Begin, payload.End = extractTimeRange(r)
		if payload.Begin == nil || payload.End == nil {
			httpError(w, http.StatusBadRequest, "invalid time range")
			return
		}
		filter["begin"] = *payload.Begin
		filter["end"] = *payload.End

		debug("filter", filter)

		hits, err := db.GetHits(*site, filter)
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		hitSet, err := sws.NewHitSet(sws.FromHits(hits))
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if hitSet != nil {
			hitSet.Fill(payload.Begin, payload.End)
			hitSet.SortByDate()
			payload.Hits = hitSet
		}

		if _, ok := filter["path"]; !ok {
			if ps := sws.NewPageSet(hitSet); ps != nil {
				ps.SortByHits()
				payload.PageSet = ps
			}
		}
		if _, ok := filter["referrer"]; !ok {
			if rs := sws.NewReferrerSet(hitSet, *site); rs != nil {
				rs.SortByHits()
				payload.ReferrerSet = rs
			}
		}
		if _, ok := filter["country"]; !ok {
			if cs := sws.NewCountrySet(hitSet); cs != nil {
				cs.SortByHits()
				payload.CountrySet = cs
			}
		}
		if _, ok := filter["browser"]; !ok {
			if bs := sws.NewBrowserSet(hitSet); bs != nil {
				bs.SortByHits()
				payload.BrowserSet = bs
			}
		}

		rndr.HTML(
			w, 200, payload,
			render.Template("site.tmpl"),
			render.Layout("layout.tmpl"),
		)
	}
}

func handleSiteEdit(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		site, ok := ctx.Value("site").(*sws.Site)
		if !ok {
			httpError(w, 422, "no site in context")
			return
		}

		if r.Method == "POST" {
			site.Name = r.FormValue("name")
			site.Description = r.FormValue("description")
			site.Aliases = r.FormValue("aliases")

			if errs := site.Validate(); len(errs) > 0 {
				log("invalid site:", errs)
				r = flashSet(r, flashError, strings.Join(errs, "<br>"))
			} else if err := db.SaveSite(site); err != nil {
				httpError(w, 500, err.Error())
				return
			}
			r = flashSet(r, flashSuccess, "site updated")
		}

		payload := newTemplateData(r)
		payload.Site = site

		rndr.HTML(
			w, 200, payload,
			render.Template("site.tmpl"),
			render.Layout("layout.tmpl"),
		)
	}
}
