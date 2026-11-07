package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DB struct {
	Path string `yaml:"path"`
}

type Config struct {
	HTTP HTTP `yaml:"http"`
	DB   DB   `yaml:"db"`
}

func Setup(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("empty config path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	var cfg Config
	yaml.Unmarshal(data, &cfg)

	return &cfg, nil
}
