package main

import (
	"github.com/lilythecat859/fractal-rpc/internal/config"
	"github.com/lilythecat859/fractal-rpc/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	if err := server.Run(cfg, logger); err != nil {
		logger.Fatal("server exited", zap.Error(err))
	}
}
