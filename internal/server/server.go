package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lilythecat859/fractal-rpc/internal/auth"
	"github.com/lilythecat859/fractal-rpc/internal/config"
	"github.com/lilythecat859/fractal-rpc/internal/rpc"
	"github.com/lilythecat859/fractal-rpc/internal/store"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	st, err := store.NewClickHouse(
		cfg.ClickHouse.Addr,
		cfg.ClickHouse.Database,
		cfg.ClickHouse.Auth.Username,
		cfg.ClickHouse.Auth.Password,
	)
	if err != nil {
		return err
	}
	defer st.Close()

	h := rpc.NewHandler(st, logger)

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"status":"ok"}`)) })
	r.PathPrefix("/").Handler(auth.JWT(cfg.JWTSecret)(h))

	srv := &http.Server{
		Addr: ":8899",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	logger.Info("listening", zap.String("addr", ln.Addr().String()))

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return srv.Serve(tls.NewListener(ln, srv.TLSConfig)) })
	g.Go(func() error {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutCtx)
	})
	return g.Wait()
}
