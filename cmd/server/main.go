package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/stdlib"
	"src.userspace.com.au/go-migrate"
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

	if err := upgradeDB(db); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		//AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/sws.js", handleSnippet())
	r.Get("/sws.gif", handlePageView(db))
	r.Get("/", handleIndex())
	http.ListenAndServe(":5000", r)
}

func upgradeDB(db *sql.DB) error {
	// Relative path to migration files
	migrator, err := migrate.NewFileMigrator(db, "file://migrations/")
	if err != nil {
		return fmt.Errorf("failed to create migrator: %s", err)
	}

	// Migrate all the way
	err = migrator.Migrate()
	if err != nil {
		return fmt.Errorf("failed to migrate: %s", err)
	}

	v, err := migrator.Version()
	fmt.Printf("database at version %d\n", v)
	return err
}

// Get envvar string with default
func envString(v, d string) string {
	out := os.Getenv(v)
	if out == "" {
		out = d
	}
	return out
}
