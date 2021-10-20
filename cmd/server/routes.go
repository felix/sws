package main

import (
	"context"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"src.userspace.com.au/render"

	"src.userspace.com.au/sws"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("lkjasd0f9u203ijsldkfj"), nil)
}

func createRouter(db sws.Store, mmdbPath string) (chi.Router, error) {

	src, err := fs.Sub(fs.FS(tmpl), "tmpl")
	if err != nil {
		return nil, err
	}
	rndr, err := render.New(render.AddTemplates(src, render.Root("root")))
	if err != nil {
		return nil, err
	}

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
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Get("/sws.gif", handleHitCounter(db, mmdbPath))
	})
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
					rndr.HTML(
						w, 200, payload,
						render.Template("site.tmpl"),
						render.Layout("layout.tmpl"),
					)
				}))
				r.Route("/{siteID}", func(r chi.Router) {
					siteHandler := handleSite(db, rndr)
					// Populate context with site if present
					r.Use(getSiteCtx(db))
					r.Get("/", siteHandler)
					r.Post("/", siteHandler)
					r.Get("/edit", handleSiteEdit(db, rndr))

					r.Route("/charts", func(r chi.Router) {
						r.Get("/{type:(p|s|b)}-{data:(h|b|c)}.svg", chartHandler(db))
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

		// Static files
		r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := strings.TrimPrefix(r.URL.Path, "/")
			debug("loading static", p)
			// b, err := loadTemplate(p)
			// if err == nil {
			// 	name := filepath.Base(p)
			// 	etag := fmt.Sprintf(`"%x"`, sha1.Sum(b))

			// 	if match := r.Header.Get("If-None-Match"); match != "" {
			// 		if strings.Contains(match, etag) {
			// 			w.WriteHeader(http.StatusNotModified)
			// 			return
			// 		}
			// 	}

			// 	w.Header().Set("Etag", etag)
			// 	w.Header().Set("Cache-Control", "no-cache")
			// 	http.ServeContent(w, r, name, time.Now(), bytes.NewReader(b))
			// }
			debug("no template found, trying files")
			fs := http.FileServer(http.Dir("public"))
			fs.ServeHTTP(w, r)
			//log("file not found:", p, err)
		}))
	})

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := newTemplateData(r)
		rndr.HTML(
			w, 404, payload,
			render.Template("404.tmpl"),
			render.Layout("layout.tmpl"),
		)
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
				//log("failed to extract user from context", err)
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
