#!/usr/bin/env bash
set -euo pipefail
VERSION=${1:?usage: release.sh v0.1.0}
git tag -a "$VERSION" -m "release $VERSION"
git push origin "$VERSION"
echo "Tagged $VERSION – GitHub Actions will build & push the image."
