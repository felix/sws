package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"src.userspace.com.au/sws"
	"src.userspace.com.au/templates"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("lkjasd0f9u203ijsldkfj"), nil)
}

func createRouter(db sws.Store, mmdbPath string) (chi.Router, error) {
	tmplsCommon := []string{"flash.tmpl", "navbar.tmpl"}
	tmplsAuthed := append(tmplsCommon, []string{"layout.tmpl", "charts.tmpl", "timerange.tmpl"}...)
	tmplsPublic := append(tmplsCommon, "layout.tmpl")

	if override != "" {
		log("using overrider", override)
		loadOverrider = func(s string) string {
			return filepath.Join(override, s)
		}
	}

	tmpls, err := loadHTMLTemplateMap(map[string][]string{
		"sites": append([]string{"sites.tmpl"}, tmplsAuthed...),
		"site":  append([]string{"site.tmpl", "worldMap.tmpl"}, tmplsAuthed...),
		"home":  append([]string{"home.tmpl"}, tmplsPublic...),
		"login": append([]string{"login.tmpl"}, tmplsPublic...),
		"user":  append([]string{"user.tmpl"}, tmplsAuthed...),
		"404":   append([]string{"404.tmpl"}, tmplsPublic...),
	}, funcMap)
	if err != nil {
		return nil, err
	}
	debug(tmpls["login"].DefinedTemplates())
	debug(tmpls["home"].DefinedTemplates())

	rndr := templates.NewRenderer(tmpls)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestID)
	compressor := middleware.NewCompressor(5, "text/html", "text/css")
	r.Use(compressor.Handler())
	if verbose {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	// For counter
	r.Get("/sws.js", handleCounter(addr))
	r.Get("/sws.gif", handleHitCounter(db, mmdbPath))
	//r.Get("/hits", handleHits(db))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		// Populate contect with user if present
		r.Use(getUserCtx(db))

		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/", handleIndex(rndr))
			r.Get(loginURL, handleLogin(db, rndr))
		})

		r.Post(loginURL, handleLogin(db, rndr))

		// Static files
		r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := strings.TrimPrefix(r.URL.Path, "/")
			debug("loading", p)
			if b, err := loadTemplate(p); err == nil {
				name := filepath.Base(p)
				etag := fmt.Sprintf(`"%x"`, sha1.Sum(b))

				if match := r.Header.Get("If-None-Match"); match != "" {
					if strings.Contains(match, etag) {
						w.WriteHeader(http.StatusNotModified)
						return
					}
				}

				w.Header().Set("Etag", etag)
				w.Header().Set("Cache-Control", "no-cache")
				http.ServeContent(w, r, name, time.Now(), bytes.NewReader(b))
			}
		}))

		// Authed routes
		r.Group(func(r chi.Router) {
			// Ensure we have a user in context
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if user := r.Context().Value("user"); user == nil {
						authRedirect(w, r, "authentication required")
					}
					next.ServeHTTP(w, r)
				})
			})

			r.Get(logoutURL, handleLogout(rndr))
			r.Route("/sites", func(r chi.Router) {
				sitesHandler := handleSites(db, rndr)
				r.Get("/", sitesHandler)
				r.Post("/", sitesHandler)
				r.Get("/new", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					payload := newTemplateData(r)
					payload.Site = &sws.Site{}
					if err := rndr.Render(w, "site", payload); err != nil {
						httpError(w, 500, err.Error())
						return
					}
				}))
				r.Route("/{siteID}", func(r chi.Router) {
					siteHandler := handleSite(db, rndr)
					// Populate context with site if present
					r.Use(getSiteCtx(db))
					r.Get("/", siteHandler)
					r.Post("/", siteHandler)
					r.Get("/edit", handleSiteEdit(db, rndr))
					r.Route("/sparklines", func(r chi.Router) {
						r.Get("/{b:\\d+}-{e:\\d+}.svg", sparklineHandler(db))
					})
					r.Route("/charts", func(r chi.Router) {
						r.Get("/{b:\\d+}-{e:\\d+}.svg", svgChartHandler(db))
						//r.Get("/{b:\\d+}-{e:\\d+}.png", svgChartHandler(db))
					})
				})
			})
			r.Route("/users", func(r chi.Router) {
				userH := handleUsers(db, rndr)
				r.Route("/{email}", func(r chi.Router) {
					r.Get("/", userH)
					r.Post("/", userH)
				})
			})
		})
	})

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := newTemplateData(r)
		if err := rndr.Render(w, "404", payload); err != nil {
			httpError(w, 500, err.Error())
			return
		}
	}))

	return r, nil
}

func getUserCtx(db sws.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				next.ServeHTTP(w, r)
			}()

			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				log("failed to extract user from context", err)
				return
			}

			// Token is authenticated, get claims

			id, ok := claims["user_id"]
			if !ok {
				log("missing user ID")
				return
			}

			user, err := db.GetUserByID(int(id.(float64)))
			if err != nil {
				log("missing user")
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)
		})
	}
}

func getSiteCtx(db sws.SiteStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(chi.URLParam(r, "siteID"))
			if err != nil {
				panic(err)
			}
			site, err := db.GetSiteByID(id)
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			ctx := context.WithValue(r.Context(), "site", site)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
