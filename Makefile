fortio_port := 9080

fortio:
	docker run --name fortio -p $(fortio_port):8080 -p 8079:8079 -d fortio/fortio server

docker:
	docker build -t layer5/meshery .

docker-run-local-saas:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link fortio:fortio \
	--link meshery-saas:meshery-saas \
	-e SAAS_BASE_URL="http://meshery-saas:9876" \
	-e EVENT=istioPlay01 \
	-e FORTIO_URL="http://fortio:8080/fortio/" \
	-e DEBUG=true \
	-p 9081:8080 \
	layer5/meshery ./meshery

docker-run-saas:
	(docker rm -f meshery) || true
	docker run --name meshery -d \
	--link fortio:fortio \
	-e SAAS_BASE_URL="https://meshery.layer5.io" \
	-e EVENT=istioPlay01 \
	-e FORTIO_URL="http://fortio:8080/fortio/" \
	-e DEBUG=true \
	-p 9081:8080 \
	layer5/meshery ./meshery

run-local-saas:
	cd cmd; go clean; go build -a -o meshery; \
	SAAS_BASE_URL="http://meshery-saas:9876" \
	EVENT=istioPlay01 \
	FORTIO_URL="http://localhost:9080/fortio/" \
	PORT=9081 \
	DEBUG=true \
	./meshery; \
	cd ..

run-saas:
	cd cmd; go clean; go build -a -o meshery; \
	SAAS_BASE_URL="https://meshery.layer5.io" \
	EVENT=istioPlay01 \
	FORTIO_URL="http://localhost:9080/fortio/" \
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