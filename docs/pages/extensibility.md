---
layout: page
title: Extensibility
permalink: extensibility
---

# Preface
Meshery provides several extension points for working with different service meshes, load generators and providers. Meshery also offers a REST API.

## Guiding Principles for Extensibility

1. Recognize that different deployment environments have different systems to integrate with.
1. Offer a default experience that provides the optimal user experience.

# Extension Points
Meshery is not just an application. It is a set of microservices where the central component is itself called Meshery. The following points of extension are currently incorporated into Meshery:

1. Providers
1. Load Generators
1. Adapters
1. REST API

## Providers
Meshery interfaces with Providers through a Go interface. The Provider implementations have to be placed in the code and compiled together today. A Provider instance will have to be injected into Meshery when the program starts.

Eventually, we will be looking to keep the implementation of Providers separate so that they are brought in through a separate process and injected into Meshery at runtime (OR) change the way the code works to make the Providers invoke Meshery.

Providers as an object have the following attributes:

```
{
 id: aower2-1234xd-1234 [guid] 	# future consideration
 type: local, [ local | remote ]
 display-name: "None", [ string ]
 description: "Default Provider
   - feature 1 
   - feature 2", [ multi-line string ] 
 capabilities: [ 			# future consideration
    {featureName:"perf_results", present: true}, 
    {featureName:"two_factor_auth", present: true}
   ]
}
```

### What functionality do Providers perform? 
- Authentication and Authorization
 - Examples: session management, two factor authentication, LDAP integration
- Long-Term Persistence
 - Storage and retrieval of performance test results
 - Storage and retrieval of user preferences

### Types of providers
Two types of providers are defined in Meshery: local and remote.

**Remote Providers**
Name: “Meshery” (default)
- Authentication and Authorization
- Results long term persistence
- Save environment setup
- Retrieve performance test results
- Free to use

**Local Provider**
Name: “None”
- No Authentication
- Local storage of results. Mainly ephemeral.
- Environment setup not saved.
- No performance test result history.
- Free to use.


Meshery provides the ability for you as a service mesh manager to customize your service mesh deployment.

## Load Generators
So, that the load generator can be chosen at runtime based on user preference. Users may prefer to use one load generator over the next given the difference of capabilities between load generators.

### What function do load generators in Meshery provide? 
Load generators will provide the capability to run load tests from Meshery. As of today the load generators are embedded as libraries in Meshery and Meshery invokes the load generators APIs with the right load test options to run the load test. At the moment, Meshery has support for HTTP load generators. Support for GRPC and TCP load testing is on the roadmap. Meshery has functional integration with fortio and wrk2. 

### Why support multiple load generators?
Different use cases and different opinions call for different approaches to statistical analysis of the performance results. For example, wrk2 accounts for a concept called Coordinated Omission.

### Which are currently supported?
1. [fortio](https://github.com/fortio/fortio) - Fortio load testing library, command line tool, advanced echo server and web UI in go (golang). Allows to specify a set query-per-second load and record latency histograms and other useful stats.
1. [wrk2](https://github.com/giltene/wrk2) - A constant throughput, correct latency recording variant of wrk.


## Service Mesh Adapters
### What are Meshery adapters?
Adapters allow Meshery to interface with the different service meshes. See a list of all available [service mesh adapters](service-meshes/adapters).

### Adapter Capabilities
Meshery communicates with adapters over grpc. Adapters establish communication with the service mesh. Adapters have a predefined set of operations which are grouped based on predefined operation types. 

The predefined operation types are:
- Install
- Sample application
- Config
- Validate
- Custom

### How to create new adapter?

Meshery uses adapters to provision and interact with different service meshes. Follow these instructions to create a new adapter or modify and existing adapter.

See the [CONTRIBUTING.md](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#adapter) for additional information and specific steps.

Tip: The [Meshery adapter for Istio](https://github.com/layer5io/meshery-istio) is a good reference adapter to use as an example of a Meshery adapter written in Go.

## REST API
Meshery provides a REST API availble through the default port of 9081/tcp.

### Authentication
Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers.
Type of authentication is determined by the selected [Provider](#providers).

### Authorization
Currently, Meshery only requires a valid token in order to allow clients to invoke its APIs.

### Endpoints
Each of the API endpoints are exposed through [server.go]](https://github.com/layer5io/meshery/blob/master/router/server.go).
Endpoints are grouped by function (e.g. /api/mesh or /api/perf).
