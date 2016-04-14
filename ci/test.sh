#!/usr/bin/env bash

set -e

if [[ $GOPATH == "/go" ]]; then
  GOPATH=$PWD/gopath
fi

PATH=${GOPATH//://bin:}/bin:$PATH

go get -v github.com/tools/godep

pushd $GOPATH/src/github.com/cloudfoundry/java-buildpack-memory-calculator
  GOPATH=$(godep path):$GOPATH
  PATH=${GOPATH//://bin:}/bin:$PATH

  go install -v github.com/onsi/ginkgo/ginkgo
  ginkgo -r -failOnPending -randomizeAllSpecs -skipMeasurements=true -race "$@"
popd
