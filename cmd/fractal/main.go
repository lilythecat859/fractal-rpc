package main

import (
	"log"
	"os"

	"github.com/lilythecat859/fractal-rpc/internal/config"
	"github.com/lilythecat859/fractal-rpc/internal/server"
)

func main() {
	cfg := config.MustLoad()
	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
