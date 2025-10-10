// internal/store/clickhouse.go
package store

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
        solana "github.com/lilythecat859/fractal-rpc/solana"

)

type clickHouse struct {
	conn driver.Conn
}

func NewClickHouse(addr, db, user, pass string, codec map[string]string) (Store, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
	})
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &clickHouse{conn: conn}, nil
}

func (c *clickHouse) Close() error { return c.conn.Close() }

func (c *clickHouse) GetSlot(ctx context.Context) (uint64, error) {
	row := c.conn.QueryRow(ctx, `SELECT max(slot) FROM blocks`)
	var slot *uint64
	if err := row.Scan(&slot); err != nil {
		return 0, err
	}
	if slot == nil {
		return 0, fmt.Errorf("no blocks")
	}
	return *slot, nil
}

func (c *clickHouse) GetBlockTime(ctx context.Context, slot uint64) (int64, error) {
	row := c.conn.QueryRow(ctx, `SELECT block_time FROM blocks WHERE slot = ?`, slot)
	var t *int64
	if err := row.Scan(&t); err != nil {
		return 0, err
	}
	if t == nil {
		return 0, fmt.Errorf("slot %d not found", slot)
	}
	return *t, nil
}

func (c *clickHouse) GetBlock(ctx context.Context, slot uint64) (*solana.Block, error) {
	row := c.conn.QueryRow(ctx, `SELECT data FROM blocks WHERE slot = ?`, slot)
	var raw []byte
	if err := row.Scan(&raw); err != nil {
		return nil, err
	}
	var block solana.Block
	if err := block.UnmarshalBinary(raw); err != nil {
		return nil, err
	}
	return &block, nil
}
