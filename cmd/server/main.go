package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
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

	// For counter
	r.Get("/sws.js", handleCounter(*addr))
	r.Get("/sws.gif", handleHitCounter(st))

	// For UI
	r.Get("/hits", handleHits(st))
	r.Get("/domains", handleDomains(st))
	r.Get("/", handleIndex())

	// Example
	r.Get("/test.html", handleExample())

	log("listening at", *addr)
	http.ListenAndServe(*addr, r)
}

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
