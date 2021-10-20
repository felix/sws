package store

import (
	"database/sql"
	"embed"
	"fmt"

	"src.userspace.com.au/migrate"
)

//go:embed sql/*
var ms embed.FS

func Migrate(driver, dsn string) (int, error) {
	var version int
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return version, err
	}
	defer db.Close()

	//debug("found", len(ms), "migrations for driver", driver)
	migrator, err := migrate.NewFSMigrator(db, ms)
	if err != nil {
		return version, fmt.Errorf("failed to initialise: %w", err)
	}

	err = migrator.Migrate()
	if err != nil {
		return version, fmt.Errorf("failed to migrate: %w", err)
	}

	return migrator.Version()
}
