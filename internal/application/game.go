package application

import (
	"context"

	"github.com/braunkc/blindfold-chess/internal/domain"
)

type GameRepository interface {
	Save(ctx context.Context, g *domain.Record) (*domain.Record, error)
	GetByID(ctx context.Context, id int64) (*domain.Record, error)
	List(ctx context.Context) ([]domain.Record, error)
}

type GameService struct {
	gameRepo GameRepository
}

func New(gameRepo GameRepository) *GameService {
	return &GameService{
		gameRepo: gameRepo,
	}
}
