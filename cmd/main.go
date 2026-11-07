package main

import (
	"log"
	"os"

	"github.com/braunkc/blindfold-chess/internal/application"
	"github.com/braunkc/blindfold-chess/internal/config"
	"github.com/braunkc/blindfold-chess/internal/infrastructure"
	httpapi "github.com/braunkc/blindfold-chess/internal/infrastructure/http"
	"github.com/braunkc/blindfold-chess/internal/infrastructure/persistence/game"
	"github.com/braunkc/blindfold-chess/internal/session"
)

func main() {
	cfg, err := config.Setup("config/config.yml")
	if err != nil {
		log.Println("failed to setup config:", err)
		os.Exit(1)
	}

	session := session.New()

	db, err := infrastructure.Setup(cfg.DB)
	if err != nil {
		log.Println("failed to setup db:", err)
		os.Exit(1)
	}
	gameRepo := game.New(db)
	service := application.New(gameRepo)

	handler := httpapi.NewHandler(session, service)
	server := httpapi.NewServer(handler, &cfg.HTTP)

	if err := server.Start(); err != nil {
		log.Println(err)
	}
}
