---
layout: default
title: "Extensibility: APIs"
permalink: extensibility/api
type: Extensibility
abstract: "Meshery architecture is extensible, offering an array of extension points and REST and GraphQL APIs."
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

### How to get your token

There are two ways to get your token:
1. Meshery UI
2. Meshery CLI

- What are tokens?
Tokens allow a user to authenticate with cloud apps and bypass two step verification and SSO, retrieve data from the instance through REST APIS
An authentication token (security token) is a “trusted device” used to access an electronically restricted resource (usually an application or a corporate network). It can be seen as an electronic key that enables a user to authenticate and prove his identity by storing some sort of personal information.


- Firstly, lets discuss how to get the token using Meshery-UI
1. Login/Signup to the dashboard of Meshery-UI
2. Click on Get token from the dropdown of profile section, 
3. You can download your token from there

- Secondly, you can download your token using Meshery-CLI

- mesheryctl system config --token [path-to-file]

Using this command, you can get your tokens easily
  

#### GraphQL

Meshery provides its GraphQl API at `hostname:9081/api/graphql/query`. A GraphQL request can be made as a POST request to the endpoint with the query as the payload.

Explore the Meshery GraphQL API using the `interactive Playground` provided with Meshery instance at `localhost:9081/api/system/graphql/playground`.

Meshery GrahphQL API can be used to perform three operations:

- Queries for data retrieval.
- Mutations for creating, updating, and deleting data.
- Subscriptions for watching for any data changes.

{% include alert.html type="dark" title="Meshery's GraphQL Schema" content="See <a href='/reference/graphql-apis'>GraphQL API Reference</a> Self-generated API documentation for Meshery’s GraphQL API. The API can be explored interactively using the GraphQL Playground. Documentation is generated from Meshery’s GraphQL schema. Each table below documents a GraphQL type." %}

#### REST

Meshery provides a REST API available through the default port of `9081/tcp`.

{% include alert.html type="dark" title="Meshery's REST API Reference" content="See <a href='/reference/rest-apis'>REST API Reference</a> Self-generated API documentation for Meshery’s REST API. Documentation is generated from Meshery’s Open API schema." %}

## See Also

- [Extension Points]({{site.baseurl}}/extensibility}})