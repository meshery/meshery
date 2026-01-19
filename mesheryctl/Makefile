###########################################################################
# Copyright Meshery Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
###########################################################################

#-----------------------------------------------------------------------------
# Includes
#-----------------------------------------------------------------------------
include ../install/Makefile.core.mk
include ../install/Makefile.show-help.mk

#-----------------------------------------------------------------------------
# Variables
#-----------------------------------------------------------------------------
GIT_VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_COMMITSHA := $(shell git rev-list -1 HEAD)
RELEASE_CHANNEL ?= edge
OUTDIR ?= bin
ARCH ?= amd64
GCFLAGS ?=
BINARY := mesheryctl
BINARY_WINDOWS := $(BINARY).exe

BINNAME_DARWIN ?= $(BINARY)-darwin-$(ARCH)
BINNAME_LINUX ?= $(BINARY)-linux-$(ARCH)
BINNAME_WINDOWS ?= $(BINARY)-windows-$(ARCH).exe

LDFLAGS := "-w -s -X github.com/meshery/meshery/mesheryctl/internal/cli/root/constants.version=${GIT_VERSION} -X github.com/meshery/meshery/mesheryctl/internal/cli/root/constants.commitsha=${GIT_COMMITSHA} -X github.com/meshery/meshery/mesheryctl/internal/cli/root/constants.releasechannel=${RELEASE_CHANNEL}"

BATS_LIB_PATH ?= helpers/bats-libs
BATS_FORMATTER ?= pretty
BATS_FOLDER_PATTERN ?= *-*
BATS_FILE_PATTERN ?= *-*
BATS_TEMP_DATA_DIR ?= $(shell mktemp -d --suffix="-bats")

#-----------------------------------------------------------------------------
# Build
#-----------------------------------------------------------------------------
.PHONY: build build-darwin build-linux build-windows build-all fmt clean forwin

## Format and build mesheryctl for the host platform.
# Default target: build mesheryctl for the host platform.
.DEFAULT_GOAL := build

build: fmt
	go build -ldflags=$(LDFLAGS) -o $(BINARY) cmd/mesheryctl/main.go

## Format source files with gofmt.
fmt:
	gofmt -l -s -w .

## Build mesheryctl for macOS (Darwin).
build-darwin:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=darwin GO111MODULE=on go build -ldflags=$(LDFLAGS) -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_DARWIN) ./cmd/mesheryctl/main.go

## Build mesheryctl for Linux.
build-linux:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=linux GO111MODULE=on go build -ldflags=$(LDFLAGS) -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_LINUX) ./cmd/mesheryctl/main.go

## Build mesheryctl for Windows.
build-windows:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=windows GO111MODULE=on go build -ldflags=$(LDFLAGS) -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_WINDOWS) ./cmd/mesheryctl/main.go

## Build mesheryctl binaries for all supported platforms.
build-all: build-darwin build-linux build-windows

## Legacy alias for Windows build to preserve existing workflows.
forwin: build-windows

## Remove compiled binaries and artifacts.
clean:
	rm -rf $(OUTDIR) $(BINARY) $(BINARY_WINDOWS) coverage.out coverage.html

#-----------------------------------------------------------------------------
# Error Codes
#-----------------------------------------------------------------------------
.PHONY: error error-util
## Analyze mesheryctl error codes.
error:
	go run github.com/meshery/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers

## Update mesheryctl error codes.
error-util:
	go run github.com/meshery/meshkit/cmd/errorutil -d . update -i ./helpers/ -o ./helpers

#-----------------------------------------------------------------------------
# Quality
#-----------------------------------------------------------------------------
.PHONY: lint coverage
## Lint mesheryctl using golangci-lint.
lint:
	golangci-lint run --timeout 5m

## Run tests with coverage and generate HTML report.
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
## Run all Mesheryctl integration tests (Golang)
tests-int:
	go test ./...
#-----------------------------------------------------------------------------
# Documentation
#-----------------------------------------------------------------------------
.PHONY: docs
## Generate markdown documentation for mesheryctl commands.
docs:
	cd doc; go run doc.go

