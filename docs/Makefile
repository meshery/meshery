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


## Install docs.meshery.io dependencies on your local machine.
## See https://gohugo.io/categories/installation
setup:
	npm install

## Verify required commands and local dependencies are present.
check-deps:
	@echo "Checking if 'npx' and local 'hugo' binary are present..."
	@command -v npx > /dev/null || { echo "Error: 'npx' not found. Please install Node.js and npm."; exit 1; }
	@test -x node_modules/.bin/hugo || { echo "Error: Hugo binary not found in node_modules. Please run 'make setup' first."; exit 1; }
	@echo "Dependencies check passed."

## Run docs.meshery.io on your local machine with draft and future content enabled.
site: check-deps check-go
	npx hugo server -D -F

## Run docs.meshery.io on your local machine in serve mode (without file watching).
serve: check-deps check-go
	npx hugo server -D -F --watch=false

## Build docs.meshery.io on your local machine.
build: check-deps check-go
	npx hugo

## Empty build cache and run docs.meshery.io on your local machine.
clean: check-deps
	npx hugo --cleanDestinationDir
	make site

## Validate Go is installed by delegating to the root Makefile check.
check-go:
	@$(MAKE) -C .. check-go

## Install base OS dependencies needed in Docker-based docs builds.
docker-install-base-deps:
	apt-get update
	apt-get install -y --no-install-recommends nodejs npm ca-certificates wget xz-utils
	rm -rf /var/lib/apt/lists/*

## Install Node.js dependencies in CI or container contexts.
docker-deps:
	npm ci

## Download and unpack a specific Hugo Extended binary for Docker builds.
docker-download-hugo:
	@test -n "$(HUGO_VERSION)" || (echo "HUGO_VERSION is required"; exit 1)
	@test -n "$(TARGETARCH)" || (echo "TARGETARCH is required"; exit 1)
	mkdir -p /tmp/hugo
	wget -O /tmp/hugo/hugo.tar.gz "https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VERSION)/hugo_extended_$(HUGO_VERSION)_linux-$(TARGETARCH).tar.gz"
	tar -xf /tmp/hugo/hugo.tar.gz -C /tmp/hugo hugo

## Build the docs site for Docker output with optional environment/base URL overrides.
docker-build:
	set -e; \
	args="--gc --minify -d /out"; \
	if [ -n "$(HUGO_ENV)" ]; then args="$$args -e $(HUGO_ENV)"; fi; \
	if [ -n "$(DOCS_URL)" ]; then args="$$args -b $(DOCS_URL)"; fi; \
	hugo $$args

## Ensure Hugo generated the stats file required by downstream tooling.
docker-validate-stats:
	@test -f /src/hugo_stats.json || (echo "hugo_stats.json was not generated"; exit 1)

## Lint Markdown content during Docker-based checks.
docker-lint:
	npx --yes markdownlint-cli2 \
		"content/**/*.md" \
		"#content/engine/release-notes/*.md" \
		"#content/desktop/previous-versions/*.md"

## Install htmltest binary for link validation in container environments.
docker-install-htmltest:
	@test -n "$(HTMLTEST_VERSION)" || (echo "HTMLTEST_VERSION is required"; exit 1)
	GOBIN=/usr/local/bin go install github.com/wjdp/htmltest@v$(HTMLTEST_VERSION)

## Run htmltest against the generated site output.
docker-htmltest:
	@test -n "$(HTMLTEST_DIR)" || (echo "HTMLTEST_DIR is required"; exit 1)
	@if [ -f /src/.htmltest.yml ]; then \
		cd "$(HTMLTEST_DIR)" && htmltest --conf /src/.htmltest.yml; \
	else \
		cd "$(HTMLTEST_DIR)" && htmltest; \
	fi

## Update Hugo modules and vendor dependencies for Docker builds.
docker-update-modules:
	set -ex; \
	if [ -n "$(MODULE)" ]; then \
		module_base="$${MODULE%@*}"; \
		hugo mod get "$(MODULE)"; \
		resolved=$$(grep -m 1 "$$module_base" go.mod | awk '{print $$1 "@" $$2}'); \
		go mod edit -replace "$$module_base=$$resolved"; \
	else \
		echo "no module set"; \
	fi; \
	hugo mod vendor

## Build docs while ignoring vendored paths for an upstream module.
docker-build-upstream:
	@test -n "$(UPSTREAM_MODULE_NAME)" || (echo "UPSTREAM_MODULE_NAME is required"; exit 1)
	hugo --ignoreVendorPaths "github.com/$(UPSTREAM_MODULE_NAME)" -d /out

## Build and run docs website within a Docker container
docker:
	docker compose watch

.PHONY: \
	setup \
	check-deps \
	build \
	site \
	serve \
	clean \
	check-go \
	docker \
	docker-install-base-deps \
	docker-deps \
	docker-download-hugo \
	docker-build \
	docker-validate-stats \
	docker-lint \
	docker-install-htmltest \
	docker-htmltest \
	docker-update-modules \
	docker-build-upstream