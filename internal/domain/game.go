package domain

import (
	"time"
)

type Record struct {
	ID        int64
	PGN       string
	CreatedAt time.Time
}

func NewRecord(pgn string) *Record {
	return &Record{
		PGN:       pgn,
		CreatedAt: time.Now().UTC(),
	}
}
