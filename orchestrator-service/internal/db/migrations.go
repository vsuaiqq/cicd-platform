package db

import (
	"database/sql"
	"embed"

	sharedPostgres "github.com/vsuaiqq/cicd/shared/postgres"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func RunMigrations(db *sql.DB) error {
	return sharedPostgres.RunMigrations(db, migrationFS, "migrations")
}
