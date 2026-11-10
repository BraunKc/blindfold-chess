package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/braunkc/blindfold-chess/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	http *http.Server
}

func NewServer(h *Handler, cfg *config.HTTP) *Server {
	r := gin.New()
	r.LoadHTMLGlob("./web/templates/*")
	r.Static("/css", "./web/static/css")
	r.Static("/js", "./web/static/js")

	setupRoutes(r, h)

	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: r,
		},
	}
}

func setupRoutes(r *gin.Engine, h *Handler) {
	r.GET("/", h.RenderMain())

	api := r.Group("/api")
	{
		api.POST("/move", h.MakeMove())
		api.POST("/reset", h.Reset())
		api.POST("/save", h.SaveGame())
		api.GET("/pgn", h.GetPGN())
		api.GET("/turn", h.GetTurn())
		api.GET("/last-moves", h.LastMoves())
		api.GET("/outcome", h.Outcome())
		api.GET("/is-over", h.IsOver())
		api.GET("/state", h.State())
		api.GET("/is-promotion", h.IsPromotion())
	}
}

func (s *Server) Start() error {
	if err := s.http.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
