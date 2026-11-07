package httpapi

type MakeMoveReq struct {
	Notation string `json:"notation"`
}

type GetPGNResp struct {
	PGN string `json:"pgn"`
}

type SaveGameResp struct {
	ID        int64  `json:"id"`
	PGN       string `json:"pgn"`
	CreatedAt int64  `json:"created_at"`
}
