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

all: error run

#-----------------------------------------------------------------------------
# Meshery Server, CLI and API Errors
#-----------------------------------------------------------------------------
.PHONY: error
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers --skip-dirs mesheryctl

#-----------------------------------------------------------------------------
# Meshery CLI Native Build
#-----------------------------------------------------------------------------
# Build the CLI for Meshery - `mesheryctl`.
mesheryctl:
	cd mesheryctl; go build -o mesheryctl cmd/mesheryctl/main.go

#-----------------------------------------------------------------------------
# Meshery Server and UI Containerized Builds
#-----------------------------------------------------------------------------
# `make docker` builds Meshery inside of a multi-stage Docker container.
# This method does NOT require that you have Go, NPM, etc. installed locally.
docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery --build-arg TOKEN=$(GLOBAL_TOKEN) --build-arg GIT_COMMITSHA=$(GIT_COMMITSHA) --build-arg GIT_VERSION=$(GIT_VERSION) --build-arg RELEASE_CHANNEL=${RELEASE_CHANNEL} .

# Runs Meshery in a container locally and points to locally-running
# Meshery Cloud for user authentication.
docker-run-local-cloud:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link meshery-cloud:meshery-cloud \
	-e PROVIDER_BASE_URLS=$(REMOTE_PROVIDER_LOCAL) \
	-e DEBUG=true \
	-e ADAPTER_URLS=$(ADAPTER_URLS) \
	-p 9081:8080 \
	layer5/meshery ./meshery

# Runs Meshery in a container locally and points to remote
# Remote Provider for user authentication.
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
# setup wrk2 for local dev
# NOTE: setup-wrk does not work on Mac Catalina at the moment
setup-wrk2:
	cd cmd; git clone https://github.com/layer5io/wrk2.git; cd wrk2; make; cd ..

setup-nighthawk:
	cd cmd; git clone https://github.com/layer5io/nighthawk-go.git; cd nighthawk-go; make setup; cd ..

# Runs Meshery on your local machine and points to locally-running
#  Meshery Cloud for user authentication.
run-local-cloud: error
	cd cmd; go clean; rm meshery; go mod tidy; \
	go build -ldflags="-w -s -X main.version=${GIT_VERSION} -X main.commitsha=${GIT_COMMITSHA} -X main.releasechannel=${RELEASE_CHANNEL}" -tags draft -a -o meshery; \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_DEV) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	./meshery; \
	cd ..

# Builds and runs Meshery to run on your local machine.
#  and points to remote Meshery Cloud for user authentication.
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

run-fast:
	cd cmd; go mod tidy; \
	BUILD="$(GIT_VERSION)" \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	APP_PATH=$(APPLICATIONCONFIGPATH) \
	go run main.go;

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

golangci-run: error
	GO111MODULE=off GOPROXY=direct GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0;
	$(GOPATH)/bin/golangci-lint run

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
# Installs dependencies for building the user interface.
setup-ui-libs:
	cd ui; npm i; cd ..
	cd provider-ui; npm i; cd ..

# Runs the Meshery UI interface on your local machine.
run-ui-dev:
	cd ui; npm run dev; cd ..

# Runs the Provider UI interface on your local machine.
run-provider-ui-dev:
	cd provider-ui; npm run dev; cd ..

# Runs the lint on Meshery UI interface on your local machine.
run-ui-lint:
	cd ui; npm run lint; cd ..

# Runs the lint on Meshery UI interface on your local machine.
run-provider-ui-lint:
	cd provider-ui; npm run lint; cd ..

# Runs the test on Provider UI interface on your local machine.
run-provider-ui-test:
	cd provider-ui; npm run test; cd ..

# Builds all user interfaces on your local machine.
build-ui:
	cd ui; npm run build && npm run export; cd ..
	cd provider-ui; npm run build && npm run export; cd ..

# Builds only the meshery user interface on your local machine
build-meshery-ui:
	cd ui; npm run build && npm run export; cd ..

# Builds only the provider user interface on your local machine
build-provider-ui:
	cd provider-ui; npm run build && npm run export; cd ..

#-----------------------------------------------------------------------------
# Meshery Docs
#-----------------------------------------------------------------------------
#Incorporating Make docs commands from the Docs Makefile
jekyll=bundle exec jekyll

site:
	cd docs; bundle install; $(jekyll) serve --drafts --livereload --config _config_dev.yml

build-docs:
	cd docs; $(jekyll) build --drafts

docker-docs:
	cd docs; docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:4.0.0 bash -c "bundle install; jekyll serve --drafts --livereload"

#-----------------------------------------------------------------------------
# Meshery Helm Charts
#-----------------------------------------------------------------------------
.PHONY: chart-readme
chart-readme: helm-lint

.PHONY: chart-readme
chart-readme: chart-readme-operator
chart-readme: chart-readme-meshery
chart-readme-operator:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs 
	$(GOPATH)/bin/helm-docs -c install/kubernetes/helm/meshery-operator
chart-readme-meshery:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs 
	$(GOPATH)/bin/helm-docs -c install/kubernetes/helm/meshery

.PHONY: helm-lint
helm-lint: helm-lint-operator
helm-lint: helm-lint-meshery
helm-lint-operator:
	helm lint install/kubernetes/helm/meshery-operator --with-subcharts
helm-lint-meshery:
	helm lint install/kubernetes/helm/meshery --with-subcharts

#-----------------------------------------------------------------------------
# Meshery APIs
#-----------------------------------------------------------------------------
swagger-spec:
	swagger generate spec -o ./helpers/swagger.yaml --scan-models

swagger-run:swagger-spec
	swagger serve ./helpers/swagger.yaml

swagger-docs:
	swagger generate spec -o ./docs/_data/swagger.yml --scan-models; \
	swagger flatten ./docs/_data/swagger.yml -o ./docs/_data/swagger.yml --with-expand --format=yaml

graphql-docs:
	cd docs; build-docs; bundle exec rake graphql:compile_docs

gqlgen-generate:
	cd internal/graphql; go run -mod=mod github.com/99designs/gqlgen generate

