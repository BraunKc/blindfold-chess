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

func (gs *GameService) SaveGame(ctx context.Context, pgn string) (*domain.Record, error) {
	g := domain.NewRecord(pgn)
	return gs.gameRepo.Save(ctx, g)
}

func (gs *GameService) GetGameByID(ctx context.Context, id int64) (*domain.Record, error) {
	return gs.gameRepo.GetByID(ctx, id)
}

func (gs *GameService) ListGames(ctx context.Context) ([]domain.Record, error) {
	return gs.gameRepo.List(ctx)
}
