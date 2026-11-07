package httpapi

type MakeMoveReq struct {
	Notation       string `json:"notation"`
	LastMovesCount int    `json:"last_moves_count"`
}

type SaveGameResp struct {
	ID        int64  `json:"id"`
	PGN       string `json:"pgn"`
	CreatedAt int64  `json:"created_at"`
}
