#!/usr/bin/env bash
# Requires: bq (gcloud sdk), clickhouse-client
set -euo pipefail
DATASET="bigquery-public-data.crypto_solana"
CH_HOST="${CH_HOST:-localhost}"
CH_DB="solana_history"

echo "=> Importing blocks 30M..30.1M (demo slice)"
bq query --max_rows 100000 --format csv \
  "SELECT slot,block_hash,prev_hash,block_time,block_height,transaction_count
   FROM \`${DATASET}.blocks\`
   WHERE slot BETWEEN 30000000 AND 30100000
   ORDER BY slot" \
| clickhouse-client --host="$CH_HOST" --query="
  INSERT INTO ${CH_DB}.blocks FORMAT CSVWithNames"
  
