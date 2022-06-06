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

.PHONY: test
test:
	$(GOTEST) $(GOTEST_OPT) ./...

.PHONY: tidy
tidy:
	rm -fr go.sum
	$(GOCMD) mod tidy -compat=1.18

