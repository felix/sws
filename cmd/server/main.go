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
	"src.userspace.com.au/flags"
	"src.userspace.com.au/go-migrate"
	"src.userspace.com.au/sws"
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
	verbose = flags.Bool("verbose", "v", false, "VERBOSE", "enable verbose output")
	addr = flags.String("listen", "l", "localhost:5000", "LISTEN", "listen address")
	dsn = flags.String("dsn", "", "file:sws.db?cache=shared", "DSN", "database password")
	noMigrate = flags.Bool("no-migrate", "m", false, "NOMIGRATE", "disable migrations")

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

	if noMigrate == nil || !*noMigrate {
		version, err := migrateDatabase(db, driver)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to migrate: %s", err)
			os.Exit(2)
		}
		log("database at version", version)
	}

	r := chi.NewRouter()

	r.Get("/sws.js", handleCounter(*addr))
	r.Get("/sws.gif", handleHitCounter(db))
	r.Get("/hits", handleHits(db))
	r.Get("/domains", handleDomains(db))
	r.Get("/test.html", handleTest())
	r.Get("/", handleIndex())

	log("listening at", *addr)
	http.ListenAndServe(*addr, r)
}

func migrateDatabase(db *sql.DB, driver string) (int64, error) {
	var v int64
	// Load migrations
	ms, err := decodeMigrations(driver)
	if err != nil {
		return 0, err
	}
	debug("found", len(ms), "migrations for driver", driver)
	migrator, err := migrate.NewStringMigrator(db, ms)
	if err != nil {
		return v, fmt.Errorf("failed to initialise: %w", err)
	}

	err = migrator.Migrate()
	if err != nil {
		return v, fmt.Errorf("failed to migrate: %w", err)
	}

	v, err = migrator.Version()
	return v, nil
}
