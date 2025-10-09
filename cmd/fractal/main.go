package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/lilythecat859/fractal-rpc/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := config.MustLoad("example.toml")

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
	logger.Info("listening", zap.String("addr", ln.Addr().String()))

	go srv.Serve(ln)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
