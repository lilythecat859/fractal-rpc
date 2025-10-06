// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lilythecat859/solana-historical-rpc/internal/config"
	"github.com/lilythecat859/solana-historical-rpc/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	log := zap.Must(zap.NewProduction()).Named("fractal")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Fatal("start server", zap.Error(err))
	}
	if err := srv.Start(ctx); err != nil {
		log.Fatal("run server", zap.Error(err))
	}
	log.Info("graceful shutdown complete")
}
