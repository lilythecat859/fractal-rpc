// cmd/fractal/main.go
package main

import (
	"context"
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

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := config.MustLoad()

	db, err := store.NewClickHouse(
		cfg.ClickHouse.Addr,
		cfg.ClickHouse.Database,
		cfg.ClickHouse.Auth.Username,
		cfg.ClickHouse.Auth.Password,
		cfg.ClickHouse.Codec,
	)
	if err != nil {
		logger.Fatal("store", zap.Error(err))
	}
	defer db.Close()

	h := rpc.NewHandler(db)
	mux := http.NewServeMux()
	h.Install(mux, cfg.RPCPath)

	srv := &http.Server{Addr: ":" + fmt.Sprint(cfg.HTTPPort), Handler: mux}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
	logger.Info("listening", zap.String("addr", ln.Addr().String()))

	go srv.Serve(ln)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
