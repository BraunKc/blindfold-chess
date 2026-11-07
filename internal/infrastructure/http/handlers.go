package httpapi

import (
	"net/http"
	"strconv"

	"github.com/braunkc/blindfold-chess/internal/application"
	"github.com/braunkc/blindfold-chess/internal/session"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	session *session.Session
	service *application.GameService
}

func NewHandler(session *session.Session, service *application.GameService) *Handler {
	return &Handler{
		session: session,
		service: service,
	}
}

func (h *Handler) RenderMain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "main.html", nil)
	}
}

func (h *Handler) MakeMove() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req MakeMoveReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
			return
		}

		if err := h.session.Move(req.Notation); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"turn": h.session.Turn(),
			"pgn":  h.session.PGN(),
		})
	}
}

func (h *Handler) Reset() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.session.Reset()
		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) SaveGame() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		record, err := h.service.SaveGame(ctx.Request.Context(), h.session.PGN())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save game"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"game": SaveGameResp{
				ID:        record.ID,
				PGN:       record.PGN,
				CreatedAt: record.CreatedAt.Unix(),
			},
		})
	}
}

func (h *Handler) GetPGN() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"pgn": h.session.PGN(),
		})
	}
}

func (h *Handler) GetTurn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"turn": h.session.Turn(),
		})
	}
}

func (h *Handler) LastMoves() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		n := ctx.Query("n")
		nInt, err := strconv.Atoi(n)
		if err != nil {
			nInt = 3
		}

		ctx.JSON(http.StatusOK, gin.H{
			"last_moves": h.session.LastMoves(nInt),
		})
	}
}
