#!/usr/bin/env bash

set -euo pipefail

readonly filepath="${1}"
readonly architecture="${2}"
readonly version="${3}"

set +x
readonly base_url="${DELIVERY_ARTIFACTORY_GENERIC_BASE_URL}"
readonly org_path=com/instana
readonly artifact_id="memory_calculator"
readonly target_url="${base_url}/${org_path}/${artifact_id}/${version}/${artifact_id}-${version}-${architecture}"
echo -e "[\e[1m\e[34mINFO\e[39m\e[21m] Uploading \e[92m$target_url\e[39m"
readonly status_code=$(curl --silent --output /dev/stderr --write-out "%{http_code}" -u "${DELIVERY_ARTIFACTORY_USERNAME}:${DELIVERY_ARTIFACTORY_PASSWORD}" -X PUT "$target_url" -T "${filepath}")
echo # curl doesn't output a line break on error
if test "$status_code" -ne 201; then
  echo -e "[\e[1m\e[91mERROR\e[39m\e[21m] curl returned $status_code, exiting..."
  exit 1
fi
set -x
