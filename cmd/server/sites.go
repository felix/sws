package main

import (
	"net/http"
	"strings"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore, rndr Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			site := &sws.Site{
				Name:        r.FormValue("name"),
				Description: r.FormValue("description"),
				Aliases:     r.FormValue("aliases"),
			}
			if errs := site.Validate(); len(errs) > 0 {
				log("invalid site:", errs)
				r = flashSet(r, flashError, strings.Join(errs, "<br>"))
			} else if err := db.SaveSite(site); err != nil {
				httpError(w, http.StatusInternalServerError, err.Error())
				return
			}
			r = flashSet(r, flashSuccess, "site created")
		}

		sites, err := db.GetSites()
		if err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}

		payload := newTemplateData(r)
		payload.Sites = sites

		if err := rndr.Render(w, "sites", payload); err != nil {
			httpError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

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

		begin, end := extractTimeRange(r)
		if begin == nil || end == nil {
			httpError(w, http.StatusBadRequest, "invalid time range")
			return
		}
		payload.Begin = *begin
		payload.End = *end
		debug("begin", *begin)
		debug("end", *end)

		hits, err := db.GetHits(*site, map[string]interface{}{
			"begin": *begin,
			"end":   *end,
		})
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
			hitSet.Fill(begin, end)
			hitSet.SortByDate()
			if err := expandPayload(hitSet, payload); err != nil {
				httpError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if err := rndr.Render(w, "site", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
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

		if err := rndr.Render(w, "site", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
	}
}
