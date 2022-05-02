# Build
.PHONY: build
build:
	@$(MAKE) -C ./otelcolbuilder/ build

BUILD_TAG ?= latest
BUILD_CACHE_TAG = latest-builder-cache
IMAGE_NAME = logzio-otel-collector
IMAGE_NAME_DEV = logzio-otel-collector-dev

DOCKERHUB_ORG = yotamloe
REPO_URL = $(DOCKERHUB_ORG)/$(IMAGE_NAME)
REPO_URL_DEV = $(DOCKERHUB_ORG)/$(IMAGE_NAME_DEV)

.PHONY: _build
_build:
	DOCKER_BUILDKIT=0 docker build \
		--file $(DOCKERFILE) \
		--build-arg BUILD_TAG=$(TAG) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--tag $(IMG):$(TAG) \
		.

.PHONY: build-container-dev
build-container-dev:
	$(MAKE) _build \
		IMG="$(REPO_URL)-dev" \
		DOCKERFILE="Dockerfile_dev" \
		TAG="$(BUILD_TAG)"

.PHONY: build-container-arm64
build-container-arm64:
	$(MAKE) _build \
		IMG="$(REPO_URL)-arm64" \
		DOCKERFILE="Dockerfile_arm" \
		TAG="$(BUILD_TAG)"

.PHONY: build-container-amd64
build-container-amd64:
	$(MAKE) _build \
		IMG="$(REPO_URL)-amd64" \
		DOCKERFILE="Dockerfile" \
		TAG="$(BUILD_TAG)"

.PHONY: build-container
build-container:
	$(MAKE) _build \
		IMG="$(REPO_URL)" \
		DOCKERFILE="Dockerfile" \
		TAG="$(BUILD_TAG)"

.PHONY: build-all-containers
build-all-containers: build-container build-container-amd64 build-container-arm64