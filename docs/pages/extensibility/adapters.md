---
layout: default
title: "Extensibility: Service Mesh Adapters"
permalink: extensibility/adapters
type: Extensibility
abstract: 'Meshery architecture is extensible. Meshery provides several extension points for working with different service meshes via <a href="extensibility#adapters">adapters</a>, <a href="extensibility#load-generators">load generators</a> and <a href="extensibility#providers">providers</a>.'
#redirect_from: extensibility
---
## Guiding Principles of Adapter Design

Adapters allow Meshery to interface with the different service meshes. Review the list of all available [service mesh adapters](../concepts/architecture/service-meshes/adapters.md). See the [Meshery Architecture](../concepts/architecture/architecture.md) diagrams for visuals on how adapters relate to other Meshery components.

Meshery upholds the following guiding principles for adapter design:

1. **Adapters allow Meshery to interface with the different service meshes, exposing their differentiated value to users.**
  - Service mesh projects should be encouraged to maintain their own adapters. Allowing them to expose their differentiated capabilities encourages this.
1. **Adapters should avoid wheel reinvention, but seek to leverage the functionality provided by service meshes under management.**
  - This both reduces sustaining costs and improves reliability.


### Adapter Capabilities

Meshery communicates with adapters over grpc. Adapters establish communication with the service mesh. Adapters have a predefined set of operations which are grouped based on predefined operation types. 

The predefined operation types are:

- Install
- Sample application
- Config
- Validate
- Custom

## Meshery Adapter Codebase Overview

[Common libraries](https://docs.google.com/presentation/d/1uQU7e_evJ8IMIzlLoBi3jQSRvpKsl_-K1COVGjJVs30/edit#) are used to avoid code duplication and apply DRY.

### [MeshKit](https://github.com/layer5io/meshkit)

The code hierarchy is pluggable and independent from one another. There can be N number of packages depending upon the use case.
- `errors/` - holds the implementations and the error handlers and error codes which are used across projects.
- `logger/` - holds the implementations of logging handler and custom attributes to add if any.
- `utils/` - holds all the utility functions that are specific to meshery projects and are to be used generically across all of them.
- `tracing/` - holds the implementations of tracing handlers with different tracing providers like jaeger,newrelic, etc.

Each package inside a meshkit is a handler interface implementation, the implementation could be from any third-party packages or the go-kit.

### [Meshery Adapter Libary](https://github.com/meshery/meshery-adapter-library)

This section contains a high level overview of the meshery-adapter-library, its purpose and architecture. For details, the reader is referred to the documentation and the code in the repository.

The main purpose of the meshery-adapter-library is to:
- provide a set of interfaces, some with default implementations, to be used and extended by adapters.
- implement common cross cutting concerns like logging, errors, and tracing
- provide a mini framework implementing the gRPC server that allows plugging in the mesh specific configuration and - operations implemented in the adapters.
- provide middleware extension points

#### Overview and Usage

The library consists of interfaces and default implementations for the main and common functionality of an adapter. It also provides a mini-framework that runs the gRPC adapter service, calling the functions of handlers injected by the adapter code. This is represented in an UML-ish style in the figure below. The library is used in all of Meshery's adapters.

## Contributing to Meshery Adapters
With the [CONTRIBUTING.md](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#adapter) in mind, understand that development follows the usual fork-and-pull request workflow described here, see also GitHub Process. On forking GitHub deactivates all workflows. It is safe and good practice to activate them such that the code is validated on each push. This requires that branches filter for “on push” is set to ‘**’ to be triggered also on branches containing ‘/’  in their name.  The actions are parameterized using secrets (see Build & Release Strategy). The Docker image is only built and pushed to Docker Hub if a tag is pushed and the corresponding authentication information is configured. The only secret that should be set in each fork is GO_VERSION, specified in Build & Release Strategy, otherwise, the corresponding action’s default version is used.

Each commit has to be signed off, see [Contributing Overview]({{site.baseurl}}/project/contributing).

### Running an adapter as a container

Testing your local changes running as a container can be accomplished in two ways:

1. Define the adapter’s address in the UI: Unless the running container is named as specified in the docker-run target in the Makefile, the container has to be removed manually first. Then, run `make docker` followed by `make docker-run`. Then, connect to the adapter in the UI in “Settings>Service Meshes” using `localhost:<port>` if the meshery server is running as a binary, or <docker IP address>:<port> if it is running as a docker container.
1. Using mesheryctl: In `~/.meshery/meshery.yaml`, change the tag specifying the image of the adapter to “latest”. Run make docker, followed by `mesheryctl system start --skip-update`. This assumes mesheryctl system start has been executed at least once before.

### Running an adapter as a process

Another way to test your local changes is to run the adapter as a process. To do this, clone the `meshery/meshery` repository, and start Meshery Server using `make server`. Start the adapter from your IDE, or by executing `make run`. Then, in Meshery UI, add the adapter using “localhost:<PORT>”.

### Creating a new Meshery Adapter

Meshery uses adapters to manage and interact with different service meshes. Meshery adapters are written in Go. Whether you are creating a new adapter or modifying an existing adapter, be sure to read the [Meshery Adapters](https://docs.google.com/document/d/1b8JAMzr3Rntu7CudRaYv6r6ccACJONAB5t7ISCaPNuA/edit#) design specification. For new adapters, start with the Repository Template(https://github.com/layer5io/layer5-repo-template). 

1. Get the proto buf spec file from Meshery repo:
   `wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto`
1. Generate code
   1. Using Go as an example, do the following:
      - adding GOPATH to PATH: `export PATH=$PATH:$GOPATH/bin`
      - install grpc: `go get -u google.golang.org/grpc`
      - install protoc plugin for go: `go get -u github.com/golang/protobuf/protoc-gen-go`
      - Generate Go code: `protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/`
   1. For other languages, please refer to gRPC.io for language-specific guides.
1. Implement the service methods and expose the gRPC server on a port of your choice (e.g. 10000).

Tip: The [Meshery Adapter for Istio](https://github.com/layer5io/meshery-istio) is a good reference adapter to use as an example of a Meshery Adapter.
