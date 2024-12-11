#!/usr/bin/env bash

set -ex
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

# push to kubesphere with default latest tag
TAG=${TAG:-}
REPO=${REPO:-}
PUSH=${PUSH:-}
# supported platforms
PLATFORMS=${PLATFORMS:-linux/amd64,linux/arm64}

# support other container tools. e.g. podman
CONTAINER_CLI=${CONTAINER_CLI:-docker}
CONTAINER_BUILDER=${CONTAINER_BUILDER:-"buildx build"}

# If set, just building, no pushing
if [[ -z "${DRY_RUN:-}" ]]; then
  PUSH="--push"
fi

# shellcheck disable=SC2086 # inteneded splitting of CONTAINER_BUILDER
${CONTAINER_CLI} ${CONTAINER_BUILDER} \
  --platform ${PLATFORMS} \
  ${PUSH} \
  -f build/apiserver/Dockerfile \
  -t "${REPO}"/gateway-apiserver:"${TAG}" .

# shellcheck disable=SC2086 # intended splitting of CONTAINER_BUILDER
${CONTAINER_CLI} ${CONTAINER_BUILDER} \
  --platform ${PLATFORMS} \
  ${PUSH} \
  -f build/controller-manager/Dockerfile \
  -t "${REPO}"/gateway-controller-manager:"${TAG}" .

