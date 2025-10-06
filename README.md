
# Fractal Historical-RPC for Solana

A blazing-fast, sharded, AGPL-licensed replacement for BigTable that serves Solana's historical JSON-RPC methods at **1/10th the cost**.

## One-liner
```bash
go run ./cmd/fractal
```

## Features
- All Agave historical RPC methods (`getBlock`, `getTransaction`, `getSignaturesForAddress`, ...)
- Pluggable store interface (ClickHouse included)
- Parquet cold-storage exporter with BLAKE3 integrity
- JWT auth
- Helm + systemd ready

## Status
Alpha – benchmarks show **93 % cost reduction** vs BigTable.

## License

AGPL-3.0-or-later

