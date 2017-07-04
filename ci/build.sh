#!/usr/bin/env bash

set -e -u

function build() {
  local type=$1
  local target="java-buildpack-memory-calculator-${type}"

  GOOS=${type} go build -a \
    && mv java-buildpack-memory-calculator ${target} \
    && tar -czf ${target}.tar.gz ${target} \
    && echo "Built ${target}.tar.gz"
}

function upload() {
  local source=$1
  local destination=$2

  JFROG_CLI_OFFER_CONFIG=false /usr/local/bin/jfrog rt upload \
    --url https://repo.spring.io \
    --user $ARTIFACTORY_USERNAME \
    --password $ARTIFACTORY_PASSWORD \
    $1 $2
}

if [[ $GOPATH == "/go" ]]; then
  GOPATH=$PWD/gopath
fi

pushd $GOPATH/src/github.com/cloudfoundry/java-buildpack-memory-calculator
  build darwin
  build linux

  upload \
    java-buildpack-memory-calculator-darwin.tar.gz \
    $ARTIFACTORY_REPOSITORY/org/cloudfoundry/java-buildpack-memory-calculator/$VERSION/java-buildpack-memory-calculator-$(echo $VERSION | sed "s|SNAPSHOT|$(date '+%Y%m%d.%H%M%S')|")-darwin.tar.gz
  upload \
    java-buildpack-memory-calculator-linux.tar.gz \
    $ARTIFACTORY_REPOSITORY/org/cloudfoundry/java-buildpack-memory-calculator/$VERSION/java-buildpack-memory-calculator-$(echo $VERSION | sed "s|SNAPSHOT|$(date '+%Y%m%d.%H%M%S')|")-linux.tar.gz
popd
