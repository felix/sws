package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
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
	"src.userspace.com.au/go-migrate"
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
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var st sws.Store
	if driver == "pgx" {
		//st = store.P{db}
	} else {
		st = store.NewSqlite3Store(db)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if *verbose {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	domainCtx := getDomainCtx(st)

	// For counter
	r.Get("/sws.js", handleCounter(*addr))
	r.Get("/sws.gif", handleHitCounter(st))

	// For UI
	r.Get("/hits", handleHits(st))
	r.Route("/domains", func(r chi.Router) {
		r.Get("/", handleDomains(st))
		r.Route("/{domainID}", func(r chi.Router) {
			r.Use(domainCtx)
			r.Get("/", handleDomain(st))
			r.Route("/sparklines", func(r chi.Router) {
				r.Get("/{s:\\d+}-{e:\\d+}.svg", sparklineHandler(st))
			})
			r.Route("/charts", func(r chi.Router) {
				r.Get("/{s:\\d+}-{e:\\d+}.svg", svgChartHandler(st))
				r.Get("/{s:\\d+}-{e:\\d+}.png", svgChartHandler(st))
			})
		})
	})
	r.Get("/", handleIndex())

	// Example
	r.Get("/test.html", handleExample())

	log("listening at", *addr)
	http.ListenAndServe(*addr, r)
}

func getDomainCtx(db sws.DomainStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(chi.URLParam(r, "domainID"))
			if err != nil {
				panic(err)
			}
			domain, err := db.GetDomainByID(id)
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			ctx := context.WithValue(r.Context(), "domain", domain)
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

func migrateDatabase(driver, dsn string) (int64, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	var version int64
	// Load migrations from generated file
	ms, err := decodeMigrations(driver)
	if err != nil {
		return 0, err
	}
	debug("found", len(ms), "migrations for driver", driver)
	migrator, err := migrate.NewStringMigrator(db, ms)
	if err != nil {
		return version, fmt.Errorf("failed to initialise: %w", err)
	}

	err = migrator.Migrate()
	if err != nil {
		return version, fmt.Errorf("failed to migrate: %w", err)
	}

	return migrator.Version()
}
