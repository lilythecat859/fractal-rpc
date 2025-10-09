package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lilythecat859/fractal-rpc/internal/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *config.Config, logger *zap.Logger) error {
	mux := mux.NewRouter()

	// --- health check route ---
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// --- mount your RPC handler on everything else ---
	rpc := NewHandler(nil) // replace with real handler when ready
	mux.PathPrefix("/").Handler(rpc)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: mux,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	logger.Info("listening", zap.String("addr", ln.Addr().String()))

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return srv.Serve(ln) })
	g.Go(func() error {
		<-ctx.Done()
		shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shCtx)
	})
	return g.Wait()
}
