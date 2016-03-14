#!/usr/bin/env bash

set -e

export GOPATH=$PWD/gopath
export PATH=${GOPATH//://bin:}/bin:$PATH

pushd $GOPATH/src/github.com/cloudfoundry/java-buildpack-memory-calculator
 go get -v github.com/tools/godep
 scripts/runTests
popd
