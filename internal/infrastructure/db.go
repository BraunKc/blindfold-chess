package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/braunkc/blindfold-chess/internal/config"
	_ "modernc.org/sqlite"
)

func Setup(cfg config.DB) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
