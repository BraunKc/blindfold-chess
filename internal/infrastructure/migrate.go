package infrastructure

import (
	"errors"
	"fmt"

	"github.com/braunkc/blindfold-chess/internal/config"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *config.DB) error {
	migrationsPath, err := cfg.MigrationsPath()
	if err != nil {
		return fmt.Errorf("get migrations path: %w", err)
	}

	dsn, err := cfg.DSN()
	if err != nil {
		return fmt.Errorf("get DSN: %w", err)
	}

	m, err := migrate.New("file://"+migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("create migration: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
