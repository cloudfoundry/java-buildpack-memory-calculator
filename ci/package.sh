#!/usr/bin/env bash

set -euo pipefail

if [[ -d $PWD/go-module-cache && ! -d ${GOPATH}/pkg/mod ]]; then
  mkdir -p ${GOPATH}/pkg
  ln -s $PWD/go-module-cache ${GOPATH}/pkg/mod
fi

TARGET="${PWD}/artifactory/org/cloudfoundry/java-buildpack-memory-calculator/${VERSION}/java-buildpack-memory-calculator-$(echo ${VERSION} | sed "s|SNAPSHOT|$(date '+%Y%m%d.%H%M%S')-1|").tgz"

cd "$(dirname "${BASH_SOURCE[0]}")/.."
go build -ldflags='-s -w' -o bin/java-buildpack-memory-calculator main.go

cd bin
mkdir -p $(dirname ${TARGET})
tar czf ${TARGET} java-buildpack-memory-calculator
