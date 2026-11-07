package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/braunkc/blindfold-chess/internal/domain"
)

type GameRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

func (gr *GameRepository) Save(ctx context.Context, g *domain.Record) (*domain.Record, error) {
	query, args, err := squirrel.Insert("games").
		Columns("pgn", "created_at").
		Values(g.PGN, g.CreatedAt.Unix()).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	res, err := gr.db.ExecContext(ctx, query, args...)
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return gr.GetByID(ctx, id)
}

func (gr *GameRepository) GetByID(ctx context.Context, id int64) (*domain.Record, error) {
	query, args, err := squirrel.Select("id", "pgn", "created_at").
		From("games").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	row := gr.db.QueryRowContext(ctx, query, args...)

	return gr.parseRow(row)
}

func (gr *GameRepository) List(ctx context.Context) ([]domain.Record, error) {
	query, args, err := squirrel.Select("id", "pgn", "created_at").
		From("games").
		OrderBy("created_at DESC").
		Limit(50).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	rows, err := gr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	return gr.parseRows(rows)
}

func (gr *GameRepository) parseRow(row *sql.Row) (*domain.Record, error) {
	var g domain.Record
	var createdAt int64
	if err := row.Scan(&g.ID, &g.PGN, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to parse row: %w", err)
	}

	g.CreatedAt = time.Unix(createdAt, 0)

	return &g, nil
}

func (gr *GameRepository) parseRows(rows *sql.Rows) ([]domain.Record, error) {
	defer rows.Close()

	var records []domain.Record
	for rows.Next() {
		var g domain.Record
		var createdAt int64
		if err := rows.Scan(&g.ID, &g.PGN, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to parse rows: %w", err)
		}

		g.CreatedAt = time.Unix(createdAt, 0)
		records = append(records, g)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to do query: %w", rows.Err())
	}

	return records, nil
}
