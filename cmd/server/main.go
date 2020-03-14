package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"src.userspace.com.au/sws"
	"src.userspace.com.au/sws/store"
)

var (
	Version    string
	log, debug sws.Logger
)

// Flags
var (
	verbose   bool
	addr      string
	dsn       string
	domain    string
	logFile   string
	override  string
	maxmind   string
	noMigrate bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "enable verbose output")
	flag.StringVar(&addr, "listen", "localhost:5000", "listen address")
	flag.StringVar(&dsn, "dsn", "file:sws.db?cache=shared", "database password")
	flag.StringVar(&domain, "domain", "stats.userspace.com.au", "stats domain")
	flag.StringVar(&logFile, "l", "", "log to file")
	flag.StringVar(&override, "override", "", "override path")
	flag.StringVar(&maxmind, "maxmind", "", "maxmind country DB path")
	flag.BoolVar(&noMigrate, "no-migrate", false, "disable migrations")

	// Default to no log
	log = func(v ...interface{}) {}
	debug = func(v ...interface{}) {}

}

type Renderer interface {
	Render(http.ResponseWriter, string, interface{}) error
}

func main() {
	var err error
	flag.Parse()

	var output io.Writer = os.Stdout
	if logFile != "" {
		if output, err = os.Create(logFile); err != nil {
			fmt.Fprintf(os.Stderr, "failed to open log file: %s", err)
			os.Exit(1)
		}
	}

	if verbose {
		log = func(v ...interface{}) {
			fmt.Fprintf(output, "[%s] ", time.Now().Format(time.RFC3339))
			fmt.Fprintln(output, v...)
		}
	}
	if d := os.Getenv("DEBUG"); d != "" {
		debug = func(v ...interface{}) {
			fmt.Fprintf(output, "[%s] ", time.Now().Format(time.RFC3339))
			fmt.Fprintln(output, v...)
		}
	}
	log("version", Version)

	driver := strings.SplitN(dsn, ":", 2)[0]
	if driver == "file" {
		driver = "sqlite3"
	}

	if !noMigrate {
		v, err := migrateDatabase(driver, dsn)
		if err != nil {
			log("failed to migrate:", err)
			os.Exit(2)
		}
		log("database at version", v)
	}

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log("failed to open database:", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log("failed to connect to database:", err)
		os.Exit(1)
	}
	var st sws.Store
	if driver == "pgx" {
		//st = store.P{db}
	} else {
		st = store.NewSqlite3Store(db)
	}

	r, err := createRouter(st, maxmind)
	if err != nil {
		log("failed to create router:", err)
		os.Exit(1)
	}

	log("listening at", addr)
	http.ListenAndServe(addr, r)
}
