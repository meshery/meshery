# Meshery

A service mesh playground to faciliate learning about functionality and performance of different service meshes. Meshery incorporates the collection and display of metrics from applications running in the playground.

- [Website](https://layer5.io/meshery)
- [Performance benchmark design document](https://docs.google.com/document/d/1nV8TunLmVC8j5cBELT42YfEXYmhG3ZqFtHxeG3-w9t0/edit?usp=sharing)
- [Architecture](https://docs.google.com/presentation/d/1UbuYMpn-e-mWVYwEASy4dzyZlrSgZX6MUfNtokraT9o/edit?usp=sharing)

![Service Mesh Playground](/public/static/img/meshery.png?raw=true "Service Mesh Playground")

## Functionality
1. Multi-mesh Performannce Benchmark
Meshery is intended to be a vendor and project-neutral utility for uniformly benchmarking the performance of service meshes. Between service mesh and proxy projects, a number of different tools *and results* exist. For example, Istio's [Performance and Scalability WG](https://github.com/istio/community/blob/master/WORKING-GROUPS.md#performance-and-scalability) currently uses a couple of different tools to measure Istio performance: [BluePerf](https://ibmcloud-perf.istio.io/regpatrol/) and [Fortio](https://fortio.istio.io).

1. Multi-mesh Functionalty Playground

### Running Meshery
## Service Mesh Playground
Sample applications will be included in Meshery. 

#### General Prerequisites
1. Docker engine (e.g. Docker for Desktop).
1. Kubernetes cluster (preferably version 1.10+).

#### Istio Playground Prerequisites
1. Istio version 1.0.3+ in `istio-system` namespace along with the Istio ingress gateway.
1. Istio Solarwinds Mixer adapter is configured with a valid AppOptics token.
1. The canonical Istio _bookinfo_ sample application deployed in the `default` namespace.

#### Run Istio Playground
To run the service mesh playground either:
- On Kubernetes
  - You can deploy Meshery to an existing kubernetes cluster using the provided yaml file into any namespace of your choice. For now let us deploy it to a namespace `meshery`: 

    ```
    kubectl create ns meshery
    kubectl -n meshery apply -f deployment_yamls/deployment.yaml
    ```
    If you want to use a different namespace, please change the name of the namespace in the `ClusterRoleBinding` section appropriately.
  - Meshery can be deployed either on/off the mesh.
  - If deployed on the same Kubernetes cluster as the mesh, you dont have to provide a kubeconfig file.
  - please review the yaml and make necessary changes as needed for your cluster
- On Docker
  - we will have to first deploy fortio: 
  
  ```
  docker run --name fortio -p $(fortio_port):8080 -p 8079:8079 -d fortio/fortio server
  ```
  - next to run Meshery we will need to place the kubeconfig file with name `config` in the current directory
  - next run Meshery <small>(please remember to update PRODUCT_PAGE_URL environment variable with a valid value)</small>:
  ```
  docker run --name meshery -d \
	--link fortio:fortio \
	-e BYPASS_AUTH=true \
	-e EVENT=istioPlay01 \
	-e FORTIO_URL="http://fortio:8080/fortio/" \
	-e PRODUCT_PAGE_URL="<URL TO THE BOOK INFO APP ON THE MESH>" \
	-p 9081:8080 \
	-v `pwd`/config:/kubeconfig \
	layer5/meshery ./meshery --kubeconfig /kubeconfig
    ```
  - If you prefer to use a Docker network you can remove the link from the above command
  - Now you should be able to access Meshery in your browser at `http://localhost:9081/play`

## Linkerd Playground App
_coming soon for Linkerd_
### Running Meshery
_coming soon for Linkerd_
#### Linkerd Prerequisites
_coming soon for Linkerd_
#### Run Linkerd Playground
_coming soon for Linkerd_

## Contributing
Contributions, updates, [discrepancy reports](/../../issues) and [pull requests](/../../pulls) are welcome. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

### Building Meshery
Meshery is written in `Go` (Golang) and leverages Go Modules. The `deployment_yaml` folder contains the configuration yaml to deploy Meshery on Kubernetes, which includes a Deployment, Service, Service Entries and Virtual Services configurations.

A sample Makefile is included to build and package the app as a Docker image.
1. `Docker` to build the image.
1. `Go` version 1.11+ installed if you want to make changes to the existing code.
1. Clone this repository (`git clone https://github.com/layer5io/meshery.git`).
1. Build the Meshery Docker image (`docker build -t layer5/meshery .`).
    1. _pre-built images available: https://hub.docker.com/u/layer5/_

# About Layer5
[Layer5.io](https://layer5.io) is a service mesh community, serving as a repository for information pertaining to the surrounding technology ecosystem (service meshes, api gateways, edge proxies, ingress and egress controllers) of microservice management in cloud native environments.

## License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
