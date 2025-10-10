// internal/store/store.go
package store

import (
	"context"
	"github.com/gagliardetto/solana-go"
)

type Store interface {
	GetSlot(ctx context.Context) (uint64, error)
	GetBlock(ctx context.Context, slot uint64) (*solana.Block, error)
	GetBlockTime(ctx context.Context, slot uint64) (int64, error)
	Close() error
}
