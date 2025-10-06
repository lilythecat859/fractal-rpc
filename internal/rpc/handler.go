// SPDX-License-Identifier: AGPL-3.0-or-later
package rpc

import (
	"net/http"

	"github.com/gagliardetto/solana-go"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/lilythecat859/solana-historical-rpc/internal/store"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Install(mux *http.ServeMux, path string) {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(h, "")
	mux.Handle(path, s)
}

type GetBlockArgs struct {
	Slot       uint64 `json:"slot,string"`
	Commitment string `json:"commitment,omitempty"`
}

type GetBlockReply struct {
	Block *store.Block `json:"block"`
}

func (h *Handler) GetBlock(r *http.Request, args *GetBlockArgs, reply *GetBlockReply) error {
	ctx := r.Context()
	b, err := h.store.GetBlock(ctx, args.Slot)
	if err != nil {
		return err
	}
	reply.Block = b
	return nil
}

type GetTransactionArgs struct {
	Signature solana.Signature `json:"signature"`
}

type GetTransactionReply struct {
	Tx *store.Transaction `json:"transaction"`
}

func (h *Handler) GetTransaction(r *http.Request, args *GetTransactionArgs, reply *GetTransactionReply) error {
	ctx := r.Context()
	tx, err := h.store.GetTransaction(ctx, args.Signature)
	if err != nil {
		return err
	}
	reply.Tx = tx
	return nil
}

type GetSignaturesForAddressArgs struct {
	Address solana.PublicKey `json:"address"`
	Limit   int              `json:"limit"`
	Before  solana.Signature `json:"before"`
	Until   solana.Signature `json:"until"`
}

type GetSignaturesForAddressReply struct {
	Signatures []solana.Signature `json:"signatures"`
}

func (h *Handler) GetSignaturesForAddress(r *http.Request, args *GetSignaturesForAddressArgs, reply *GetSignaturesForAddressReply) error {
	ctx := r.Context()
	if args.Limit == 0 || args.Limit > 1000 {
		args.Limit = 1000
	}
	sigs, err := h.store.GetSignaturesForAddress(ctx, args.Address, args.Limit, args.Before, args.Until)
	if err != nil {
		return err
	}
	reply.Signatures = sigs
	return nil
}

type GetBlocksWithLimitArgs struct {
	Start uint64 `json:"start_slot,string"`
	Limit uint64 `json:"limit"`
}

type GetBlocksWithLimitReply struct {
	Slots []uint64 `json:"slots"`
}

func (h *Handler) GetBlocksWithLimit(r *http.Request, args *GetBlocksWithLimitArgs, reply *GetBlocksWithLimitReply) error {
	ctx := r.Context()
	if args.Limit == 0 || args.Limit > 5000 {
		args.Limit = 5000
	}
	slots, err := h.store.GetBlocksWithLimit(ctx, args.Start, args.Limit)
	if err != nil {
		return err
	}
	reply.Slots = slots
	return nil
}

type GetBlockTimeArgs struct {
	Slot uint64 `json:"slot,string"`
}

type GetBlockTimeReply struct {
	Time *int64 `json:"blockTime"`
}

func (h *Handler) GetBlockTime(r *http.Request, args *GetBlockTimeArgs, reply *GetBlockTimeReply) error {
	ctx := r.Context()
	t, err := h.store.GetBlockTime(ctx, args.Slot)
	if err != nil {
		return err
	}
	if t != nil {
		reply.Time = new(int64)
		*reply.Time = t.Unix()
	}
	return nil
}

type GetSlotReply struct {
	Slot uint64 `json:"slot"`
}

func (h *Handler) GetSlot(r *http.Request, _ *struct{}, reply *GetSlotReply) error {
	ctx := r.Context()
	s, err := h.store.GetSlot(ctx)
	if err != nil {
		return err
	}
	reply.Slot = s
	return nil
}
