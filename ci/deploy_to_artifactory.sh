#!/usr/bin/env bash

set -euo pipefail

readonly filename="${1}"
readonly architecture="${2}"
readonly version="${3}"

mvn deploy:deploy-file \
	-Durl=https://artifact.instana.io/artifactory/agent-releases  \
	-Dfile=${filename} \
	-DrepositoryId=agent-releases \
	-DgroupId=com.instana \
	-DartifactId=memory_calculator \
    -Dclassifier=${architecture}\
	-Dpackaging=bin \
	-Dversion=${version}

mvn deploy:deploy-file \
	-Durl=https://delivery.instana.io/rel-generic-agent-local  \
	-Dfile=${filename} \
	-DrepositoryId=agent-releases \
	-DgroupId=com.instana \
	-DartifactId=memory_calculator \
    -Dclassifier=${architecture}\
	-Dpackaging=bin \
	-Dversion=${version} \
	-Dusername=${DELIVERY_ARTIFACTORY_USERNAME} \
	-Dpassword=${DELIVERY_ARTIFACTORY_PASSWORD}
