package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	detector "github.com/mssola/user_agent"

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
	dsn       string
	domain    string
	logFile   string
	maxmind   string
	noMigrate bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "enable verbose output")
	flag.StringVar(&dsn, "dsn", "file:sws.db?cache=shared", "database password")
	flag.StringVar(&domain, "domain", "stats.userspace.com.au", "stats domain")
	flag.StringVar(&logFile, "log", "", "log to file")
	flag.StringVar(&maxmind, "maxmind", "", "maxmind country DB path")
	flag.BoolVar(&noMigrate, "no-migrate", false, "disable migrations")

	// Default to no log
	log = func(v ...interface{}) {}
	debug = func(v ...interface{}) {}
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

	cache, err := lru.New(100)
	if err != nil {
		panic(err)
	}

	seenUAs := make(map[string]bool)

	toUpdate := make([]*sws.Hit, 0)
	log("updating country code")
	err = st.HitCursor(func(h *sws.Hit) error {
		// Populate missing country codes
		if maxmind != "" && h.CountryCode == nil {
			var cc *string
			if v, ok := cache.Get(h.Addr); ok {
				cc = v.(*string)
			} else {
				if cc, err = sws.FetchCountryCode(maxmind, h.Addr); err != nil {
					log("geoip lookup failed:", err)
				}
				cache.Add(h.Addr, cc)
			}
			h.CountryCode = cc
			toUpdate = append(toUpdate, h)
		}

		// Populate user agent
		if h.UserAgent != nil {
			log("hit ID", *h.ID)
			ua := h.UserAgent
			if ok := seenUAs[ua.Hash]; ok {
				return nil
			}
			det := detector.New(ua.Name)
			browser, version := det.Browser()

			h.UserAgent.Browser = browser
			h.UserAgent.Platform = det.Platform()
			h.UserAgent.Version = version
			h.UserAgent.Bot = det.Bot()
			h.UserAgent.Mobile = strings.Contains(ua.Name, "Mobi") || det.Mobile()

			seenUAs[ua.Hash] = true

			toUpdate = append(toUpdate, h)
		}
		return nil
	})
	if err != nil {
		log(err)
		os.Exit(1)
	}
	if len(toUpdate) > 0 {
		for _, h := range toUpdate {
			if err := st.SaveHit(h); err != nil {
				log(err)
			}
		}
	}
}
