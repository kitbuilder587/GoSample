package migrations

import (
	"embed"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationFiles embed.FS

func newMigrator(db *sqlx.DB) (*migrate.Migrate, error) {
	src, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return nil, errors.Newf("iofs source: %w", err)
	}

	dbDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, errors.Newf("postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", dbDriver)
	if err != nil {
		return nil, errors.Newf("construct migrator: %w", err)
	}
	return m, nil
}

func Up(db *sqlx.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Newf("apply migrations: %w", err)
	}
	return nil
}

func Run(db *sqlx.DB) error {
	return Up(db)
}

func Down(db *sqlx.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Newf("revert migrations: %w", err)
	}
	return nil
}
