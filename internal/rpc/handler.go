package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/lilythecat859/fractal-rpc/internal/store"
)

type Handler struct{ store store.Store }

func NewHandler(s store.Store) *Handler { return &Handler{store: s} }

func (h *Handler) Install(mux *http.ServeMux, base string) {
	mux.HandleFunc(base+"health", h.health)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
