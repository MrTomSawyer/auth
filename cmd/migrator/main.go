package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable, migrationType string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&migrationType, "migration-type", "", "up or down migration")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)

	if err != nil {
		panic(err)
	}

	switch migrationType {
	case "up":
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no up migrations to apply")
				return
			}
			panic(err)
		}
	case "down":
		if err = m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no down migrations to apply")
				return
			}
			panic(err)
		}
	default:
		panic("migration type is not specified")
	}
}

// go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations --migration-type=up
// go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test
