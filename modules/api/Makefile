GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
CRD_OPTIONS ?= "crd:allowDangerousTypes=true"

REPO ?= kubesphere
TAG ?= latest

.PHONY: lint test

lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v1.51.2

test:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "total: " $$3 }'

.PHONY: binary
# Build all of binary
binary: | apiserver controller-manager; $(info $(M)...Build all of binary.) @ ## Build all of binary.

# Build gateway-apiserver binary
apiserver: ; $(info $(M)...Begin to build apiserver binary.)  @ ## Build apiserver.
	 hack/gobuild.sh cmd/apiserver;

# Build gateway-controller-manager binary
controller-manager: ; $(info $(M)...Begin to build controller-manager binary.)  @ ## Build gateway-controller-manager.
	hack/gobuild.sh cmd/controller-manager

# Run go fmt against code
fmt: ;$(info $(M)...Begin to run go fmt against code.)  @ ## Run go fmt against code.
	gofmt -w ./pkg ./cmd

# Generate manifests e.g. CRD, RBAC etc.
manifests: ;$(info $(M)...Begin to generate manifests e.g. CRD, RBAC etc..)  @ ## Generate manifests e.g. CRD, RBAC etc.
	hack/generate_manifests.sh ${CRD_OPTIONS}

# Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
deepcopy: ;$(info $(M)...Begin to deepcopy.)  @ ## Deepcopy.
	hack/generate_manifests.sh ${CRD_OPTIONS} "deepcopy"

# Format all import, `goimports` is required.
goimports: ;$(info $(M)...Begin to Format all import.)  @ ## Format all import, `goimports` is required.
	@hack/update-goimports.sh

# Run go vet against code
vet: ;$(info $(M)...Begin to run go vet against code.)  @ ## Run go vet against code.
	go vet ./pkg/... ./cmd/...

container: ;$(info $(M)...Begin to build the docker image.)  @ ## Build the docker image.
	TAG=${TAG} REPO=${REPO} DRY_RUN=true hack/docker_build.sh

container-push: ;$(info $(M)...Begin to build and push.)  @ ## Build and Push.
	TAG=${TAG} REPO=${REPO} hack/docker_build.sh

container-cross: ; $(info $(M)...Begin to build container images for multiple platforms.)  @ ## Build container images for multiple platforms. Currently, only linux/amd64,linux/arm64 are supported.
	TAG=${TAG} REPO=${REPO} DRY_RUN=true hack/docker_build_multiarch.sh

container-cross-push: ; $(info $(M)...Begin to build and push.)  @ ## Build and Push.
	TAG=${TAG} REPO=${REPO} PLATFORMS=${PLATFORMS} hack/docker_build_multiarch.sh
