// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"context"
	"encoding/json"
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

// ---------- JSON-RPC stubs (append only) ----------
type rpcReq struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}
type rpcResp struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	var req rpcReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	var result interface{}
	var rpcErr *rpcError
	switch req.Method {
	case "getSlot":
		result = uint64(42)
	default:
		rpcErr = &rpcError{Code: -32601, Message: "Method not found"}
	}
	resp := rpcResp{JSONRPC: "2.0", ID: req.ID}
	if rpcErr != nil {
		resp.Error = rpcErr
	} else {
		resp.Result = result
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
// ---------- end append ----------

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := config.MustLoad("example.toml")

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// single JSON-RPC entry point (append)
	r.HandleFunc("/", rpcHandler).Methods("POST")

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
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
