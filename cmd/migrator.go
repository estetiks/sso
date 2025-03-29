package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func main() {

	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storagePath", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrationsPath", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrationsTable", "migrations", "name of migrations table")

	flag.Parse()

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)

	if err != nil {
		if errors.Is(err, migrate.Err) {
			fmt.Println("")
		}
	}
}
