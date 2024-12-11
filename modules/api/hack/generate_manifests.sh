#!/bin/bash
set -ex

CRD_OPTIONS="$1"
GENS="$2"

export GOFLAGS=-mod=readonly

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${KUBE_ROOT}" || exit

if grep -qw "deepcopy" <<<"${GENS}"; then
  go run ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go object:headerFile=./hack/boilerplate.go.txt paths=./pkg/api/...
else
  go run ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go object:headerFile=./hack/boilerplate.go.txt paths=./pkg/api/... rbac:roleName=controller-perms "${CRD_OPTIONS}" output:crd:artifacts:config=config/crds
fi
