package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lilythecat859/fractal-rpc/internal/config"
	"github.com/lilythecat859/fractal-rpc/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	srv := server.New(cfg, logger)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		logger.Fatal("forced shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}
