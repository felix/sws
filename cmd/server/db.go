package main

import (
	"database/sql"
	"fmt"

	"src.userspace.com.au/go-migrate"
)

func migrateDatabase(driver, dsn string) (int64, error) {
	var version int64
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return version, err
	}
	defer db.Close()

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
