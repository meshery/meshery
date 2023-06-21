---
layout: default
title: "Extensibility: Service Mesh Adapters"
permalink: extensibility/adapters
type: Extensibility
abstract: 'Meshery architecture is extensible. Meshery provides several extension points for working with different service meshes via <a href="extensibility#adapters">adapters</a>, <a href="extensibility#load-generators">load generators</a> and <a href="extensibility#providers">providers</a>.'
language: en
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

### [Meshery Adapter Library](https://github.com/meshery/meshery-adapter-library)

This section contains a high level overview of the meshery-adapter-library, its purpose and architecture. For details, the reader is referred to the documentation and the code in the repository.

The main purpose of the meshery-adapter-library is to:
- provide a set of interfaces, some with default implementations, to be used and extended by adapters.
- implement common cross cutting concerns like logging, errors, and tracing
- provide a mini framework implementing the gRPC server that allows plugging in the mesh specific configuration and - operations implemented in the adapters.
- provide middleware extension points

[![Meshery Adapter Library]({{ site.baseurl }}/assets/img/adapters/meshery-adapter-library.svg
)]({{ site.baseurl }}/assets/img/adapters/meshery-adapter-library.svg)

#### Overview and Usage

The library consists of interfaces and default implementations for the main and common functionality of an adapter. It also provides a mini-framework that runs the gRPC adapter service, calling the functions of handlers injected by the adapter code. This is represented in an UML-ish style in the figure below. The library is used in all of Meshery's adapters.
