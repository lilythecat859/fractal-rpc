// SPDX-License-Identifier: AGPL-3.0-or-later
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/lilythecat859/solana-historical-rpc/internal/config"
	"github.com/lilythecat859/solana-historical-rpc/internal/rpc"
	"github.com/lilythecat859/solana-historical-rpc/internal/store"
	"go.uber.org/zap"
)

type Server struct {
	cfg   *config.Config
	log   *zap.Logger
	store store.Store
	srv   *http.Server
}

func New(cfg *config.Config, log *zap.Logger) (*Server, error) {
	st, err := store.NewClickHouse(
		cfg.ClickHouse.Addr,
		cfg.ClickHouse.Database,
		cfg.ClickHouse.Auth.Username,
		cfg.ClickHouse.Auth.Password,
		cfg.ClickHouse.Codec,
	)
	if err != nil {
		return nil, err
	}
	return &Server{cfg: cfg, log: log, store: st}, nil
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	h := rpc.NewHandler(s.store)
	h.Install(mux, s.cfg.RPCPath)

	s.srv = &http.Server{
		Addr:        fmt.Sprintf(":%d", s.cfg.HTTPPort),
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return ctx },
		ReadTimeout: 5 * time.Second,
	}
	s.log.Info("listening", zap.Int("port", s.cfg.HTTPPort))
	go func() {
		<-ctx.Done()
		s.log.Info("shutting down http")
		_ = s.srv.Shutdown(context.Background())
	}()
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
