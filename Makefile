
BUILD_TAG ?= latest
BUILD_CACHE_TAG = latest-builder-cache
IMAGE_NAME = otel-collector-distro
IMAGE_NAME_DEV = otel-collector-distro-dev

DOCKERHUB_ORG = logzio
REPO_URL = $(DOCKERHUB_ORG)/$(IMAGE_NAME)
REPO_URL_DEV = $(DOCKERHUB_ORG)/$(IMAGE_NAME_DEV)

ALL_MODULES := $(shell find ./logzio/ -type f -name "go.mod" -exec dirname {} \; | sort )


# Build
.PHONY: build
build:
	@$(MAKE) -C ./otelbuilder/ build


.PHONY: _build
_build:
	DOCKER_BUILDKIT=0 docker build \
		--file $(DOCKERFILE) \
		--build-arg BUILD_TAG=$(TAG) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--tag $(IMG):$(TAG) \
		.

.PHONY: _buildx
_buildx:
	DOCKER_BUILDKIT=0 docker buildx build \
		--platform $(PLATFORM) \
		--file $(DOCKERFILE) \
		--build-arg BUILD_TAG=$(TAG) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--tag $(IMG):$(TAG) \
		--push \
		.

.PHONY: build-container-dev
build-container-dev:
	$(MAKE) _build \
		IMG="$(REPO_URL)-dev" \
		DOCKERFILE="Dockerfile.dev" \
		TAG="$(BUILD_TAG)"

.PHONY: build-container-arm64
build-container-arm64:
	$(MAKE) _buildx \
		PLATFORM="linux/arm64" \
		IMG="$(REPO_URL)" \
		DOCKERFILE="Dockerfile.arm" \
		TAG="$(BUILD_TAG)-arm64"

.PHONY: build-container-amd64
build-container-amd64:
	$(MAKE) _buildx \
		PLATFORM="linux/amd64" \
		IMG="$(REPO_URL)" \
		DOCKERFILE="Dockerfile" \
		TAG="$(BUILD_TAG)-amd64"

.PHONY: build-container
build-container:
	$(MAKE) _build \
		IMG="$(REPO_URL)" \
		DOCKERFILE="Dockerfile" \
		TAG="$(BUILD_TAG)"

# build and push multi arch docker images to logzio docker hub
.PHONY: build-container-multi-platform
build-container-multi-platform: build-container-amd64 build-container-arm64
	$(MAKE) multi-platform-manifest-create-push

.PHONY: multi-platform-manifest-create
multi-platform-manifest-create:
	docker manifest create $(REPO_URL):$(BUILD_TAG) \
		--amend $(REPO_URL):$(BUILD_TAG)-arm64 \
		--amend $(REPO_URL):$(BUILD_TAG)-amd64


.PHONY: multi-platform-manifest-create-latest
multi-platform-manifest-create-latest:
	docker manifest create $(REPO_URL):latest \
		--amend $(REPO_URL):$(BUILD_TAG)-arm64 \
		--amend $(REPO_URL):$(BUILD_TAG)-amd64


.PHONY: multi-platform-manifest-push
 multi-platform-manifest-push:
	docker manifest push $(REPO_URL):$(BUILD_TAG)

.PHONY: multi-platform-manifest-push-latest
 multi-platform-manifest-push-latest:
	docker manifest push $(REPO_URL):latest

.PHONY: multi-platform-manifest-create-push
multi-platform-manifest-create-push: multi-platform-manifest-create multi-platform-manifest-create-latest multi-platform-manifest-push multi-platform-manifest-push-latest

.PHONY: install-tools
install-tools:
	go install github.com/google/addlicense@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/pavius/impi/cmd/impi@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.3

.PHONY: test-componenets
test-components:
	$(MAKE) install-tools
	$(MAKE) for-all CMD="make test"

.PHONY: format-components
format-components:
	$(MAKE) install-tools
	$(MAKE) for-all CMD="make fmt"


.PHONY: lint-components
lint-components:
	$(MAKE) install-tools
	$(MAKE) for-all CMD="make lint"

.PHONY: for-all
for-all:
	@set -e; for dir in $(ALL_MODULES); do \
	  (cd "$${dir}" && \
	  	echo "running $${CMD} in $${dir}" && \
	 	$${CMD} ); \
	done








