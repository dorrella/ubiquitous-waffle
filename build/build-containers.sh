#!/usr/bin/env bash

BUILD_EXEC=${BUILD_EXEC:-'buildah'}
OCI_REGISTRY=${OCI_REGISTRY:-'localhost:31000'}
BUILD_TAG=${BUILD_TAG:-'latest'}

SERVICE_TAG="${OCI_REGISTRY}/service:${BUILD_TAG}"
INIT_TAG="${OCI_REGISTRY}/init:${BUILD_TAG}"
COLLECTOR_TAG="${OCI_REGISTRY}/otel-collector:${BUILD_TAG}"

#make sure we are in the right dir
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd ${SCRIPT_DIR}/..

#start with otel collector
${BUILD_EXEC} build -f build/collector.Dockerfile -t ${COLLECTOR_TAG} .
${BUILD_EXEC} push --tls-verify=false ${COLLECTOR_TAG}

#todo common apk base

#service
${BUILD_EXEC} build -f build/service.Dockerfile -t ${SERVICE_TAG} .
${BUILD_EXEC} push --tls-verify=false ${SERVICE_TAG}

#init
${BUILD_EXEC} build -f build/init.Dockerfile -t ${INIT_TAG} .
${BUILD_EXEC} push --tls-verify=false ${INIT_TAG}
