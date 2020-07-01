ADAPTER_URLS := "mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004 mesherylocal.layer5.io:10008"

MESHERY_CLOUD_LOCAL=http://mesherylocal.layer5.io:9876
MESHERY_CLOUD_DEV=http://localhost:9876
MESHERY_CLOUD_PROD=https://meshery.layer5.io
MESHERY_CLOUD_STAGING=https://staging-meshery.layer5.io

# Build the CLI for Meshery - `mesheryctl`.
# Build Meshery inside of a multi-stage Docker container.
mesheryctl:
	cd mesheryctl; go build -o mesheryctl cmd/mesheryctl/main.go
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery .

# `make docker` builds Meshery inside of a multi-stage Docker container.
# This method does NOT require that you have Go, NPM, etc. installed locally.
docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery --build-arg TOKEN=$(GLOBAL_TOKEN) .

# Runs Meshery in a container locally and points to locally-running 
#  Meshery Cloud for user authentication.
docker-run-local-cloud:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link meshery-cloud:meshery-cloud \
	-e SAAS_BASE_URL=$(MESHERY_CLOUD_LOCAL) \
	-e DEBUG=true \
	-e ADAPTER_URLS=$(ADAPTER_URLS) \
	-p 9081:8080 \
	layer5/meshery ./meshery

# Runs Meshery in a container locally and points to remote 
#  Meshery Cloud for user authentication.
docker-run-cloud:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	-e SAAS_BASE_URL=$(MESHERY_CLOUD_PROD) \
	-e DEBUG=true \
	-e ADAPTER_URLS=$(ADAPTER_URLS) \
	-v meshery-config:/home/appuser/.meshery/config \
  -v $(HOME)/.kube:/home/appuser/.kube:ro \
	-p 9081:8080 \
	layer5/meshery ./meshery

# Runs Meshery on your local machine and points to locally-running  
#  Meshery Cloud for user authentication.

run-local-cloud:
	cd cmd; go clean; rm meshery; go mod tidy; go build -tags draft -a -o meshery; \
	SAAS_BASE_URL=$(MESHERY_CLOUD_DEV) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	./meshery; \
	cd ..

# Builds and runs Meshery to run on your local machine.
#  and points to remote Meshery Cloud for user authentication.
run-local:
	cd cmd; go clean; rm meshery; go mod tidy; go build -tags draft -a -o meshery; \
	SAAS_BASE_URL=$(MESHERY_CLOUD_PROD) \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS=$(ADAPTER_URLS) \
	./meshery; \
	cd ..

run-tests:
	# -----Linting Check-----
	# GOPROXY=direct GOSUMDB=off go get github.com/mgechev/revive
	revive -config tools-config/revive-lint.toml -formatter friendly ./... \

	# ------Error Check------
	# GOPROXY=direct GOSUMDB=off GO111MODULE=on go get github.com/kisielk/errcheck
	errcheck -tags draft ./... \

	# ------Static Check------
	# GOPROXY=direct GOSUMDB=off GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck
	staticcheck -tags draft -checks all,-ST1003,-ST1000,-U1000 ./... \

	# -------Vet Check-------
	GOPROXY=direct GOSUMDB=off GO111MODULE=on go vet -tags draft ./... \

	# ----Security Check-----
	# go get github.com/securego/gosec/cmd/gosec
	gosec -exclude=G301,G304,G107 ./...
proto:
	# go get -u google.golang.org/grpc
	# go get -u github.com/golang/protobuf/protoc-gen-go
	# PATH=$(PATH):`pwd`/../protoc/bin:$(GOPATH)/bin
	# export PATH=$PATH:`pwd`/../protoc/bin:$GOPATH/bin
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

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

# Builds all user interfaces on your local machine.
build-ui:
	cd ui; npm run build && npm run export; cd ..
	cd provider-ui; npm run build && npm run export; cd ..

# setup wrk2 for local dev 
# NOTE: setup-wrk does not work on Mac Catalina at the moment
setup-wrk2:
	cd cmd; git clone git@github.com:layer5io/wrk2.git; cd wrk2; make; cd ..

#Incorporating Make docs commands from the Docs Makefile	
jekyll=bundle exec jekyll

docs:
	$(jekyll) serve --drafts --livereload

build-docs:
	$(jekyll) build --drafts --livereload

docker-docs:
	docker run --name meshery-docs --rm -p 4000:4000 -v `pwd`:"/srv/jekyll" jekyll/jekyll:3.8.5 bash -c "bundle install; jekyll serve --drafts --livereload"	


