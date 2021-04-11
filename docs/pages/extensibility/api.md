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

Meshery provides its GraphQl API at `localhost:9081/api/system/graphql/query`. A GraphQL request can be made as a POST request to the endpoint with the query as the payload.

Explore the Meshery GraphQL API using the `interactive Playground` provided with meshery instance at `localhost:9081/api/system/graphql/playground`.

Meshery GrahphQL API can be used to peform three operations:

- Queries for data retrival.
- Mutations for creating, updating, and deleting data.
- Subscriptions for watching data changes.

### Queries

{% for data in site.data.GraphQL.Queries %}
#### `{{data.name}}`

{{data.description}}

{% if data.arguments %}
**Arguments**

<table>
<thead>
    <tr>
        <th>Name</th>
        <th>Type</th>
        <th>Description</th>
    </tr>
</thead>
    {% for arg in data.arguments %}
    <tr>
        <td><code>{{arg.name}}</code></td>
        <td><code>{{arg.type}}</code></td>
        <td>{{arg.desc}}</td>
    </tr>
    {% endfor %}
</table>
{% else %}
_No Arguments needed._
{% endif %}
{% endfor %}

### Mutations

{% for data in site.data.GraphQL.Mutations %}
#### `{{data.name}}`

{{data.description}}

{% if data.arguments %}
**Arguments**

<table>
<thead>
    <tr>
        <th>Name</th>
        <th>Type</th>
        <th>Description</th>
    </tr>
</thead>
    {% for arg in data.arguments %}
    <tr>
        <td><code>{{arg.name}}</code></td>
        <td><code>{{arg.type}}</code></td>
        <td>{{arg.desc}}</td>
    </tr>
    {% endfor %}
</table>
{% else %}
_No Arguments needed._
{% endif %}
{% endfor %}

### Subscriptions

{% for data in site.data.GraphQL.Subscriptions %}
#### `{{data.name}}`

{{data.description}}

{% if data.arguments %}
**Arguments**

<table>
<thead>
    <tr>
        <th>Name</th>
        <th>Type</th>
        <th>Description</th>
    </tr>
</thead>
    {% for arg in data.arguments %}
    <tr>
        <td><code>{{arg.name}}</code></td>
        <td><code>{{arg.type}}</code></td>
        <td>{{arg.desc}}</td>
    </tr>
    {% endfor %}
</table>
{% else %}
_No Arguments needed._
{% endif %}
{% endfor %}

### Object Types

{% for data in site.data.GraphQL.ObjectTypes %}
#### `{{data.name}}`

{{data.description}}

{% if data.arguments %}
**Arguments**

<table>
<thead>
    <tr>
        <th>Field</th>
        <th>Type</th>
    </tr>
</thead>
    {% for arg in data.arguments %}
    <tr>
        <td><code>{{arg.name}}</code></td>
        <td><code>{{arg.type}}</code></td>
    </tr>
    {% endfor %}
</table>
{% endif %}
{% endfor %}

### Enumeration Types

{% for data in site.data.GraphQL.EnumerationTypes %}
#### `{{data.name}}`

{{data.description}}

{% if data.arguments %}
**Arguments**

<table>
<thead>
    <tr>
        <th>Value</th>
        <th>Description</th>
    </tr>
</thead>
    {% for arg in data.arguments %}
    <tr>
        <td><code>{{arg.value}}</code></td>
        <td>{{arg.desc}}</td>
    </tr>
    {% endfor %}
</table>
{% endif %}
{% endfor %}

## REST

Meshery provides a REST API available through the default port of `9081/tcp`.
