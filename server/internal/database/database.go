package database

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(ps string) (*sql.DB, *migrate.Migrate, error) {

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Error when creating the driver:" + err.Error())
		return nil, nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error initializing migrations:" + err.Error())
		return nil, nil, err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Debug("No changes during migration")
	} else if err != nil {
		slog.Error("Error during migration: " + err.Error())
		return nil, nil, err
	}

	return db, m, nil
}
