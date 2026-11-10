package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

var (
	ErrEmptyField   = errors.New("empty field")
	ErrInvalidField = errors.New("invalid field")
)

type Driver string

const (
	DriverSQLite   Driver = "sqlite"
	DriverPostgres Driver = "postgres"
)

type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (h *HTTP) validate() error {
	if h.Host == "" {
		return fmt.Errorf("host: %w", ErrEmptyField)
	}
	if h.Port <= 0 || h.Port > 65535 {
		return fmt.Errorf("port: %w", ErrInvalidField)
	}

	return nil
}

type SQLite struct {
	Path           string `yaml:"path"`
	MigrationsPath string `yaml:"migrations-path"`
}

func (s *SQLite) validate() error {
	if s.Path == "" {
		return fmt.Errorf("path: %w", ErrEmptyField)
	}
	if s.MigrationsPath == "" {
		return fmt.Errorf("migrations-path: %w", ErrEmptyField)
	}

	return nil
}

func (s *SQLite) DSN() string {
	return "sqlite://" + s.Path
}

type Postgres struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DB             string `yaml:"db"`
	MigrationsPath string `yaml:"migrations-path"`

	MaxOpenConns    int           `yaml:"max-open-conns"`
	MaxIdleConns    int           `yaml:"max-idle-conns"`
	ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime"`
}

func (p *Postgres) validate() error {
	if p.Host == "" {
		return fmt.Errorf("host: %w", ErrEmptyField)
	}
	if p.Port <= 0 || p.Port > 65535 {
		return fmt.Errorf("port: %w", ErrInvalidField)
	}
	if p.User == "" {
		return fmt.Errorf("user: %w", ErrEmptyField)
	}
	if p.DB == "" {
		return fmt.Errorf("db: %w", ErrEmptyField)
	}
	if p.MigrationsPath == "" {
		return fmt.Errorf("migrations-path: %w", ErrEmptyField)
	}

	if p.MaxOpenConns <= 0 {
		return fmt.Errorf("max-open-conns: %w", ErrInvalidField)
	}
	if p.MaxIdleConns < 0 {
		return fmt.Errorf("max-idle-conns: %w", ErrInvalidField)
	}
	if p.MaxIdleConns > p.MaxOpenConns {
		return fmt.Errorf("max-idle-conns must be <= max-open-conns: %w", ErrInvalidField)
	}
	if p.ConnMaxLifetime <= 0 {
		return fmt.Errorf("conn-max-lifetime: %w", ErrInvalidField)
	}

	return nil
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.DB)
}

type DB struct {
	Driver   Driver   `yaml:"driver"`
	SQLite   SQLite   `yaml:"sqlite"`
	Postgres Postgres `yaml:"postgres"`
}

func (d *DB) validate() error {
	switch d.Driver {
	case DriverSQLite:
		return d.SQLite.validate()
	case DriverPostgres:
		return d.Postgres.validate()
	default:
		return fmt.Errorf("driver: %w", ErrInvalidField)
	}
}

func (d *DB) DSN() (string, error) {
	switch d.Driver {
	case DriverSQLite:
		return d.SQLite.DSN(), nil
	case DriverPostgres:
		return d.Postgres.DSN(), nil
	default:
		return "", fmt.Errorf("driver: %w", ErrInvalidField)
	}
}

func (d *DB) MigrationsPath() (string, error) {
	switch d.Driver {
	case DriverSQLite:
		return d.SQLite.MigrationsPath, nil
	case DriverPostgres:
		return d.Postgres.MigrationsPath, nil
	default:
		return "", fmt.Errorf("driver: %w", ErrInvalidField)
	}
}

type Config struct {
	HTTP HTTP `yaml:"http"`
	DB   DB   `yaml:"db"`
}

func (c *Config) validate() error {
	if err := c.HTTP.validate(); err != nil {
		return fmt.Errorf("validate http: %w", err)
	}
	if err := c.DB.validate(); err != nil {
		return fmt.Errorf("validate db: %w", err)
	}

	return nil
}

func New(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("empty config path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	data = []byte(os.ExpandEnv(string(data)))

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	if err = cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config %q: %w", path, err)
	}

	return &cfg, nil
}
