# SRC_ROOT is the top of the source tree.
SRC_ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

# build tags required by any component should be defined as an independent variables and later added to GO_BUILD_TAGS below
GO_BUILD_TAGS=""
GOTEST_OPT?= -race -v -timeout 300s --tags=$(GO_BUILD_TAGS)
GOTEST_INTEGRATION_OPT?= -race -timeout 120s
GOTEST_OPT_WITH_COVERAGE = $(GOTEST_OPT) -coverprofile=coverage.txt -covermode=atomic
GOTEST_OPT_WITH_INTEGRATION=$(GOTEST_INTEGRATION_OPT) -v -tags=integration,$(GO_BUILD_TAGS) -run=Integration -coverprofile=integration-coverage.txt -covermode=atomic
GOCMD?= go
GOTEST=$(GOCMD) test
GOOS=$(shell $(GOCMD) env GOOS)
GOARCH=$(shell $(GOCMD) env GOARCH)
ADDLICENCESE= addlicense
MISSPELL=misspell -error
MISSPELL_CORRECTION=misspell -w
LINT=golangci-lint
IMPI=impi
# BUILD_TYPE should be one of (dev, release).
BUILD_TYPE?=release
ALL_PKG_DIRS := $(shell $(GOCMD) list -f '{{ .Dir }}' ./... | sort)

ALL_SRC := $(shell find $(ALL_PKG_DIRS) -name '*.go' \
                                -not -path '*/third_party/*' \
                                -type f | sort)

# All source code and documents. Used in spell check.
ALL_SRC_AND_DOC := $(shell find $(ALL_PKG_DIRS) -name "*.md" -o -name "*.go" -o -name "*.yaml" \
                                -not -path '*/third_party/*' \
                                -type f | sort)

# ALL_PKGS is used with 'go cover'
ALL_PKGS := $(shell $(GOCMD) list $(sort $(dir $(ALL_SRC))))

.PHONY: common
common: checklicense impi lint misspell

.PHONY: test
test:
	$(GOTEST) $(GOTEST_OPT) ./...

.PHONY: tidy
tidy:
	rm -fr go.sum
	$(GOCMD) mod tidy -compat=1.17

.PHONY: misspell
misspell:
	@echo "running $(MISSPELL)"
	@$(MISSPELL) $(ALL_SRC_AND_DOC)

.PHONY: misspell-correction
misspell-correction:
	$(MISSPELL_CORRECTION) $(ALL_SRC_AND_DOC)

.PHONY: lint
lint: impi misspell
	$(LINT) run --allow-parallel-runners

.PHONY: impi
impi:
	@$(IMPI) --local github.com/logzio/otel-collector-distro --scheme stdThirdPartyLocal $(ALL_PKGS)

.PHONY: checklicense
checklicense:
	@ADDLICENCESEOUT=`$(ADDLICENCESE) -check $(ALL_SRC) 2>&1`; \
		if [ "$$ADDLICENCESEOUT" ]; then \
			echo "$(ADDLICENCESE) FAILED => add License errors:\n"; \
			echo "$$ADDLICENCESEOUT\n"; \
			echo "Use 'make addlicense' to fix this."; \
			exit 1; \
		else \
			echo "Check License finished successfully"; \
		fi

.PHONY: addlicense
addlicense:
	@ADDLICENCESEOUT=`$(ADDLICENCESE) -y "" -c 'The OpenTelemetry Authors' $(ALL_SRC) 2>&1`; \
		if [ "$$ADDLICENCESEOUT" ]; then \
			echo "$(ADDLICENCESE) FAILED => add License errors:\n"; \
			echo "$$ADDLICENCESEOUT\n"; \
			exit 1; \
		else \
			echo "Add License finished successfully"; \
		fi


.PHONY: fmt
fmt:
	gofmt  -w -s ./
	goimports -w  -local github.com/logzio/otel-collector-distro ./