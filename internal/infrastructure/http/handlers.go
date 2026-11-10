package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/braunkc/blindfold-chess/internal/application"
	"github.com/braunkc/blindfold-chess/internal/domain"
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if err := h.session.Move(req.Notation); err != nil {
			handleError(ctx, err)
			return
		}

		n := req.LastMovesCount
		if n <= 0 {
			n = 3
		}

		ctx.JSON(http.StatusOK, h.session.State(n))
	}
}

func (h *Handler) Reset() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.session.Reset()
		ctx.JSON(http.StatusOK, h.session.State(0))
	}
}

func (h *Handler) SaveGame() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		record, err := h.service.SaveGame(ctx.Request.Context(), h.session.PGN())
		if err != nil {
			handleError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, SaveGameResp{
			ID:        record.ID,
			PGN:       record.PGN,
			CreatedAt: record.CreatedAt.Unix(),
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
		ctx.JSON(http.StatusOK, gin.H{
			"last_moves": h.session.LastMoves(parseLMC(ctx, 3)),
		})
	}
}

func (h *Handler) Outcome() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"outcome": h.session.Outcome(),
		})
	}
}

func (h *Handler) IsOver() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"is_over": h.session.IsOver(),
		})
	}
}

func (h *Handler) State() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, h.session.State(parseLMC(ctx, 3)))
	}
}

func (h *Handler) IsPromotion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"is_promotion": h.session.IsPromotion(ctx.Query("notation"))})
	}
}

func parseLMC(ctx *gin.Context, def int) int {
	n, err := strconv.Atoi(ctx.Query("lmc"))
	if err != nil || n <= 0 || n > 100 {
		n = def
	}

	return n
}

func handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, session.ErrInvalidMove):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid move"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		slog.Error("internal error", slog.Any("err", err))
	}
}
