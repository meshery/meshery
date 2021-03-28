ADAPTER_URLS := "mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004 mesherylocal.layer5.io:10006 mesherylocal.layer5.io:10008 mesherylocal.layer5.io:10009"

MESHERY_CLOUD_LOCAL="http://mesherylocal.layer5.io:9876"
MESHERY_CLOUD_DEV="http://localhost:9876"
MESHERY_CLOUD_PROD="https://meshery.layer5.io"
MESHERY_CLOUD_STAGING="https://staging-meshery.layer5.io"
GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_COMMITSHA=$(shell git rev-list -1 HEAD)
RELEASE_CHANNEL="edge"

# Build the CLI for Meshery - `mesheryctl`.
# Build Meshery inside of a multi-stage Docker container.
mesheryctl:
	cd mesheryctl; go build -o mesheryctl cmd/mesheryctl/main.go
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery .

# `make docker` builds Meshery inside of a multi-stage Docker container.
# This method does NOT require that you have Go, NPM, etc. installed locally.
docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery --build-arg TOKEN=$(GLOBAL_TOKEN) --build-arg GIT_COMMITSHA=$(GIT_COMMITSHA) --build-arg GIT_VERSION=$(GIT_VERSION) --build-arg RELEASE_CHANNEL=${RELEASE_CHANNEL} .

# Runs Meshery in a container locally and points to locally-running
#  Meshery Cloud for user authentication.
docker-run-local-cloud:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link meshery-cloud:meshery-cloud \
	-e PROVIDER_BASE_URLS=$(MESHERY_CLOUD_LOCAL) \
	-e DEBUG=true \
	-e ADAPTER_URLS=$(ADAPTER_URLS) \
	-p 9081:8080 \
	layer5/meshery ./meshery

# Runs Meshery in a container locally and points to remote
#  Meshery Cloud for user authentication.
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

# Runs Meshery on your local machine and points to locally-running
#  Meshery Cloud for user authentication.

run-local-cloud:
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
run-local:
	cd cmd; go clean; rm meshery; go mod tidy; \
	go build -ldflags="-w -s -X main.version=${GIT_VERSION} -X main.commitsha=${GIT_COMMITSHA} -X main.releasechannel=${RELEASE_CHANNEL}" -tags draft -a -o meshery; \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	./meshery; \
	cd ..

run-fast:
	cd cmd; go mod tidy; \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	go run main.go;

run-fast-cloud:
	cd cmd; go mod tidy; \
	PROVIDER_BASE_URLS=$(MESHERY_CLOUD_DEV) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	go run main.go;


golangci-run:
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

# setup wrk2 for local dev
# NOTE: setup-wrk does not work on Mac Catalina at the moment
setup-wrk2:
	cd cmd; git clone https://github.com/layer5io/wrk2.git; cd wrk2; make; cd ..

setup-nighthawk:
	cd cmd; git clone https://github.com/layer5io/nighthawk-go.git; cd wrk2; make setup; cd ..

#Incorporating Make docs commands from the Docs Makefile
jekyll=bundle exec jekyll

docs:
	$(jekyll) serve --drafts --livereload

build-docs:
	$(jekyll) build --drafts --livereload

docker-docs:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:3.8.5 bash -c "bundle install; jekyll serve --drafts --livereload"


.PHONY: chart-readme
chart-readme:
	go run github.com/norwoodj/helm-docs/cmd/helm-docs -c install/kubernetes/helm/