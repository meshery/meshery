# Copyright Layer5, Inc.
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


## Install docs.layer5.io dependencies on your local machine.
## See https://gohugo.io/categories/installation
setup:
	npm install

## Run docs.meshery.io on your local machine with draft and future content enabled.
site: check-go
	hugo server -D -F
	
## Run docs.meshery.io on your local machine. Alternate method.
site-fast:
	gatsby develop

## Build docs.meshery.io on your local machine.
build:
	hugo

## Empty build cache and run docs.meshery.io on your local machine.
clean: 
	hugo --cleanDestinationDir
	make site

check-go:
	@echo "Checking if Go is installed..."
	@command -v go > /dev/null || (echo "Go is not installed. Please install it before proceeding."; exit 1)
	@echo "Go is installed."

docker-install-base-deps:
	apt-get update
	apt-get install -y --no-install-recommends nodejs npm ca-certificates wget xz-utils
	rm -rf /var/lib/apt/lists/*

docker-deps:
	npm ci

docker-download-hugo:
	@test -n "$(HUGO_VERSION)" || (echo "HUGO_VERSION is required"; exit 1)
	@test -n "$(TARGETARCH)" || (echo "TARGETARCH is required"; exit 1)
	mkdir -p /tmp/hugo
	wget -O /tmp/hugo/hugo.tar.gz "https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VERSION)/hugo_extended_$(HUGO_VERSION)_linux-$(TARGETARCH).tar.gz"
	tar -xf /tmp/hugo/hugo.tar.gz -C /tmp/hugo hugo

docker-build:
	set -e; \
	args="--gc --minify -d /out"; \
	if [ -n "$(HUGO_ENV)" ]; then args="$$args -e $(HUGO_ENV)"; fi; \
	if [ -n "$(DOCS_URL)" ]; then args="$$args -b $(DOCS_URL)"; fi; \
	hugo $$args

docker-validate-stats:
	@test -f /src/hugo_stats.json || (echo "hugo_stats.json was not generated"; exit 1)

docker-lint:
	npx --yes markdownlint-cli2 \
		"content/**/*.md" \
		"#content/engine/release-notes/*.md" \
		"#content/desktop/previous-versions/*.md"

docker-install-htmltest:
	@test -n "$(HTMLTEST_VERSION)" || (echo "HTMLTEST_VERSION is required"; exit 1)
	GOBIN=/usr/local/bin go install github.com/wjdp/htmltest@v$(HTMLTEST_VERSION)

docker-htmltest:
	@test -n "$(HTMLTEST_DIR)" || (echo "HTMLTEST_DIR is required"; exit 1)
	@if [ -f /src/.htmltest.yml ]; then \
		cd "$(HTMLTEST_DIR)" && htmltest --conf /src/.htmltest.yml; \
	else \
		cd "$(HTMLTEST_DIR)" && htmltest; \
	fi

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

docker-build-upstream:
	@test -n "$(UPSTREAM_MODULE_NAME)" || (echo "UPSTREAM_MODULE_NAME is required"; exit 1)
	hugo --ignoreVendorPaths "github.com/$(UPSTREAM_MODULE_NAME)" -d /out

## Build and run docs website within a Docker container
docker:
	docker compose watch

.PHONY: \
	setup \
	build \
	site \
	clean \
	site-fast \
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