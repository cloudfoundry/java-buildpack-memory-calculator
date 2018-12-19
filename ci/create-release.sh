#!/usr/bin/env bash

set -euo pipefail

RELEASE=$1
SNAPSHOT=$2

update_version() {
  sed -E -i '' "s|(^[ ]*VERSION:[ ]*).+$|\1${1}|" ci/package.yml
}

git branch -f release
git checkout release
update_version ${RELEASE}
git add .
git commit --message "v${RELEASE} Release"
git tag -s v${RELEASE} -m "v${RELEASE}"

git master
update_version ${SNAPSHOT}
git add .
git commit --message "v${SNAPSHOT} Development"
