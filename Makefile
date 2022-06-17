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

include install/Makefile.core.mk
include install/Makefile.show-help.mk

#-----------------------------------------------------------------------------
# Docker-based Builds
#-----------------------------------------------------------------------------
.PHONY: docker-build docker-local-cloud docker-cloud
## Build Meshery Server and UI containers.
docker-build:
	# `make docker` builds Meshery inside of a multi-stage Docker container.
	# This method does NOT require that you have Go, NPM, etc. installed locally.
	DOCKER_BUILDKIT=1 docker build -f install/docker/Dockerfile -t layer5/meshery --build-arg TOKEN=$(GLOBAL_TOKEN) --build-arg GIT_COMMITSHA=$(GIT_COMMITSHA) --build-arg GIT_VERSION=$(GIT_VERSION) --build-arg RELEASE_CHANNEL=${RELEASE_CHANNEL} .


## Meshery Cloud for user authentication.
## Runs Meshery in a container locally and points to locally-running
docker-local-cloud:
	
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link meshery-cloud:meshery-cloud \
	-e PROVIDER_BASE_URLS=$(REMOTE_PROVIDER_LOCAL) \
	-e DEBUG=true \
	-e ADAPTER_URLS=$(ADAPTER_URLS) \
	-p 9081:8080 \
	layer5/meshery ./meshery

## Runs Meshery in a container locally and points to remote
## Remote Provider for user authentication.
docker-cloud:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	-e PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	-e DEBUG=true \
	-e ADAPTER_URLS=$(ADAPTER_URLS) \
	-v meshery-config:/home/appuser/.meshery/config \
  -v $(HOME)/.kube:/home/appuser/.kube:ro \
	-p 9081:8080 \
	layer5/meshery ./meshery

#-----------------------------------------------------------------------------
# Meshery Server Native Builds
#-----------------------------------------------------------------------------
.PHONY: server wrk2-setup nighthawk-setup server-local server-skip-compgen server-no-content golangci proto-build error
## Setup wrk2 for local development.
wrk2-setup:
	echo "setup-wrk does not work on Mac Catalina at the moment"
	cd cmd; git clone https://github.com/layer5io/wrk2.git; cd wrk2; make; cd ..

## ## Setup nighthawk for local development.
nighthawk-setup:
	cd cmd; git clone https://github.com/layer5io/nighthawk-go.git; cd nighthawk-go; make setup; cd ..

run-local: server-local error
## Build and run Meshery Server on your local machine
## and point to (expect) a locally running Meshery Cloud or other Provider(s)
## for user authentication (requires go${GOVERSION}).
server-local:
	cd cmd; go$(GOVERSION) clean; rm meshery; go$(GOVERSION) mod tidy; \
	go$(GOVERSION) build -ldflags="-w -s -X main.version=${GIT_VERSION} -X main.commitsha=${GIT_COMMITSHA} -X main.releasechannel=${RELEASE_CHANNEL}" -tags draft -a -o meshery; \
	PROVIDER_BASE_URLS=$(REMOTE_PROVIDER_LOCAL) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	./meshery; \
	cd ..

run-fast: 
	## "DEPRECATED: This target is deprecated. Use `make server`.

## Build and run Meshery Server on your local machine (requires go${GOVERSION}).
server:
	cd cmd; go$(GOVERSION) mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	go$(GOVERSION) run main.go;

## Build and run Meshery Server with no Kubernetes components on your local machine (requires go${GOVERSION}).
server-skip-compgen:
	cd cmd; go$(GOVERSION) mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
 	SKIP_COMP_GEN=true \
	go$(GOVERSION) run main.go;
		
## Build and run Meshery Server with no seed content (requires go$(GOVERSION)).
server-no-content:
	cd cmd; go$(GOVERSION) mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	SKIP_DOWNLOAD_CONTENT=true \
	go$(GOVERSION) run main.go;

## Lint check Meshery Server.
golangci: error
	GO111MODULE=off GOPROXY=direct GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2;
	$(GOPATH)/bin/golangci-lint run

## Build Meshery's protobufs (requires go$(GOVERSION)).
proto-build:
	# see https://grpc.io/docs/languages/go/quickstart/
	# go$(GOVERSION) get -u google.golang.org/grpc
	# go$(GOVERSION) get -u google.golang.org/protobuf/cmd/protoc-gen-go \
	#         google.golang.org/grpc/cmd/protoc-gen-go-grpc
	# PATH=$(PATH):`pwd`/../protoc/bin:$(GOPATH)/bin
	# export PATH=$PATH:`pwd`/../protoc/bin:$GOPATH/bin
	protoc -I meshes/ meshes/meshops.proto --go-grpc_out=./meshes/ --go_out=./meshes/

## Analyze error codes
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers --skip-dirs mesheryctl

