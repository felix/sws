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
	verbose   *bool
	addr      *string
	dsn       *string
	domain    *string
	logFile   *string
	noMigrate *bool
)

func init() {
	verbose = boolFlag("verbose", "v", false, "VERBOSE", "enable verbose output")
	addr = stringFlag("listen", "l", "localhost:5000", "LISTEN", "listen address")
	dsn = stringFlag("dsn", "", "file:sws.db?cache=shared", "DSN", "database password")
	domain = stringFlag("domain", "", "stats.userspace.com.au", "DOMAIN", "stats domain")
	logFile = stringFlag("log", "", "", "LOGFILE", "log to file")
	noMigrate = boolFlag("no-migrate", "m", false, "NOMIGRATE", "disable migrations")

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
	if logFile != nil && *logFile != "" {
		if output, err = os.Create(*logFile); err != nil {
			fmt.Fprintf(os.Stderr, "failed to open log file: %s", err)
			os.Exit(1)
		}
	}

	if *verbose {
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

	driver := strings.SplitN(*dsn, ":", 2)[0]
	if driver == "file" {
		driver = "sqlite3"
	}

	if noMigrate == nil || !*noMigrate {
		v, err := migrateDatabase(driver, *dsn)
		if err != nil {
			log("failed to migrate:", err)
			os.Exit(2)
		}
		log("database at version", v)
	}

	db, err := sqlx.Open(driver, *dsn)
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

	r, err := createRouter(st)
	if err != nil {
		log("failed to create router:", err)
		os.Exit(1)
	}

	log("listening at", *addr)
	http.ListenAndServe(*addr, r)
}
