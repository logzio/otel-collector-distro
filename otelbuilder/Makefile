BINARY_NAME ?= otelcol-logzio
BUILDER_VERSION ?= 0.54.0
BUILDER_REPO ?= github.com/open-telemetry/opentelemetry-collector
BUILDER_BIN_NAME ?= opentelemetry-collector-builder
VERSION ?= 0.0.1

GO ?= go
OS ?= $(shell uname -s | tr A-Z a-z)

ifeq ($(shell go env GOOS),darwin)
CGO_ENABLED ?= 1
else
CGO_ENABLED ?= 0
endif

.PHONY: _install-bin
_install-bin:
	@curl -L -o $(BUILDER_BIN_NAME) https://$(BUILDER_REPO)/releases/download/v$(BUILDER_VERSION)/ocb_$(BUILDER_VERSION)_$(PLATFORM)_amd64
	@chmod +x $(BUILDER_BIN_NAME)

# install otelcol builder
.PHONY: install-builder
install-builder:
	@echo "Installing $(BUILDER_REPO)/cmd/builder@v$(BUILDER_VERSION)... (in $(BUILDER_BIN_PATH))"
	@$(MAKE) _install-bin PLATFORM=$(OS)


.PHONY: _builder
_builder:
# Need to specify go path because otherwise opentelemetry-collector-builder
# uses /usr/bin/go which on Github Actions is using preinstalled 1.15.12 by default.
	CGO_ENABLED=$(CGO_ENABLED) ./$(BUILDER_BIN_NAME) \
		--go $(GO) \
		--version "$(VERSION)" \
		--config otelcol-builder.yaml \
		--output-path ./cmd \
		--skip-compilation=$(SKIP_COMPILATION)

.PHONY: _gobuild
_gobuild:
	(cd cmd && \
		CGO_ENABLED=$(CGO_ENABLED) go build -v \
		-tags enable_unstable \
		-ldflags="-s -w" \
		-trimpath \
		-o ./$(BINARY_NAME) . \
	)
	chmod +x ./cmd/$(BINARY_NAME)

.PHONY: _gobuild_debug
_gobuild_debug:
	(cd cmd && \
		CGO_ENABLED=$(CGO_ENABLED) go build -v \
		-tags enable_unstable \
		-race \
		-gcflags "all=-N -l" \
		-o ./$(BINARY_NAME)-debug . \
	)
	chmod +x ./cmd/$(BINARY_NAME)-debug

.PHONY: build
build:
	@$(MAKE) generate-sources
	@$(MAKE) _gobuild

.PHONY: build-debug
build-debug:
	@$(MAKE) generate-sources
	@$(MAKE) _gobuild_debug

.PHONY: generate-sources
generate-sources:
	@$(MAKE) _builder SKIP_COMPILATION=true

.PHONY: test
test:
	@$(MAKE) install-builder
	@$(MAKE) generate-sources
	@$(MAKE) -C cmd test

.PHONY: lint
lint: install-builder generate-sources
	@$(MAKE) -C cmd lint

# build collector distro binaries for all supported os arch combinations
.PHONY: otelcol-logzio-all
otelcol-logzio-all-sys: otelcol-logzio-darwin_amd64 otelcol-logzio-linux_amd64 otelcol-logzio-windows_amd64 otelcol-logzio-linux_arm64 otelcol-logzio-linux_arm

.PHONY: otelcol-logzio-darwin_amd64
otelcol-logzio-darwin_amd64:
	GOOS=darwin  GOARCH=amd64 $(MAKE) build BINARY_NAME=$(BINARY_NAME)-darwin_amd64

.PHONY: otelcol-logzio-linux_amd64
otelcol-logzio-linux_amd64:
	GOOS=linux   GOARCH=amd64 $(MAKE) build BINARY_NAME=$(BINARY_NAME)-linux_amd64

.PHONY: otelcol-logzio-linux_arm64
otelcol-logzio-linux_arm64:
	GOOS=linux   GOARCH=arm64 $(MAKE) build BINARY_NAME=$(BINARY_NAME)-linux_arm64

.PHONY: otelcol-logzio-linux_arm
otelcol-logzio-linux_arm:
	GOOS=linux   GOARCH=arm $(MAKE) build BINARY_NAME=$(BINARY_NAME)-linux_arm

.PHONY: otelcol-logzio-windows_amd64
otelcol-logzio-windows_amd64:
	GOOS=windows GOARCH=amd64 $(MAKE) build BINARY_NAME=$(BINARY_NAME)-windows_amd64.exe