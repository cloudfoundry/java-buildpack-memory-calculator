#!/usr/bin/env bash

set -euo pipefail

cd "${SOURCE_DIR}"

if [[ -d "${PWD}/go-module-cache" && ! -d "${GOPATH}/pkg/mod" ]]; then
  mkdir -p "${GOPATH}/pkg"
  ln -s "${PWD}/go-module-cache" "${GOPATH}/pkg/mod"
fi

declare -a ARCHITECTURES_LINUX=("amd64" "arm64" "s390x")

for arch in "${ARCHITECTURES_LINUX[@]}"
do
  GOOS=linux GOARCH=${arch} go build -ldflags='-s -w' -o "${TARGET_DIR}/memory-calculator_linux_${arch}" main.go
done

GOOS=windows GOARCH=amd64 go build -ldflags='-s -w' -o "${TARGET_DIR}/memory-calculator_win_amd64" main.go
#GOOS=zos GOARCH=s390 go build -ldflags='-s -w' -o "${TARGET_DIR}/memory-calculator_zos_s390x" main.go
#GOOS=solaris GOARCH=sparc64 go build -ldflags='-s -w' -o "${TARGET_DIR}/memory-calculator_solaris_sparc64" main.go
