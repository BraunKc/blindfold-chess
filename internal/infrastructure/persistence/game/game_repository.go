package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/braunkc/blindfold-chess/internal/config"
	"github.com/braunkc/blindfold-chess/internal/domain"
)

type GameRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func New(db *sql.DB, cfg config.DB) *GameRepository {
	var builder squirrel.StatementBuilderType
	switch cfg.Driver {
	case config.DriverSQLite:
		builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	case config.DriverPostgres:
		builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	}

	return &GameRepository{
		db:      db,
		builder: builder,
	}
}

func (gr *GameRepository) Save(ctx context.Context, g *domain.Record) (*domain.Record, error) {
	query, args, err := gr.builder.Insert("games").
		Columns("pgn", "created_at").
		Values(g.PGN, g.CreatedAt.Unix()).
		Suffix("RETURNING id, pgn, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("create query: %w", err)
	}

	var r domain.Record
	var createdAt int64
	if err = gr.db.QueryRowContext(ctx, query, args...).Scan(
		&r.ID, &r.PGN, &createdAt,
	); err != nil {
		return nil, fmt.Errorf("query %q: %w", query, err)
	}
	r.CreatedAt = time.Unix(createdAt, 0)

	return &r, nil
}

func (gr *GameRepository) GetByID(ctx context.Context, id int64) (*domain.Record, error) {
	query, args, err := gr.builder.Select("id", "pgn", "created_at").
		From("games").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("create query %q: %w", query, err)
	}

	row := gr.db.QueryRowContext(ctx, query, args...)

	return gr.parseRow(row)
}

func (gr *GameRepository) List(ctx context.Context) ([]domain.Record, error) {
	query, args, err := gr.builder.Select("id", "pgn", "created_at").
		From("games").
		OrderBy("created_at DESC").
		Limit(100).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("create query %q: %w", query, err)
	}

	rows, err := gr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query %q: %w", query, err)
	}

	return gr.parseRows(rows)
}

func (gr *GameRepository) parseRow(row *sql.Row) (*domain.Record, error) {
	var r domain.Record
	var createdAt int64
	if err := row.Scan(&r.ID, &r.PGN, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("parse row: %w", err)
	}
	r.CreatedAt = time.Unix(createdAt, 0)

	return &r, nil
}

func (gr *GameRepository) parseRows(rows *sql.Rows) ([]domain.Record, error) {
	defer rows.Close()

	records := make([]domain.Record, 0)
	for rows.Next() {
		var r domain.Record
		var createdAt int64
		if err := rows.Scan(&r.ID, &r.PGN, &createdAt); err != nil {
			return nil, fmt.Errorf("parse rows: %w", err)
		}
		r.CreatedAt = time.Unix(createdAt, 0)
		records = append(records, r)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("do query: %w", rows.Err())
	}

	return records, nil
}
