fortio:
	docker run -p $(fortio_port):8080 -p 8079:8079 -d fortio/fortio server

docker:
	docker build -t layer5/meshery-istio .
