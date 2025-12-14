package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"time2meet/internal/infrastructure/config"
	"time2meet/internal/presentation/http/middleware"

	"go.uber.org/zap"
)

func NewServer(cfg config.HTTPConfig, db *sqlx.DB, log *zap.Logger) *http.Server {
	r := NewRouter(Dependencies{DB: db, Log: log})
	r.Use(middleware.ContextFromHeaders())

	r.GET("/healthz", func(c *gin.Context) {
		type resp struct {
			Status string `json:"status"`
		}
		if db != nil {
			if err := db.Ping(); err != nil {
				log.Warn("healthz db ping failed", zap.Error(err))
				c.JSON(http.StatusServiceUnavailable, resp{Status: "db_down"})
				return
			}
		}
		c.JSON(http.StatusOK, resp{Status: "ok"})
	})

	s := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}
