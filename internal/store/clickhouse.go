// SPDX-License-Identifier: AGPL-3.0-or-later
package store

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gagliardetto/solana-go"
)

type ClickHouse struct {
	conn driver.Conn
}

func NewClickHouse(addr, database, user, pass string, codec map[string]string) (*ClickHouse, error) {
	opt := &clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: pass,
		},
	}
	conn, err := clickhouse.Open(opt)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	ch := &ClickHouse{conn: conn}
	if err := ch.createTables(codec); err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *ClickHouse) createTables(codec map[string]string) error {
	blockCodec := codec["blocks"]
	if blockCodec == "" {
		blockCodec = "ZSTD(3)"
	}
	q := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS blocks (
	slot UInt64,
	blockhash FixedString(32),
	parent_slot UInt64,
	block_time Int64
) ENGINE = MergeTree
ORDER BY slot
SETTINGS index_granularity = 8192
%s`, blockCodec)
	if err := c.conn.Exec(context.Background(), q); err != nil {
		return err
	}

	txCodec := codec["transactions"]
	if txCodec == "" {
		txCodec = "ZSTD(3)"
	}
	q = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS transactions (
	signature FixedString(64),
	slot UInt64,
	idx UInt32,
	tx_bin String,
	meta_bin String
) ENGINE = MergeTree
ORDER BY (slot, idx)
SETTINGS index_granularity = 8192
%s`, txCodec)
	if err := c.conn.Exec(context.Background(), q); err != nil {
		return err
	}

	sigCodec := codec["signatures"]
	if sigCodec == "" {
		sigCodec = "ZSTD(3)"
	}
	q = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS sig_index (
	address FixedString(32),
	slot UInt64,
	idx UInt32,
	signature FixedString(64)
) ENGINE = MergeTree
ORDER BY (address, slot, idx)
SETTINGS index_granularity = 8192
%s`, sigCodec)
	return c.conn.Exec(context.Background(), q)
}

func (c *ClickHouse) GetBlock(ctx context.Context, slot uint64) (*Block, error) {
	row := c.conn.QueryRow(ctx, `SELECT blockhash, parent_slot, block_time FROM blocks WHERE slot = ?`, slot)
	var b Block
	b.Slot = slot
	if err := row.Scan(&b.Blockhash, &b.ParentSlot, &b.BlockTime); err != nil {
		return nil, err
	}
	rows, err := c.conn.Query(ctx, `SELECT signature, idx, tx_bin, meta_bin FROM transactions WHERE slot = ? ORDER BY idx`, slot)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.Signature, &tx.Index, &tx.Tx, &tx.Meta); err != nil {
			return nil, err
		}
		tx.Slot = slot
		b.Txs = append(b.Txs, tx)
	}
	return &b, rows.Err()
}

func (c *ClickHouse) GetTransaction(ctx context.Context, sig solana.Signature) (*Transaction, error) {
	row := c.conn.QueryRow(ctx, `SELECT slot, idx, tx_bin, meta_bin FROM transactions WHERE signature = ?`, sig[:])
	var tx Transaction
	tx.Signature = sig
	if err := row.Scan(&tx.Slot, &tx.Index, &tx.Tx, &tx.Meta); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (c *ClickHouse) GetSignaturesForAddress(ctx context.Context, addr solana.PublicKey, limit int, before, until solana.Signature) ([]solana.Signature, error) {
	q := `SELECT signature FROM sig_index WHERE address = ?`
	args := []interface{}{addr[:]}
	if !before.IsZero() {
		q += ` AND (slot, idx) < (SELECT slot, idx FROM sig_index WHERE signature = ? LIMIT 1)`
		args = append(args, before[:])
	}
	if !until.IsZero() {
		q += ` AND (slot, idx) >= (SELECT slot, idx FROM sig_index WHERE signature = ? LIMIT 1)`
		args = append(args, until[:])
	}
	q += ` ORDER BY slot DESC, idx DESC LIMIT ?`
	args = append(args, limit)
	rows, err := c.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	var sigs []solana.Signature
	for rows.Next() {
		var s solana.Signature
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		sigs = append(sigs, s)
	}
	return sigs, rows.Err()
}

func (c *ClickHouse) GetBlocksWithLimit(ctx context.Context, start uint64, limit uint64) ([]uint64, error) {
	rows, err := c.conn.Query(ctx, `SELECT slot FROM blocks WHERE slot >= ? ORDER BY slot LIMIT ?`, start, limit)
	if err != nil {
		return nil, err
	}
	var slots []uint64
	for rows.Next() {
		var s uint64
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}

func (c *ClickHouse) GetBlockTime(ctx context.Context, slot uint64) (*time.Time, error) {
	row := c.conn.QueryRow(ctx, `SELECT block_time FROM blocks WHERE slot = ?`, slot)
	var t int64
	if err := row.Scan(&t); err != nil {
		return nil, err
	}
	tt := time.Unix(t, 0)
	return &tt, nil
}

func (c *ClickHouse) GetSlot(ctx context.Context) (uint64, error) {
	row := c.conn.QueryRow(ctx, `SELECT max(slot) FROM blocks`)
	var s *uint64
	if err := row.Scan(&s); err != nil {
		return 0, err
	}
	if s == nil {
		return 0, nil
	}
	return *s, nil
}

func (c *ClickHouse) Close() error {
	return c.conn.Close()
}
