package main

import (
	"log"
	"os"

	"github.com/lilythecat859/fractal-rpc/internal/config"
	"github.com/lilythecat859/fractal-rpc/internal/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.MustLoad()
	if err := server.Run(cfg, logger); err != nil {
		logger.Fatal("server died", zap.Error(err))
	}
}
