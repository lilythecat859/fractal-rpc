// SPDX-License-Identifier: AGPL-3.0-or-later
package store

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
)

// Store is the pluggable backend interface.
type Store interface {
	GetBlock(ctx context.Context, slot uint64) (*Block, error)
	GetTransaction(ctx context.Context, sig solana.Signature) (*Transaction, error)
	GetSignaturesForAddress(ctx context.Context, addr solana.PublicKey, limit int, before solana.Signature, until solana.Signature) ([]solana.Signature, error)
	GetBlocksWithLimit(ctx context.Context, start uint64, limit uint64) ([]uint64, error)
	GetBlockTime(ctx context.Context, slot uint64) (*time.Time, error)
	GetSlot(ctx context.Context) (uint64, error)
	Close() error
}

type Block struct {
	Slot       uint64
	Blockhash  solana.Hash
	ParentSlot uint64
	BlockTime  int64
	Txs        []Transaction
}

type Transaction struct {
	Signature solana.Signature
	Slot      uint64
	Index     uint32
	Tx        []byte // bincode
	Meta      []byte // bincode
}
