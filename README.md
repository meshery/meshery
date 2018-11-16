# Meshery

A service mesh playground.


## Istio Playground App
__NOTE__: this is considered EXPERIMENTAL and is not yet recommended for production systems.

This repo contains a sample Istio playground app. It is written in `Go` (Golang) and leverages Go Modules.

To try out this app you will need:
- `Docker` on your local to build the image.
- (Optional) `Go` version 1.11+ installed if you want to make changes to the existing code.
- You have a functional Kubernetes (preferably version 1.10+) cluster.
- Istio is deployed in the `istio-system` namespace along with the istio ingress gateway.
- The canonical istio book info app is deployed in the `default` namespace.


A sample Makefile is included to build and package the app as a Docker image.

Deployment yaml folder contains the needed yaml to deploy this app on kubernetes with the needed deployment, service, service entries and virtual services configurations.
