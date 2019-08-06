make:
	cd mesheryctl; go build -o mesheryctl
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery .

docker:
	DOCKER_BUILDKIT=1 docker build -t layer5/meshery .

docker-run-local-saas:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link meshery-saas:meshery-saas \
	-e SAAS_BASE_URL="http://mesherylocal.layer5.io:9876" \
	-e DEBUG=true \
	-e ADAPTER_URLS="mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001" \
	-p 9081:8080 \
	layer5/meshery ./meshery

docker-run-saas:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	-e SAAS_BASE_URL="https://meshery.layer5.io" \
	-e DEBUG=true \
	-e ADAPTER_URLS="mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001" \
	-p 9081:8080 \
	layer5/meshery ./meshery

run-local-saas:
	cd cmd; go clean; go build -tags draft -a -o meshery; \
	SAAS_BASE_URL="http://mesherylocal.layer5.io:9876" \
	PORT=9081 \
	DEBUG=true \
	ADAPTER_URLS="mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001" \
	./meshery; \
	cd ..

run-local:
	cd cmd; go clean; go build -tags draft -a -o meshery; \
	SAAS_BASE_URL="https://meshery.layer5.io" \
	PORT=9081 \
	DEBUG=true \
	./meshery; \
	cd ..


proto:
	# go get -u google.golang.org/grpc
	# go get -u github.com/golang/protobuf/protoc-gen-go
	# PATH=$(PATH):`pwd`/../protoc/bin:$(GOPATH)/bin
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

setup-ui-libs:
	cd ui; npm i; cd ..

run-ui-dev:
	cd ui; npm run dev; cd ..

build-ui:
	cd ui; npm run build && npm run export; cd ..
