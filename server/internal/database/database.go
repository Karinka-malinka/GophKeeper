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

func NewDB(ps string) (*sql.DB, error) {

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("Error when creating the driver:" + err.Error())
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error initializing migrations:" + err.Error())
		return nil, err
	}

	defer m.Close()

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Debug("Error during migration:" + err.Error())
	}

	return db, nil
}
