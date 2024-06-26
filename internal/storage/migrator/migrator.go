package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type DatabaseDriver interface {
	Open(db *sql.DB) (database.Driver, error)
}

type Migrator struct {
	srcDriver source.Driver
	dbDriver  DatabaseDriver
}

func NewMigrator(sqlFiles embed.FS, dirName string, dbDriver DatabaseDriver) (*Migrator, error) {
	const op = "/internal/storage/migrator.NewMigrator"

	d, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to create new iofs: %w", op, err)
	}
	return &Migrator{
		srcDriver: d,
		dbDriver:  dbDriver,
	}, nil
}

func (m *Migrator) ApplyMigrations(db *sql.DB, dbName string) error {
	const op = "/internal/storage/migrator.ApplyMigrations"

	driver, err := m.dbDriver.Open(db)
	if err != nil {
		return fmt.Errorf("%s: unable to open database driver: %v", op, err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, dbName, driver)
	if err != nil {
		return fmt.Errorf("%s: unable to create migration: %v", op, err)
	}

	defer func() {
		if err, err2 := migrator.Close(); err != nil || err2 != nil {
			log.Printf("%s: error closing migrator: %v", op, err)
		}
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: unable to apply migrations: %v", op, err)
	}

	return nil
}
