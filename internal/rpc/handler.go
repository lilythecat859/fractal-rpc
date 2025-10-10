// internal/rpc/handler.go
package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/lilythecat859/fractal-rpc/internal/store"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler { return &Handler{store: s} }

func (h *Handler) Install(mux *http.ServeMux, base string) {
	mux.HandleFunc(base+"health", h.health)
	mux.HandleFunc(base, h.rpc)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type req struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type resp struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *errObj     `json:"error,omitempty"`
}

type errObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) rpc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var q req
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		writeError(w, q.ID, -32700, "parse error")
		return
	}
	ctx := r.Context()
	var result interface{}
	var err error
	switch q.Method {
	case "getSlot":
		result, err = h.store.GetSlot(ctx)
	case "getBlockTime":
		var slot uint64
		_ = json.Unmarshal(q.Params, &slot)
		result, err = h.store.GetBlockTime(ctx, slot)
	case "getBlock":
		var slot uint64
		_ = json.Unmarshal(q.Params, &slot)
		result, err = h.store.GetBlock(ctx, slot)
	default:
		writeError(w, q.ID, -32601, "method not found")
		return
	}
	if err != nil {
		writeError(w, q.ID, -32603, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp{JSONRPC: "2.0", ID: q.ID, Result: result})
}

func writeError(w http.ResponseWriter, id interface{}, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &errObj{Code: code, Message: msg},
	})
}
