package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"time2meet/internal/infrastructure/config"
	"time2meet/internal/infrastructure/persistence/postgres"
	httpiface "time2meet/internal/presentation/http"
	"time2meet/pkg/logger"

	"go.uber.org/zap"
)

// @title Time2Meet API
// @version 0.1.0
// @description Система управления мероприятиями: события, площадки, билеты, отчёты, аудит.
// @BasePath /api/v1
// @schemes http
// @host localhost:8080
func main() {
	log := logger.New()
	defer func() { _ = log.Sync() }()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Error("config load failed", zap.Error(err))
		os.Exit(1)
	}

	db, err := postgres.NewDB(cfg.Database, log)
	if err != nil {
		log.Error("db connect failed", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	srv := httpiface.NewServer(cfg.HTTP, db, log)

	go func() {
		log.Info("http server starting", zap.String("addr", cfg.HTTP.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", zap.Error(err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown failed", zap.Error(err))
	}
}
