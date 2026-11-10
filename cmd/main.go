package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/braunkc/blindfold-chess/internal/application"
	"github.com/braunkc/blindfold-chess/internal/config"
	"github.com/braunkc/blindfold-chess/internal/infrastructure"
	httpapi "github.com/braunkc/blindfold-chess/internal/infrastructure/http"
	"github.com/braunkc/blindfold-chess/internal/infrastructure/persistence/game"
	"github.com/braunkc/blindfold-chess/internal/session"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	cfg, err := config.New("config/config.yml")
	if err != nil {
		slog.Error("create config", slog.Any("err", err))
		os.Exit(1)
	}

	session := session.New()

	db, err := infrastructure.Setup(cfg.DB)
	if err != nil {
		slog.Error("setup db", slog.Any("err", err))
		os.Exit(1)
	}

	if err = infrastructure.Migrate(&cfg.DB); err != nil {
		slog.Error("migrate", slog.Any("err", err))
		os.Exit(1)
	}

	gameRepo := game.New(db, cfg.DB)
	service := application.New(gameRepo)

	handler := httpapi.NewHandler(session, service)
	server := httpapi.NewServer(handler, &cfg.HTTP)

	if err := server.Start(); err != nil {
		log.Println(err)
	}
}