#-----------------------------------------------------------------------------
# End-to-End Tests
#-----------------------------------------------------------------------------
.PHONY: e2e-libs e2e-clean e2e e2e-no-build e2e-run-all e2e-ci
## Fetch bats helper libraries for e2e tests.
e2e-libs:
	@mkdir -p tests/e2e/$(BATS_LIB_PATH) && \
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

## Clean temporary artifacts from e2e test runs.
e2e-clean:
	@echo "\n\n#### MESHERYCTL END-TO-END TESTING - CLEANUP ####"; \
	rm -rf /tmp/tmp.*-bats; \
	echo "#### MESHERYCTL END-TO-END TESTING - CLEANUP DONE ####\n\n"

## Build mesheryctl and run end-to-end tests.
e2e: build e2e-libs
	@echo "\n\n#### MESHERYCTL END-TO-END TESTING - START ####\n\n"; \
	cd tests/e2e; \
	BATS_LIB_PATH="$(shell pwd)/tests/e2e/$(BATS_LIB_PATH)" \
	E2E_HELPERS_PATH=$(shell pwd)/tests/e2e/helpers \
	E2E_TESTDATA_PATH=$(shell pwd)/tests/e2e/testdata \
	TEMP_DATA_DIR=$(shell mktemp -d) \
	MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml" \
	MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json" \
	MESHERYCTL_BIN="../../$(BINARY)" bats --formatter $(BATS_FORMATTER) $(BATS_FOLDER_PATTERN)/$(BATS_FILE_PATTERN).bats; \
	echo "\n\n#### MESHERYCTL END-TO-END TESTING - DONE ####\n\n"

## Build mesheryctl and run end-to-end tests.
e2e-ci: build e2e-libs
	@echo "\n\n#### MESHERYCTL END-TO-END TESTING - START ####\n\n"; \
	cd tests/e2e; \
	BATS_LIB_PATH="$(shell pwd)/tests/e2e/$(BATS_LIB_PATH)" \
	E2E_HELPERS_PATH=$(shell pwd)/tests/e2e/helpers \
	E2E_TESTDATA_PATH=$(shell pwd)/tests/e2e/testdata \
	TEMP_DATA_DIR=$(shell mktemp -d) \
	MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml" \
	MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json" \
	MESHERYCTL_BIN="../../$(BINARY)" bats  $(BATS_FOLDER_PATTERN)/$(BATS_FILE_PATTERN).bats >> bats.tap; cd ..
	echo "\n\n#### MESHERYCTL END-TO-END TESTING - DONE ####\n\n"  

	ALLURE_LABELS="epic=MesheryCtl,type=e2e,env=ci" node bats-to-allure.js tests/e2e/bats.tap

## Run end-to-end tests using an existing mesheryctl binary.
e2e-no-build: e2e-libs
	@echo "\n\n#### MESHERYCTL END-TO-END TESTING - START ####\n\n"; \
	cd tests/e2e; \
	BATS_LIB_PATH="$(shell pwd)/tests/e2e/$(BATS_LIB_PATH)" \
	E2E_HELPERS_PATH=$(shell pwd)/tests/e2e/helpers \
	E2E_TESTDATA_PATH=$(shell pwd)/tests/e2e/testdata \
	TEMP_DATA_DIR=$(BATS_TEMP_DATA_DIR) \
	MESHERY_CONFIG_FILE_PATH="${HOME}/.meshery/config.yaml" \
	MESHERY_AUTH_FILE="${HOME}/.meshery/auth.json" \
	MESHERYCTL_BIN="../../$(BINARY)" bats --formatter $(BATS_FORMATTER) $(BATS_FOLDER_PATTERN)/$(BATS_FILE_PATTERN).bats; \
	echo "\n\n#### MESHERYCTL END-TO-END TESTING - DONE ####\n\n"

## Run end-to-end tests and clean temporary artifacts.
e2e-run-all: e2e-no-build e2e-clean
