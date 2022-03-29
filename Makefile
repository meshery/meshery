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
## Build Meshery Server and UI containers.
docker:
	# `make docker` builds Meshery inside of a multi-stage Docker container.
	# This method does NOT require that you have Go, NPM, etc. installed locally.
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery --build-arg TOKEN=$(GLOBAL_TOKEN) --build-arg GIT_COMMITSHA=$(GIT_COMMITSHA) --build-arg GIT_VERSION=$(GIT_VERSION) --build-arg RELEASE_CHANNEL=${RELEASE_CHANNEL} .


## Meshery Cloud for user authentication.
## Runs Meshery in a container locally and points to locally-running
docker-run-local-cloud:
	
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
docker-run-cloud:
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
## Setup wrk2 for local development.
setup-wrk2:
	echo "setup-wrk does not work on Mac Catalina at the moment"
	cd cmd; git clone https://github.com/layer5io/wrk2.git; cd wrk2; make; cd ..

## ## Setup nighthawk for local development.
setup-nighthawk:
	cd cmd; git clone https://github.com/layer5io/nighthawk-go.git; cd nighthawk-go; make setup; cd ..

## Run Meshery on your local machine and point to locally-running
##  Meshery Cloud for user authentication.
run-local-cloud: error
	cd cmd; go clean; rm meshery; go mod tidy; \
	go build -ldflags="-w -s -X main.version=${GIT_VERSION} -X main.commitsha=${GIT_COMMITSHA} -X main.releasechannel=${RELEASE_CHANNEL}" -tags draft -a -o meshery; \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_DEV) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	./meshery; \
	cd ..

## Build and run Meshery to run on your local machine
## and point to remote Meshery Cloud for user authentication.
run-local: error
	cd cmd; go clean; rm meshery; go mod tidy; \
	go build -ldflags="-w -s -X main.version=${GIT_VERSION} -X main.commitsha=${GIT_COMMITSHA} -X main.releasechannel=${RELEASE_CHANNEL}" -tags draft -a -o meshery; \
	PROVIDER_BASE_URLS=$(REMOTE_PROVIDER_LOCAL) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	./meshery; \
	cd ..

## Buiild and run Meshery Server on your local machine.
run-fast:
	cd cmd; go mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	go run main.go;

## Build and run Meshery Server with no Kubernetes components on your local machine.
run-fast-skip-compgen:
	cd cmd; go mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
 	SKIP_COMP_GEN=true \
	go run main.go;
		
## Build and run Meshery Server with no seed content.
run-fast-no-content:
	cd cmd; go mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	SKIP_DOWNLOAD_CONTENT=true \
	go run main.go;

## Lint check Meshery Server.
golangci-run: error
	GO111MODULE=off GOPROXY=direct GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2;
	$(GOPATH)/bin/golangci-lint run

## Build Meshery's protobufs.
proto:
	# see https://grpc.io/docs/languages/go/quickstart/
	# go get -u google.golang.org/grpc
	# go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
	#         google.golang.org/grpc/cmd/protoc-gen-go-grpc
	# PATH=$(PATH):`pwd`/../protoc/bin:$(GOPATH)/bin
	# export PATH=$PATH:`pwd`/../protoc/bin:$GOPATH/bin
	protoc -I meshes/ meshes/meshops.proto --go-grpc_out=./meshes/ --go_out=./meshes/

#-----------------------------------------------------------------------------
# Meshery UI Native Builds
#-----------------------------------------------------------------------------
## Install dependencies for building Meshery UI.
setup-ui-libs:
	cd ui; npm i; cd ..
	cd provider-ui; npm i; cd ..

## Run Meshery UI on your local machine. Listen for changes.
run-ui-dev:
	cd ui; npm run dev; cd ..

## Run Meshery Provider UI  on your local machine. Listen for changes.
run-provider-ui-dev:
	cd provider-ui; npm run dev; cd ..

## Lint check Meshery UI and Provider UI on your local machine.
lint-ui:
	cd ui; npm run lint; cd ..

## Lint check Meshery Provider UI on your local machine.
lint-provider-ui:
	cd provider-ui; npm run lint; cd ..

## Test Meshery Provider UI on your local machine.
test-provider-ui:
	cd provider-ui; npm run test; cd ..

## Buils all Meshery UIs  on your local machine.
build-ui:
	cd ui; npm run build && npm run export; cd ..
	cd provider-ui; npm run build && npm run export; cd ..

## Build only Meshery UI on your local machine.
build-meshery-ui:
	cd ui; npm run build && npm run export; cd ..

## Builds only the provider user interface on your local machine
build-provider-ui:
	cd provider-ui; npm run build && npm run export; cd ..

#-----------------------------------------------------------------------------
# Meshery Docs
#-----------------------------------------------------------------------------
#Incorporating Make docs commands from the Docs Makefile
jekyll=bundle exec jekyll

## Run Meshery Docs. Listen for changes.
site:
	cd docs; bundle install; $(jekyll) serve --drafts --livereload --config _config_dev.yml

## Build Meshery Docs on your local machine.
build-docs:
	cd docs; $(jekyll) build --drafts

## Run Meshery Docs in a Docker container. Listen for changes.
docker-docs:
	cd docs; docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:4.0.0 bash -c "bundle install; jekyll serve --drafts --livereload"

#-----------------------------------------------------------------------------
# Meshery Helm Charts
#-----------------------------------------------------------------------------
.PHONY: docs-charts
## Generate all Meshery Helm Chart documentation in markdown format.
docs-charts: docs-chart-operator docs-chart-meshery

## Generate Meshery Operator Helm Chart documentation in markdown format.
docs-chart-operator:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs 
	$(GOPATH)/bin/helm-docs -c install/kubernetes/helm/meshery-operator

## Generate Meshery Server and Adapters Helm Chart documentation in markdown format.
docs-chart-meshery:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs 
	$(GOPATH)/bin/helm-docs -c install/kubernetes/helm/meshery

.PHONY: lint-helm
## Lint all of Meshery's Helm Charts
lint-helm: lint-helm-operator lint-helm-meshery
## Lint Meshery Operator Helm Chart
lint-helm-operator:
	helm lint install/kubernetes/helm/meshery-operator --with-subcharts
## Lint Meshery Server and Adapter Helm Charts
lint-helm-meshery:
	helm lint install/kubernetes/helm/meshery --with-subcharts

#-----------------------------------------------------------------------------
# Meshery APIs
#-----------------------------------------------------------------------------
## Generate Meshery REST API specifications
gen-swagger:
	swagger generate spec -o ./helpers/swagger.yaml --scan-models

## Generate and servc Meshery REST API specifications
run-swagger:swagger-spec
	swagger serve ./helpers/swagger.yaml

## Build Meshery REST API documentation
docs-swagger:
	swagger generate spec -o ./docs/_data/swagger.yml --scan-models; \
	swagger flatten ./docs/_data/swagger.yml -o ./docs/_data/swagger.yml --with-expand --format=yaml

## Build Meshery GraphQL API documentation
docs-graphql:
	cd docs; build-docs; bundle exec rake graphql:compile_docs

## Build Meshery GraphQl API specifications
gen-gqlgen:
	cd internal/graphql; go run -mod=mod github.com/99designs/gqlgen generate

