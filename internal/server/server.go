
package server

import (
	"context"
	"fmt"
	"net"
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

type Server struct {
	httpSrv *http.Server
	logger  *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Server {
	st, err := store.NewClickHouse(
		cfg.ClickHouse.Addr,
		cfg.ClickHouse.Database,
		cfg.ClickHouse.Auth.Username,
		cfg.ClickHouse.Auth.Password,
		cfg.ClickHouse.Codec,
	)
	if err != nil {
		logger.Fatal("open store", zap.Error(err))
	}

	mux := http.NewServeMux()
	h := rpc.NewHandler(st)
	path := cfg.RPCPath
	if path == "" {
		path = "/"
	}
	h.Install(mux, path)

	return &Server{
		httpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
			Handler: mux,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpSrv.Addr)
	if err != nil {
		return err
	}
	s.logger.Info("listening", zap.String("addr", ln.Addr().String()))
	return s.httpSrv.Serve(ln)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

func Run(cfg *config.Config, logger *zap.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s := New(cfg, logger)

	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server start", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Stop(shutCtx)
}
