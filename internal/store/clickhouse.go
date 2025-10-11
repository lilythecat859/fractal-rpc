package store

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type clickHouse struct {
	conn driver.Conn
}

func NewClickHouse(addr, db, user, pass string) (Store, error) {
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
	row := c.conn.QueryRow(ctx, "SELECT max(slot) FROM blocks")
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
	row := c.conn.QueryRow(ctx, "SELECT block_time FROM blocks WHERE slot = ?", slot)
	var t *int64
	if err := row.Scan(&t); err != nil {
		return 0, err
	}
	if t == nil {
		return 0, fmt.Errorf("slot %d not found", slot)
	}
	return *t, nil
}

func (c *clickHouse) GetBlock(ctx context.Context, slot uint64) ([]byte, error) {
	row := c.conn.QueryRow(ctx, "SELECT data FROM blocks WHERE slot = ?", slot)
	var raw []byte
	if err := row.Scan(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *clickHouse) GetSignaturesForAddress(ctx context.Context, addr string, before, until string, limit uint) ([]string, error) {
	rows, err := c.conn.Query(ctx,
		`SELECT tx_sig FROM transactions
		  WHERE has(signers, ?)
		    AND (?, '') = ('', '')
		    AND (?, '') = ('', '')
		  ORDER BY slot DESC
		  LIMIT ?`, addr, before, until, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var sig string
		if err := rows.Scan(&sig); err != nil {
			return nil, err
		}
		out = append(out, sig)
	}
	return out, nil
}
