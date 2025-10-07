package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lilythecat859/fractal/internal/debug"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "fractal.toml", "path to config file")
	flag.Parse()

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg config) error {
	router := http.NewServeMux()

	// ---- DEBUG / METRICS / PPROF ----
	debug.Mount(router)

	// ---- RPC ENDPOINTS ----
	// (your real handlers will be attached here)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"jsonrpc":"2.0","result":"pong","id":1}` + "\n"))
	})

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutCtx)
}
