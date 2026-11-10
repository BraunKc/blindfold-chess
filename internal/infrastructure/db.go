package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/braunkc/blindfold-chess/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

func setupSQLite(cfg config.SQLite) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("connect to db %q: %w", cfg.Path, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

func setupPostgres(cfg config.Postgres) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect to db %q: %w", cfg.DSN(), err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

func Setup(cfg config.DB) (*sql.DB, error) {
	switch cfg.Driver {
	case config.DriverSQLite:
		return setupSQLite(cfg.SQLite)
	case config.DriverPostgres:
		return setupPostgres(cfg.Postgres)
	default:
		return nil, fmt.Errorf("unknown driver: %q", cfg.Driver)
	}
}
