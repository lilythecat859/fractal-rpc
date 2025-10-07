#!/usr/bin/env bash
set -euo pipefail
sudo useradd -r -s /bin/false fractal || true
sudo cp fractal.toml /etc/fractal.toml
sudo cp fractal /usr/local/bin/fractal
sudo chmod +x /usr/local/bin/fractal
sudo cp fractal.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now fractal