#-----------------------------------------------------------------------------
# Meshery UI Native Builds
#-----------------------------------------------------------------------------
.PHONY: setup-ui-libs ui-setup run-ui-dev ui ui-meshery-build ui ui-provider ui-lint ui-provider ui-meshery ui-build ui-provider-build ui-provider-test

setup-ui-libs: ui-setup
## Install dependencies for building Meshery UI.
ui-setup:
	cd ui; npm i; cd ..
	cd provider-ui; npm i; cd ..

run-ui-dev: ui
## Run Meshery UI on your local machine. Listen for changes.
ui:
	cd ui; npm run dev; cd ..

run-provider-ui-dev: ui-provider
## Run Meshery Provider UI  on your local machine. Listen for changes.
ui-provider:
	cd provider-ui; npm run dev; cd ..

lint-ui: ui-lint
## Lint check Meshery UI and Provider UI on your local machine.
ui-lint:
	cd ui; npm run lint; cd ..

lint-provider-ui: ui-provider-lint
## Lint check Meshery Provider UI on your local machine.
ui-provider-lint:
	cd provider-ui; npm run lint; cd ..

## Test Meshery Provider UI on your local machine.
ui-provider-test:
	cd provider-ui; npm run test; cd ..

build-ui: ui-build
## Buils all Meshery UIs  on your local machine.
ui-build: 
	cd ui; npm run build && npm run export; cd ..
	cd provider-ui; npm run build && npm run export; cd ..

build-meshery-ui: ui-meshery-build
## Build only Meshery UI on your local machine.
ui-meshery-build:
	cd ui; npm run build && npm run export; cd ..

build-provider-ui: ui-provider-build
## Builds only the provider user interface on your local machine
ui-provider-build:
	cd provider-ui; npm run build && npm run export; cd ..

run-ui-integration-tests: ui-integration-tests
## Run Meshery Cypress Integration Tests against your local Meshery UI (cypress runs in non-interactive mode).
ui-integration-tests: ui-setup
	cd ui; npm run ci-test-integration; cd ..

#-----------------------------------------------------------------------------
# Meshery Docs
#-----------------------------------------------------------------------------
#Incorporating Make docs commands from the Docs Makefile
.PHONY: docs docs-build site docs-docker
jekyll=bundle exec jekyll

site: docs

## Run Meshery Docs. Listen for changes.
docs:
	cd docs; bundle install; $(jekyll) serve --drafts --livereload --config _config_dev.yml

## Build Meshery Docs on your local machine.
docs-build:
	cd docs; $(jekyll) build --drafts

## Run Meshery Docs in a Docker container. Listen for changes.
docs-docker:
	cd docs; docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:4.0.0 bash -c "bundle install; jekyll serve --drafts --livereload"

#-----------------------------------------------------------------------------
# Meshery Helm Charts
#-----------------------------------------------------------------------------
.PHONY: helm-docs helm-operator-docs helm-meshery-docs helm-operator-lint helm-lint
## Generate all Meshery Helm Chart documentation in markdown format.
helm-docs: helm-operator-docs helm-meshery-docs

## Generate Meshery Operator Helm Chart documentation in markdown format.
helm-operator-docs:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs 
	$(GOPATH)/bin/helm-docs -c install/kubernetes/helm/meshery-operator

## Generate Meshery Server and Adapters Helm Chart documentation in markdown format.
helm-meshery-docs:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs 
	$(GOPATH)/bin/helm-docs -c install/kubernetes/helm/meshery

## Lint all of Meshery's Helm Charts
helm-lint: helm-operator-lint helm-meshery-lint

## Lint Meshery Operator Helm Chart
helm-operator-lint:
	helm lint install/kubernetes/helm/meshery-operator --with-subcharts
## Lint Meshery Server and Adapter Helm Charts
helm-meshery-lint:
	helm lint install/kubernetes/helm/meshery --with-subcharts

#-----------------------------------------------------------------------------
# Meshery APIs
#-----------------------------------------------------------------------------
.PHONY: swagger-build swagger swagger-docs-build graphql-docs graphql-build
## Build Meshery REST API specifications
swagger-build:
	swagger generate spec -o ./helpers/swagger.yaml --scan-models

## Generate and serve Meshery REST API specifications
swagger: swagger-build
	swagger serve ./helpers/swagger.yaml

## Build Meshery REST API documentation
swagger-docs-build:
	swagger generate spec -o ./docs/_data/swagger.yml --scan-models; \
	swagger flatten ./docs/_data/swagger.yml -o ./docs/_data/swagger.yml --with-expand --format=yaml

## Build Meshery GraphQL API documentation
graphql-docs:
	cd docs; build-docs; bundle exec rake graphql:compile_docs

## Build Meshery GraphQl API specifications
graphql-build:
	cd internal/graphql; go run -mod=mod github.com/99designs/gqlgen generate

