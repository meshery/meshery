---
layout: enhanced
title: "Extensibility: APIs"
permalink: extensibility/api
type: Extensibility
abstract: "Meshery architecture is extensible, offering an array of extension points and REST and GraphQL APIs."
language: en
#redirect_from: extensibility
---

## Meshery's APIs

Each of Meshery's APIs are subject to the following authentication and authorization system. The REST and GraphQL APIs are available on the same network port you use to connect to Meshery UI, which is `9081/tcp` by default. Each of the API endpoints are exposed through [server.go](https://github.com/meshery/meshery/blob/master/server/router/server.go). Endpoints are grouped by function (e.g. `/api/system/version`). Additionally, [Remote Providers](./providers) can extend Meshery's endpoints behind the `/api/extensions/` endpoint.

### API Reference

<main>
  
<input id="tab1" type="radio" name="tabs" checked>
<label for="tab1">REST API</label>

<input id="tab2" type="radio" name="tabs">
<label for="tab2">GraphQL API</label>

<section class="tabbed" id="content1">
<p>Meshery Server's REST API is available at <code>[hostname]:[port]/api/</code>.</p>

{% include alert.html type="dark" title="Meshery's REST API Reference" content="See <a href='/reference/rest-apis'>REST API Reference</a> Self-generated API documentation for Meshery’s REST API. Documentation is generated from Meshery’s Open API schema." %}

</section>

<section class="tabbed" id="content2">


<p>Meshery Server's GraphQl API is available at <code>{hostname]:[port]/api/graphql/query</code>. A GraphQL request can be made as a POST request to the endpoint with the query as the payload. Meshery Server's GraphQL API can be used to perform three operations:</p>

<ul>
  <li>Queries for data retrieval.</li>
  <li>Mutations for creating, updating, and deleting data.</li>
   <li>Subscriptions for watching for any data changes.</li>
</ul>
<p>
Explore the Meshery GraphQL API using the interactive Playground provided with Meshery instance at http://localhost:9081/api/system/graphql/playground.</p>

{% include alert.html type="dark" title="Meshery's GraphQL API Reference" content="See <a href='/reference/graphql-apis'>GraphQL API Reference</a> Self-generated API documentation for Meshery’s GraphQL API. The API can be explored interactively using the GraphQL Playground. Documentation is generated from Meshery’s GraphQL schema. Each table below documents a GraphQL type." %}

  </section>
    
</main>

### Authentication

Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers. Type of authentication is determined by the selected [Provider](#providers). Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication.

{% include alert.html type="dark" title="What are authentication tokens?" content="Meshery authentication tokens allow users or systems to authenticate with Meshery Server via either its two clients, <a href='/reference/mesheryctl'>Meshery >CLI</a> and <a href='/extensibility/api#how-to-get-your-token'>UI</a>, or its two APIs: <a href='/reference/rest-apis'>REST</a> or <a href='/reference/graphql-apis'>GraphQL</a>. <p>Meshery's authentication token system provide secure access to Meshery's powerful management features.</p>" %}

### How to get your token

There are two ways to get your authentication token:

1. Meshery UI
2. Meshery CLI

Using Meshery UI, you can get a copy of your authentication token by following these steps:

1. Log into Meshery by selecting your identity provider of choice (typically found at `http:<meshery-server>:9081/provider`)
2. Navigate to your user's avatar in the upper lefthand corner and select "Get Token" from the dropdown of profile section.

Using Meshery CLI, you can get a copy of your authentication token by executing this command:

{% include code.html code="mesheryctl system login" %}

In order to use this command, you must have a web browser available on your system (this command cannot be executed on a headless system).

{% include alert.html type="dark" title="How to use the token for requests to Meshery Rest API using API clients" content=" Download your token from Meshery UI and copy the token value from the downloaded file, this will be used for authentication." %}
{% include code.html code="meshery-provider: Meshery
token: your token"%}
{% include alert.html type="dark" title="" content="Navigate to the cookies section your API testing platform, to set the appropriate name and value for the cookie, then make sure to save the configured cookies to initiate the request. <p> The request will be sent to the Meshery API using the configured cookies for authentication</p>"%}

### Endpoints

Each of the API endpoints are exposed through [server.go](https://github.com/layer5io/meshery/blob/master/router/server.go). Endpoints are grouped by function (e.g. /api/mesh or /api/perf).

Alternatively, [Remote Providers](./providers) can extend Meshery's endpoints behind the `/api/extensions/` endpoint.

## Authorization

While Meshery only requires a valid token in order to allow clients to invoke its APIs, Remote Providers can optionally enforce key-based permissions.

{% if page.suggested-reading != false and page.title and page.type and page.category and page.url %}
{% include_cached suggested-reading.html  title=page.title type=page.type category=page.category url=page.url language="en" %}
{% endif %}
