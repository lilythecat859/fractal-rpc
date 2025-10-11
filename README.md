 
# Fractal – Historical-RPC for Solana

Open-source, AGPL-3.0.  
Stores full blocks & transactions in ClickHouse and serves them over JSON-RPC.

## Features

- All Agave historical RPC methods (`getBlock`, `getTransaction`, `getSignaturesForAddress`, …)  
- Pluggable store interface (ClickHouse included)  
- Parquet cold-storage exporter with BLAKE3 integrity  
- JWT auth  
- Helm + systemd ready  

## Quick start (Docker Compose)

1. Clone repo  
2. Copy example config: `cp fractal.toml.example fractal.toml`  
3. `docker compose up -d`  
4. RPC endpoint: `http://localhost:8899/`

## Methods

- `getBlock` – returns full block with tx list  
- `getTransaction` – returns tx + meta  
- `getSignaturesForAddress` – paginated history  
- `getBlocksWithLimit` – slot list  
- `getBlockTime` – unix ts  
- `getSlot` – latest indexed slot

## Build from source

 

go 1.22+ required make build   # produces ./fractal

 

## Kubernetes

 

helm install fractal ./helm/fractal

 

## Environment variables

`CGO_ENABLED=1` (mandatory for ClickHouse driver)

## License

AGPL-3.0-or-later

---

> “The fractal nature of reality is that you can go forever in, and you get the same amount of detail forever out.”  
> — Terence McKenna, *Evolving Times* lecture, Boulder, Colorado, October 1994
 

## Screenshots  
Folder: https://drive.google.com/drive/folders/1UuVj5KKO37IEWVDH8ZlOKgx5ztBnldCh
