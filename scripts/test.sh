#!/usr/bin/env bash
set -euo pipefail
URL=${URL:-http://localhost:8899}
SLOT=${1:-30000000}  # any old slot with data

echo "==> getSlot"
curl -sX POST "$URL" -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"getSlot"}' | jq .

echo "==> getBlockTime ($SLOT)"
curl -sX POST "$URL" -H 'content-type: application/json' \
  --data "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"getBlockTime\",\"params\":[$SLOT]}" | jq .

echo "==> getBlock ($SLOT)"
curl -sX POST "$URL" -H 'content-type: application/json' \
  --data "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"getBlock\",\"params\":[$SLOT]}" | jq .
  
