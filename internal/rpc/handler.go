package rpc

import "net/http"

type Handler struct{}

func NewHandler(s Store) *Handler { return &Handler{} }

func (h *Handler) Install(mux *http.ServeMux, path string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","result":"pong","id":1}` + "\n"))
	})
}
