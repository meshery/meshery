fortio_port := 9080

fortio:
	docker run --name fortio -p $(fortio_port):8080 -p 8079:8079 -d fortio/fortio server

docker:
	docker build -t layer5/meshery .

docker-run:
	docker run -d \
	-e TWITTER_APP_HOST="http://auth.layer5.ga" \
	-e EVENT=istioPlay01 \
	-e FORTIO_URL="http://localhost:$(fortio_port)/fortio/" \
	-e PRODUCT_PAGE_URL="http://localhost:8080/" \
	layer5/meshery