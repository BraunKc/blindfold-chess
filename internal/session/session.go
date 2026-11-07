package session

import (
	"fmt"
	"sync"

	"github.com/corentings/chess/v2"
)

type Session struct {
	mu   sync.Mutex
	game *chess.Game
}

type State struct {
	Turn      string   `json:"turn"`
	Outcome   string   `json:"outcome"`
	IsOver    bool     `json:"is_over"`
	PGN       string   `json:"pgn"`
	LastMoves []string `json:"last_moves"`
}

func New() *Session {
	return &Session{
		game: chess.NewGame(),
	}
}

func (s *Session) Move(notation string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.game.PushNotationMove(notation, chess.UCINotation{}, nil); err != nil {
		return fmt.Errorf("invalid move %s: %w", notation, err)
	}

	return nil
}

func (s *Session) PGN() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.game.String()
}

func (s *Session) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.game = chess.NewGame()
}

func (s *Session) Turn() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.game.CurrentPosition().Turn().String()
}

func (s *Session) LastMoves(n int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	moves := s.game.Moves()
	if len(moves) > n {
		moves = moves[len(moves)-n:]
	}

	lastMoves := make([]string, len(moves))
	for i, m := range moves {
		lastMoves[i] = m.String()
	}

	return lastMoves
}

func (s *Session) Outcome() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.game.Outcome().String()
}

func (s *Session) IsOver() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.game.Outcome() != chess.NoOutcome
}

func (s *Session) State(lmc int) State {
	s.mu.Lock()
	defer s.mu.Unlock()

	outcome := s.game.Outcome().String()
	moves := s.game.Moves()
	if len(moves) > lmc {
		moves = moves[len(moves)-lmc:]
	}

	lastMoves := make([]string, len(moves))
	for i, m := range moves {
		lastMoves[i] = m.String()
	}

	return State{
		Turn:      s.game.CurrentPosition().Turn().String(),
		Outcome:   outcome,
		IsOver:    outcome != chess.NoOutcome.String(),
		PGN:       s.game.String(),
		LastMoves: lastMoves,
	}
}
