package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"postgres://%s:%s@%s:5432/%s?sslmode=disable",
			os.Getenv("PG_USER"),
			os.Getenv("PG_PASS"),
			envString("PG_HOST", "localhost"),
			envString("PG_DB", "swa"),
		),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Get("/page_views", handlePageViews(db))
	r.Get("/domains", handleDomains(db))
	r.Get("/", handleIndex())
	http.ListenAndServe(":5000", r)
}

// Get envvar string with default
func envString(v, d string) string {
	out := os.Getenv(v)
	if out == "" {
		out = d
	}
	return out
}
