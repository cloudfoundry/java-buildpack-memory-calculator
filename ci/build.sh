#!/usr/bin/env bash

set -e

function build() {
  local type=$1
  local target="java-buildpack-memory-calculator-${type}"

  GOOS=${type} go build -a \
    && mv java-buildpack-memory-calculator ${target} \
    && tar -czf ${target}.tar.gz ${target} \
    && echo "Built ${target}.tar.gz"
}

if [[ $GOPATH == "/go" ]]; then
  GOPATH=$PWD/gopath
fi

pushd $GOPATH/src/github.com/cloudfoundry/java-buildpack-memory-calculator
  build linux
  build darwin
popd
