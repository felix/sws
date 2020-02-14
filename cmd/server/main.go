package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"src.userspace.com.au/sws"
	"src.userspace.com.au/sws/store"
)

// Flags
var (
	verbose   *bool
	addr      *string
	dsn       *string
	noMigrate *bool
)

var log, debug sws.Logger

func init() {
	verbose = boolFlag("verbose", "v", false, "VERBOSE", "enable verbose output")
	addr = stringFlag("listen", "l", "localhost:5000", "LISTEN", "listen address")
	dsn = stringFlag("dsn", "", "file:sws.db?cache=shared", "DSN", "database password")
	noMigrate = boolFlag("no-migrate", "m", false, "NOMIGRATE", "disable migrations")

	// Default to no log
	log = func(v ...interface{}) {}
	debug = func(v ...interface{}) {}
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

	driver := strings.SplitN(*dsn, ":", 2)[0]
	if driver == "file" {
		driver = "sqlite3"
	}

	if noMigrate == nil || !*noMigrate {
		version, err := migrateDatabase(driver, *dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to migrate: %s", err)
			os.Exit(2)
		}
		log("database at version", version)
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

	tmpls := template.Must(loadTemplateHTML([]string{
		"home",
		"sites",
		"example",
		"partials/navMain",
		"partials/pageHead",
		"partials/pageFoot",
		"partials/siteForList",
	}, funcMap))
	debug(tmpls.DefinedTemplates())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if *verbose {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	siteCtx := getSiteCtx(st)

	// For counter
	r.Get("/sws.js", handleCounter(*addr))
	r.Get("/sws.gif", handleHitCounter(st))

	// For UI
	r.Get("/hits", handleHits(st))
	r.Route("/sites", func(r chi.Router) {
		r.Get("/", handleSites(st, tmpls))
		r.Route("/{siteID}", func(r chi.Router) {
			r.Use(siteCtx)
			r.Get("/", handleSite(st))
			r.Route("/sparklines", func(r chi.Router) {
				r.Get("/{s:\\d+}-{e:\\d+}.svg", sparklineHandler(st))
			})
			r.Route("/charts", func(r chi.Router) {
				r.Get("/{s:\\d+}-{e:\\d+}.svg", svgChartHandler(st))
				r.Get("/{s:\\d+}-{e:\\d+}.png", svgChartHandler(st))
			})
		})
	})
	r.Get("/", handleIndex(tmpls))

	// Example
	r.Get("/test.html", handleExample(tmpls))

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

// func getAuthCtx(db sws.UserStore) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := context.WithValue(r.Context(), "user", user)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }
