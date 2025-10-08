package rpc

import (
	"encoding/json"
	"net/http"
	"time"
)

type Store interface{}

type Handler struct{}

func NewHandler(s Store) *Handler { return &Handler{} }

func (h *Handler) Install(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, h.serveJSONRPC)
}

func (h *Handler) serveJSONRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "error": map[string]interface{}{"code": -32700, "message": "Parse error"}, "id": nil})
		return
	}

	switch req.Method {
	case "getSlot":
		writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "result": uint64(123456), "id": req.ID})

	case "getBlockTime":
		var p []uint64
		_ = json.Unmarshal(req.Params, &p)
		if len(p) == 0 {
			writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "error": map[string]interface{}{"code": -32602, "message": "missing block number"}, "id": req.ID})
			return
		}
		writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "result": time.Now().Unix(), "id": req.ID})

	case "getBlocksWithLimit":
		var p []uint64
		_ = json.Unmarshal(req.Params, &p)
		start, limit := uint64(42), uint64(3)
		if len(p) >= 1 {
			start = p[0]
		}
		if len(p) >= 2 {
			limit = p[1]
		}
		end := start + limit
		blocks := make([]uint64, 0, limit)
		for i := start; i < end; i++ {
			blocks = append(blocks, i)
		}
		writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "result": blocks, "id": req.ID})

	case "getSignaturesForAddress":
		writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "result": []string{"5VERv8NMvzbJMYkT4MbsVJdJzRkGSSkMyzyevgjr3zFv7nKfJXs11111111111111111111111111111111"}, "id": req.ID})

	default:
		writeJSON(w, map[string]interface{}{"jsonrpc": "2.0", "error": map[string]interface{}{"code": -32601, "message": "Method not found"}, "id": req.ID})
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
