CREATE DATABASE IF NOT EXISTS solana_history;

CREATE TABLE IF NOT EXISTS solana_history.blocks
(
    slot       UInt64,
    blockhash  FixedString(44),
    prev_hash  FixedString(44),
    block_time DateTime,
    height     UInt64,
    tx_count   UInt32,
    data       String CODEC(ZSTD(3))
) ENGINE = ReplacingMergeTree()
ORDER BY slot;

CREATE TABLE IF NOT EXISTS solana_history.transactions
(
    slot    UInt64,
    tx_sig  String,
    idx     UInt32,
    signers Array(String),
    fee     UInt64,
    err     UInt8,
    data    String CODEC(ZSTD(3))
) ENGINE = ReplacingMergeTree()
ORDER BY (slot, tx_sig);
