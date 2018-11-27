#!/usr/bin/env bash

set -e -u

if [[ $GOPATH == "/go" ]]; then
  GOPATH=$PWD/gopath
fi

PATH=${GOPATH//://bin:}/bin:$PATH
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
PARENT_SCRIPT_DIR="$( cd "$( dirname "${SCRIPT_DIR}" )" >/dev/null && pwd )"

pushd ${PARENT_SCRIPT_DIR}
  go install -v github.com/onsi/ginkgo/ginkgo
  ginkgo -r -failOnPending -randomizeAllSpecs -skipMeasurements=true -race "$@"
popd
