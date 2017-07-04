#!/usr/bin/env bash

set -e -u

RELEASE=$1
SNAPSHOT=$2

update_release() {
  local file=$1
  local repository=$2
  local version=$3

  sed -E -i '' "s|(^[ ]*ARTIFACTORY_REPOSITORY:[ ]*).+$|\1$repository|" $file
  sed -E -i '' "s|(^[ ]*VERSION:[ ]*).+$|\1$version|" $file
}

update_release ci/build.yml libs-release-local $RELEASE
git add .
git commit --message "v$RELEASE Release"

git tag -s v$RELEASE -m "v$RELEASE"
git reset --hard HEAD^1

update_release ci/build.yml libs-snapshot-local $SNAPSHOT
git add .
git commit --message "v$SNAPSHOT Development"
