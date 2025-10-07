#!/usr/bin/env bash
set -euo pipefail
echo "=== Fractal repo sanity check ==="
test -f go.mod || { echo "go.mod missing"; exit 1; }
test -f Dockerfile || { echo "Dockerfile missing"; exit 1; }
test -f fractal.toml.example || { echo "fractal.toml.example missing"; exit 1; }
test -f helm/fractal/Chart.yaml || { echo "helm chart missing"; exit 1; }
test -f .github/workflows/ci.yml || { echo "CI workflow missing"; exit 1; }
test -f scripts/test.sh || { echo "test.sh missing"; exit 1; }
test -f LICENSE || { echo "LICENSE missing"; exit 1; }
echo "All required files present ✅"
