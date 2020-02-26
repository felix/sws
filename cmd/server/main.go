package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"src.userspace.com.au/sws"
	"src.userspace.com.au/sws/store"
	"src.userspace.com.au/templates"
)

var (
	Version    string
	log, debug sws.Logger
	tokenAuth  *jwtauth.JWTAuth
)

const (
	endpoint  = "//stats.userspace.com.au/sws.gif"
	loginURL  = "/login"
	logoutURL = "/logout"
)

// Flags
var (
	verbose   *bool
	addr      *string
	dsn       *string
	noMigrate *bool
)

func init() {
	verbose = boolFlag("verbose", "v", false, "VERBOSE", "enable verbose output")
	addr = stringFlag("listen", "l", "localhost:5000", "LISTEN", "listen address")
	dsn = stringFlag("dsn", "", "file:sws.db?cache=shared", "DSN", "database password")
	noMigrate = boolFlag("no-migrate", "m", false, "NOMIGRATE", "disable migrations")

	// Default to no log
	log = func(v ...interface{}) {}
	debug = func(v ...interface{}) {}

	tokenAuth = jwtauth.New("HS256", []byte("lkjasd0f9u203ijsldkfj"), nil)
}

type Renderer interface {
	Render(http.ResponseWriter, string, interface{}) error
}

func main() {
	flag.Parse()

	if *verbose {
		log = func(v ...interface{}) {
			fmt.Fprintf(os.Stdout, "[%s] ", time.Now().Format(time.RFC3339))
			fmt.Fprintln(os.Stdout, v...)
		}
	}
	if d := os.Getenv("DEBUG"); d != "" {
		debug = func(v ...interface{}) {
			fmt.Fprintf(os.Stdout, "[%s] ", time.Now().Format(time.RFC3339))
			fmt.Fprintln(os.Stdout, v...)
		}
	}
	log("version", Version)

	driver := strings.SplitN(*dsn, ":", 2)[0]
	if driver == "file" {
		driver = "sqlite3"
	}

	if noMigrate == nil || !*noMigrate {
		v, err := migrateDatabase(driver, *dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to migrate: %s", err)
			os.Exit(2)
		}
		log("database at version", v)
	}

	db, err := sqlx.Open(driver, *dsn)
	if err != nil {
		log(err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log(err)
		os.Exit(1)
	}
	var st sws.Store
	if driver == "pgx" {
		//st = store.P{db}
	} else {
		st = store.NewSqlite3Store(db)
	}

	tmplsCommon := []string{"flash.tmpl", "navbar.tmpl"}
	tmplsAuthed := append(tmplsCommon, []string{"layouts/base.tmpl", "charts.tmpl"}...)
	tmplsPublic := append(tmplsCommon, "layouts/public.tmpl")

	tmpls, err := LoadHTMLTemplateMap(map[string][]string{
		"sites":   append(tmplsAuthed, "sites.tmpl"),
		"site":    append(tmplsAuthed, "site.tmpl"),
		"home":    append(tmplsPublic, "home.tmpl"),
		"login":   append(tmplsPublic, "login.tmpl"),
		"example": []string{"example.tmpl"},
	}, funcMap)
	if err != nil {
		log(err)
		os.Exit(1)
	}
	//debug(tmpls.DefinedTemplates())
	renderer := templates.NewRenderer(tmpls)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestID)
	compressor := middleware.NewCompressor(5, "text/html", "text/css")
	r.Use(compressor.Handler())
	if *verbose {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	siteCtx := getSiteCtx(st)
	userCtx := getUserCtx(st)

	// For counter
	r.Get("/sws.js", handleCounter(*addr))
	r.Get("/sws.gif", handleHitCounter(st))

	// For UI
	r.Get("/hits", handleHits(st))

	// Authed routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(userCtx)
		r.Get(logoutURL, func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				//Secure: true,
				Expires: time.Time{},
			})
			r = flashSet(r, flashSuccess, "de-authenticated successfully")
			http.Redirect(w, r, flashURL(r, "/"), http.StatusSeeOther)
		})
		r.Route("/sites", func(r chi.Router) {
			r.Get("/", handleSites(st, renderer))
			r.Route("/{siteID}", func(r chi.Router) {
				r.Use(siteCtx)
				r.Get("/", handleSite(st, renderer))
				r.Route("/sparklines", func(r chi.Router) {
					r.Get("/{b:\\d+}-{e:\\d+}.svg", sparklineHandler(st))
				})
				r.Route("/charts", func(r chi.Router) {
					r.Get("/{b:\\d+}-{e:\\d+}.svg", svgChartHandler(st))
					r.Get("/{b:\\d+}-{e:\\d+}.png", svgChartHandler(st))
				})
			})
		})
	})

	// Example
	r.Get("/test.html", handleExample(renderer))

	authHandler := handleAuth(st, renderer)

	// Public routes
	r.Route("/", func(r chi.Router) {
		r.Get("/", handleIndex(renderer))
		r.Get(loginURL, func(w http.ResponseWriter, r *http.Request) {
			flash := flashGet(r)
			if err := renderer.Render(w, "login", map[string]interface{}{
				"Flash": flash,
			}); err != nil {
				httpError(w, 500, err.Error())
				return
			}
			return
		})
		r.Post(loginURL, authHandler)

		r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := strings.TrimPrefix(r.URL.Path, "/")
			if b, err := StaticLoadTemplate(p); err == nil {
				name := filepath.Base(p)
				http.ServeContent(w, r, name, time.Now(), bytes.NewReader(b))
			}
		}))
	})

	log("listening at", *addr)
	http.ListenAndServe(*addr, r)
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

func getUserCtx(db sws.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				log("missing claims")
				http.Redirect(w, r, flashURL(r, "/login"), http.StatusUnauthorized)
				return
			}

			id, ok := claims["user_id"]
			if !ok {
				log("missing user_id")
				http.Redirect(w, r, flashURL(r, "/login"), http.StatusUnauthorized)
				return
			}

			user, err := db.GetUserByID(int(id.(float64)))
			if err != nil {
				log("missing user")
				http.Redirect(w, r, flashURL(r, "/login"), http.StatusUnauthorized)
				return
			}
			debug("found user, adding to context")
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
