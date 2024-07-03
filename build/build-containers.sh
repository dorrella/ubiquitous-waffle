#!/usr/bin/env bash

BUILD_EXEC=${BUILD_EXEC:-'buildah'}
OCI_REGISTRY=${OCI_REGISTRY:-'localhost:31000'}
BUILD_TAG=${BUILD_TAG:-'latest'}

SERVICE_TAG="${OCI_REGISTRY}/service:latest"

#make sure we are in the right dir
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd ${SCRIPT_DIR}/..

${BUILD_EXEC} build -f build/service.Dockerfile -t ${SERVICE_TAG} .
${BUILD_EXEC} push --tls-verify=false ${SERVICE_TAG}
