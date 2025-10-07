package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lilythecat859/fractal-rpc/internal/config"
	"github.com/lilythecat859/fractal-rpc/internal/rpc"
	"github.com/lilythecat859/fractal-rpc/internal/store"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	st, err := store.NewClickHouse(
		cfg.ClickHouse.Addr,
		cfg.ClickHouse.Database,
		cfg.ClickHouse.Auth.Username,
		cfg.ClickHouse.Auth.Password,
		cfg.ClickHouse.Codec,
	)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}
	defer st.Close()

	mux := http.NewServeMux()
	h := rpc.NewHandler(st)
	path := cfg.RPCPath; if path == "" { path = "/" }; h.Install(mux, path)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
	}

	go func() {
		logger.Info("listening", zap.Int("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http serve", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutCtx)
}
