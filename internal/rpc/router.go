
package rpc

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	store Store
	log   *zap.Logger
}

type Store interface {
	GetSlot(ctx context.Context) (uint64, error)
	GetBlock(ctx context.Context, slot uint64) ([]byte, error)
	GetBlockTime(ctx context.Context, slot uint64) (int64, error)
	GetSignaturesForAddress(ctx context.Context, addr string, before string, until string, limit uint) ([]string, error)
	Close() error
}

func NewHandler(s Store, log *zap.Logger) http.Handler {
	h := &Handler{store: s, log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.serve)
	return mux
}

func (h *Handler) serve(w http.ResponseWriter, r *http.Request) {
	var req jsonRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}
	ctx := r.Context()
	switch req.Method {
	case "getSlot":
		slot, _ := h.store.GetSlot(ctx)
		writeResult(w, req.ID, slot)
	case "getBlockTime":
		slot := uint64(req.Params[0].(float64))
		ts, _ := h.store.GetBlockTime(ctx, slot)
		writeResult(w, req.ID, ts)
	case "getBlock":
		slot := uint64(req.Params[0].(float64))
		raw, _ := h.store.GetBlock(ctx, slot)
		writeResult(w, req.ID, json.RawMessage(raw))
	case "getSignaturesForAddress":
		addr := req.Params[0].(string)
		limit := uint(1000)
		sigs, _ := h.store.GetSignaturesForAddress(ctx, addr, "", "", limit)
		writeResult(w, req.ID, sigs)
	default:
		writeError(w, req.ID, -32601, "Method not found")
	}
}

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  []interface{}   `json:"params"`
}

func writeResult(w http.ResponseWriter, id interface{}, result interface{}) {
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc": "2.0", "id": id, "result": result})
}
func writeError(w http.ResponseWriter, id interface{}, code int, msg string) {
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc": "2.0", "id": id, "error": map[string]interface{}{"code": code, "message": msg}})
}
