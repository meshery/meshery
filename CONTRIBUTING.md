# <a name="contributing">Contributing</a>
Please do! Contributions, updates, [discrepancy reports](/../../issues) and [pull requests](/../../pulls) are welcome. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

See the [newcomers welcome guide](https://docs.google.com/document/d/14Fofs9BysojB5igihXBI_SsFWoSUu-QRsGnnFqUvR0M/edit) for how, where and why to contribute.

Not sure where to start? Grab an open issue with the [help-wanted label](../../labels/help%20wanted).
### License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).

# <a name="adapter">Writing an Adapter</a>
Meshery uses adapters to provision and interact with different service meshes. Follow these instructions to create a new adapter or modify and existing adapter.

1. Get the proto buf spec file from Meshery repo: 
```wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto```
1. Generate code
    1. Using Go as an example, do the following:
        - adding GOPATH to PATH: `export PATH=$PATH:$GOPATH/bin`
        - install grpc: `go get -u google.golang.org/grpc`
        - install protoc plugin for go: `go get -u github.com/golang/protobuf/protoc-gen-go`
        - Generate Go code: `protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/`
    1. For other languages, please refer to gRPC.io for language-specific guides.
1. Implement the service methods and expose the gRPC server on a port of your choice (e.g. 10000). 

_Tip:_ The [Meshery adapter for Istio](https://github.com/layer5io/meshery-istio) is a good reference adapter to use as an example of a Meshery adapter written in Go.

# <a name="building">Building Meshery</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. The `deployment_yaml` folder contains the configuration yaml to deploy Meshery on Kubernetes, which includes a Deployment, Service, Service Entries and Virtual Services configurations.

A sample Makefile is included to build and package the app as a Docker image.
1. `Docker` to build the image.
1. `Go` version 1.11+ installed if you want to make changes to the existing code.
1. Clone this repository (`git clone https://github.com/layer5io/meshery.git`).
1. Build the Meshery Docker image (`docker build -t layer5/meshery .`).
    1. _pre-built images available: https://hub.docker.com/u/layer5/_
