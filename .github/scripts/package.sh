#!/usr/bin/env bash
# Builds the calculator for linux/amd64 and linux/arm64 and packs each binary into a tgz.
# Usage: package.sh <version> <output dir>
#
# Each archive is named after its Go platform: memory-calculator-<version>-linux-amd64.tgz and
# memory-calculator-<version>-linux-arm64.tgz. Each archive holds one file, java-buildpack-memory-calculator,
# as in earlier releases. A SHA256SUMS file lists the archives.
set -euo pipefail

VERSION=$1
OUT=$2

cd "$(dirname "${BASH_SOURCE[0]}")/../.."
mkdir -p "${OUT}"
OUT=$(cd "${OUT}" && pwd)

# Fixed timestamps and ordering keep the archives identical for the same commit. -buildvcs=false keeps the
# checkout state (such as an untracked dist/ directory) out of the binary.
SOURCE_DATE_EPOCH=$(git log -1 --format=%ct)
WORK=$(mktemp -d)
trap 'rm -rf "${WORK}"' EXIT

for arch in amd64 arm64; do
  name="memory-calculator-${VERSION}-linux-${arch}.tgz"

  mkdir -p "${WORK}/${arch}"
  CGO_ENABLED=0 GOOS=linux GOARCH="${arch}" \
    go build -trimpath -buildvcs=false -ldflags="-s -w -X main.version=${VERSION}" -o "${WORK}/${arch}/java-buildpack-memory-calculator" .

  tar --sort=name --mtime="@${SOURCE_DATE_EPOCH}" --owner=0 --group=0 --numeric-owner \
    -cf - -C "${WORK}/${arch}" java-buildpack-memory-calculator | gzip -n > "${OUT}/${name}"
done

(cd "${OUT}" && sha256sum -- *.tgz > SHA256SUMS)
cat "${OUT}/SHA256SUMS"
