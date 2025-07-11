# ENVs
GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_COMMITSHA=$(shell git rev-list -1 HEAD)
RELEASE_CHANNEL="edge"
LDFLAGS=-w -s -X github.com/meshery/meshery/mesheryctl/internal/cli/root/constants.version=${GIT_VERSION} -X github.com/meshery/meshery/mesheryctl/internal/cli/root/constants.commitsha=${GIT_COMMITSHA} -X github.com/meshery/meshery/mesheryctl/internal/cli/root/constants.releasechannel=${RELEASE_CHANNEL}

make:
	gofmt -l -s -w .
	go build -ldflags="${LDFLAGS}" -o mesheryctl cmd/mesheryctl/main.go
forwin:
	go build -ldflags="${LDFLAGS}" -o mesheryctl.exe cmd/mesheryctl/main.go

OUTDIR := bin

ARCH := amd64

BINNAME_DARWIN ?= mesheryctl-darwin-$(ARCH)
BINNAME_LINUX ?= mesheryctl-linux-$(ARCH)
BINNAME_WINDOWS ?= mesheryctl-windows-$(ARCH).exe

BATS_LIB_PATH ?= "helpers/bats-libs"
BATS_FORMATTER ?= "pretty"
E2E_TEMP_DATA_DIR ?= $(shell mktemp -d --suffix="-bats")

# CLI

.PHONY: $(BINNAME_DARWIN)
$(BINNAME_DARWIN):
	@echo
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=darwin GO111MODULE=on go build -ldflags="${LDFLAGS}" -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_DARWIN) ./cmd/mesheryctl/main.go

.PHONY: $(BINNAME_LINUX)
$(BINNAME_LINUX):
	@echo
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=linux GO111MODULE=on go build -ldflags="${LDFLAGS}" -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_LINUX) ./cmd/mesheryctl/main.go

.PHONY: $(BINNAME_WINDOWS)
$(BINNAME_WINDOWS):
	@echo
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=windows GO111MODULE=on go build -ldflags="${LDFLAGS}" -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_WINDOWS) ./cmd/mesheryctl/main.go

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

.PHONY: error
error:
	go run github.com/meshery/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

## Runs meshkit error utility to update error codes for mesheryctl.
error-util:
	go run github.com/meshery/meshkit/cmd/errorutil -d . update -i ./helpers/ -o ./helpers

# Check code quality at local by using golangci-lint and return the coverage
.PHONY: lint
lint:
	golangci-lint run --timeout 5m

# output a HTML coverage by using gocov and gocov-html
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Generate markdown for mesheryctl docs
.PHONY: docs
docs:
	cd doc; go run doc.go


e2e-libs:
	@ mkdir -p tests/e2e/$(BATS_LIB_PATH) && \
	cd tests/e2e/$(BATS_LIB_PATH) && \
	if [ ! -d "bats-assert" ]; then \
		git clone https://github.com/bats-core/bats-assert.git; \
	fi && \
	if [ ! -d "bats-detik" ]; then \
		git clone https://github.com/bats-core/bats-detik.git; \
	fi && \
	if [ ! -d "bats-support" ]; then \
		git clone https://github.com/bats-core/bats-support.git; \
	fi && \
	if [ ! -d "bats-file" ]; then \
		git clone https://github.com/bats-core/bats-file.git; \
	fi


e2e-clean:
	@ echo "\n\n#### MESHERYCTL END-TO-END TESTING - CLEANUP ####"; \
	rm -rf /tmp/tmp.*-bats; \
	echo "#### MESHERYCTL END-TO-END TESTING - CLEANUP DONE ####\n\n"

.PHONY: e2e
e2e: make e2e-libs
	@ echo "\n\n#### MESHERYCTL END-TO-END TESTING - START ####\n\n"; \
	cd tests/e2e; \
	BATS_LIB_PATH="$(shell pwd)/tests/e2e/$(BATS_LIB_PATH)" \
	E2E_HELPERS_PATH=$(shell pwd)/tests/e2e/helpers \
	E2E_TESTDATA_PATH=$(shell pwd)/tests/e2e/testdata \
	TEMP_DATA_DIR=$(shell mktemp -d) \
	MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml" \
	MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json" \
	MESHERYCTL_BIN="../../mesheryctl" bats --formatter $(BATS_FORMATTER) *-*/*.bats; \
	echo "\n\n#### MESHERYCTL END-TO-END TESTING - DONE ####\n\n"


.PHONY: e2e-no-build
e2e-no-build: e2e-libs
	@ echo "\n\n#### MESHERYCTL END-TO-END TESTING - START ####\n\n"; \
	cd tests/e2e; \
	BATS_LIB_PATH="$(shell pwd)/tests/e2e/$(BATS_LIB_PATH)" \
	E2E_HELPERS_PATH=$(shell pwd)/tests/e2e/helpers \
	E2E_TESTDATA_PATH=$(shell pwd)/tests/e2e/testdata \
	TEMP_DATA_DIR=$(E2E_TEMP_DATA_DIR) \
	MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml" \
	MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json" \
	MESHERYCTL_BIN="../../mesheryctl" bats --formatter $(BATS_FORMATTER) *-*/*.bats; \
	echo "\n\n#### MESHERYCTL END-TO-END TESTING - DONE ####\n\n"

e2e-run-all: e2e-no-build e2e-clean