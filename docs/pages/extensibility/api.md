---
layout: default
title: "Extensibility: APIs"
permalink: extensibility/api
type: Reference
abstract: 'Meshery architecture is extensible, offering an array of extension points and REST and GraphQL APIs.'
#redirect_from: extensibility
---
## Meshery's APIs

Each of Meshery's APIs are subject to the following authentication and authorization system.

### Authentication

Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers. Type of authentication is determined by the selected [Provider](#providers). Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication.

### Authorization

Currently, Meshery only requires a valid token in order to allow clients to invoke its APIs.

### Endpoints

Each of the API endpoints are exposed through [server.go](https://github.com/layer5io/meshery/blob/master/router/server.go). Endpoints are grouped by function (e.g. /api/mesh or /api/perf).

Alternatively, [Remote Providers](./providers) can extend Meshery's endpoints behind the `/api/extensions/` endpoint.

## GraphQL

Meshery provides a GraphQl API available through the default port of `/tcp`.

## REST

Meshery provides a REST API available through the default port of `9081/tcp`.
