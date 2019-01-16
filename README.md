# Meshery

A service mesh playground to faciliate learning about functionality and performance of different service meshes. Meshery incorporates the collection and display of metrics from applications running in the playground.

--> See the [performance benchmark design document](https://docs.google.com/document/d/1nV8TunLmVC8j5cBELT42YfEXYmhG3ZqFtHxeG3-w9t0/edit?usp=sharing)

Meshery is written in `Go` (Golang) and leverages Go Modules. The `deployment_yaml` folder contains the configuration yaml to deploy Meshery on Kubernetes, which includes a Deployment, Service, Service Entries and Virtual Services configurations.

![Istio Playground](/static/img/istio-playground.png?raw=true "Istio Playground")

## Prequisites
1. Docker engine (e.g. Docker for Desktop).
1. Kubernetes cluster (preferably version 1.10+).

## Istio Playground App
A sample Istio app is included in Meshery. 

### Running Meshery
#### Prerequisites
1. Istio version 1.0.3+ in `istio-system` namespace along with the Istio ingress gateway.
1. Istio Solarwinds Mixer adapter is configured with a valid AppOptics token.
1. The canonical Istio _bookinfo_ sample application deployed in the `default` namespace.

#### Run
To run the service mesh playground either:
1. Deploy Meshery (`kubectl apply -f deployment_yamls/deployment.yaml`).

## Linkerd Playground App
_coming soon_
### Running Meshery
#### Prerequisites
#### Run

### Building Meshery
A sample Makefile is included to build and package the app as a Docker image.
1. `Docker` to build the image.
1. `Go` version 1.11+ installed if you want to make changes to the existing code.
1. Clone this repository (`git clone https://github.com/layer5io/meshery.git`).
1. Build the Meshery Docker image (`docker build . -t meshery`).
1.1. _pre-built images available: https://hub.docker.com/u/layer5/_

# About - Layer 5 - Service Meshes

[Layer5.io](https://layer5.io) is a repository for information pertaining to the technology ecosystem surrounding service meshes, api gateways, edge proxies, ingress and egress controllers - - microservice management in cloud native environments.

## Contributing

Contributions, updates, [discrepancy reports](/../../issues) and [pull requests](/../../pulls) are welcome. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
