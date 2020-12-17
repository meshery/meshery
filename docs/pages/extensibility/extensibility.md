---
layout: default
title: Extensibility
permalink: extensibility
type: Reference
abstract: 'Meshery architecture is extensible. Meshery provides several extension points for working with different service meshes via <a href="extensibility#adapters">adapters</a>, <a href="extensibility#load-generators">load generators</a> and <a href="extensibility#providers">providers</a>.'
redirect_from: reference/extensibility
---

Meshery has an extensible architecture with several extension points. Meshery provides several extension points for working with different service meshes via [adapters](#adapters), different [load generators](#load-generators) and different [providers](#providers). Meshery also offers a REST API.

**Guiding Principles for Extensibility**

The following principles are upheld in the design of Meshery's extensibility.

1. Recognize that different deployment environments have different systems to integrate with.
1. Offer a default experience that provides the optimal user experience.

## Extension Points

Meshery is not just an application. It is a set of microservices where the central component is itself called Meshery. The following points of extension are currently incorporated into Meshery.

**Types of Extension Points:**

1. [Providers]({{site.baseurl}}/extensibility/providers)
1. [Load Generators]({{site.baseurl}}/extensibility/load-generators)
1. [Adapters]({{site.baseurl}}/extensibility/adapters)
1. [REST API](#rest-api)

## REST API
Meshery provides a REST API available through the default port of 9081/tcp.

### Authentication
Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers.
Type of authentication is determined by the selected [Provider](#providers).

### Authorization
Currently, Meshery only requires a valid token in order to allow clients to invoke its APIs.

### Endpoints
Each of the API endpoints are exposed through [server.go](https://github.com/layer5io/meshery/blob/master/router/server.go).
Endpoints are grouped by function (e.g. /api/mesh or /api/perf).
