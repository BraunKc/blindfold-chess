package application

import (
	"context"
	"errors"

	"github.com/braunkc/blindfold-chess/internal/domain"
)

var (
	ErrNotFound = errors.New("game not found")
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
	rec, err := gs.gameRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrNotFound
	}

	return rec, nil
}

func (gs *GameService) ListGames(ctx context.Context) ([]domain.Record, error) {
	return gs.gameRepo.List(ctx)
}
