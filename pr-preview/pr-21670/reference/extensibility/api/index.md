# Extensibility: APIs

> Meshery architecture is extensible, offering an array of extension points and REST and GraphQL APIs.

Source: /pr-preview/pr-21670/reference/extensibility/api/

## Meshery's APIs

Each of Meshery's APIs are subject to the following authentication and authorization system. The REST and GraphQL APIs are available on the same network port you use to connect to Meshery UI, which is `9081/tcp` by default. Each of the API endpoints are exposed through [server.go](https://github.com/meshery/meshery/blob/master/server/router/server.go). Endpoints are grouped by function (e.g. `/api/system/version`). Additionally, [Remote Providers](/pr-preview/pr-21670/reference/extensibility/providers/) can extend Meshery's endpoints behind the `/api/extensions/` endpoint.

### API Reference

<div class="tab-container tab-level-1"><input id="api-reference-tabs-3-1" type="radio" name="api-reference-tabs-3" checked>
    <label for="api-reference-tabs-3-1"><i class="fa fa-cube" aria-hidden="true"></i>REST API
    </label>
    <section class="tabbed" id="api-reference-tabs-3-content-1"><p>Meshery Server's REST API is available at <code>http://[hostname]:[port]/api/</code>.</p>
<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">Meshery's REST API Reference</div>
<p>See <a href="/pr-preview/pr-21670/reference/references/rest-apis/">REST API Reference</a> Self-generated API documentation for Meshery&rsquo;s REST API. Documentation is generated from Meshery&rsquo;s Open API schema.</p>
</div></section><input id="api-reference-tabs-3-2" type="radio" name="api-reference-tabs-3">
    <label for="api-reference-tabs-3-2"><i class="fa fa-globe" aria-hidden="true"></i>GraphQL API
    </label>
    <section class="tabbed" id="api-reference-tabs-3-content-2"><p>Meshery Server's GraphQl API is available at <code>{hostname]:[port]/api/graphql/query</code>. A GraphQL request can be made as a POST request to the endpoint with the query as the payload. Meshery Server's GraphQL API can be used to perform three operations:</p>
<ul>
<li>Queries for data retrieval.</li>
<li>Mutations for creating, updating, and deleting data.</li>
<li>Subscriptions for watching for any data changes.</li>
</ul>
<p>Explore the Meshery GraphQL API using the interactive Playground provided with Meshery instance at http://localhost:9081/api/system/graphql/playground.</p>
<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">Meshery's GraphQL API Reference</div>
<p>See <a href="/pr-preview/pr-21670/reference/references/graphql-apis/">GraphQL API Reference</a> Self-generated API documentation for Meshery&rsquo;s GraphQL API. The API can be explored interactively using the GraphQL Playground. Documentation is generated from Meshery&rsquo;s GraphQL schema. Each table below documents a GraphQL type.</p>
</div></section></div>

### Authentication

Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers. Type of authentication is determined by the selected [Provider](#providers). Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication.

<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">What are authentication tokens?</div>


Meshery authentication tokens allow users or systems to authenticate with Meshery Server via either its two clients, [Meshery CLI](/pr-preview/pr-21670/reference/references/mesheryctl/) and [UI](/pr-preview/pr-21670/reference/extensibility/api/#how-to-get-your-token), or its two APIs: [REST](/pr-preview/pr-21670/reference/references/rest-apis/) or [GraphQL](/pr-preview/pr-21670/reference/references/graphql-apis/). 

Meshery's authentication token system provide secure access to Meshery's powerful management features.
</div>


### How to get your token

There are two ways to get your authentication token:

1. Meshery UI
2. Meshery CLI

Using Meshery UI, you can get a copy of your authentication token by following these steps:

1. Log into Meshery by selecting your identity provider of choice (typically found at `http:<meshery-server>:9081/provider`)
2. Navigate to your user's avatar in the upper lefthand corner and select "Get Token" from the dropdown of profile section.

Using Meshery CLI, you can get a copy of your authentication token by executing this command:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system login</code>
	</div>
</pre>


In order to use this command, you must have a web browser available on your system (this command cannot be executed on a headless system).

<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">How to use the token for requests to Meshery Rest API using API clients</div>


Download your token from Meshery UI and copy the token value from the downloaded file, this will be used for authentication.
</div>



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">meshery-provider: Meshery
token: your token</code>
	</div>
</pre>

<div class="alert alert-dark" role="alert">


Navigate to the cookies section your API testing platform, to set the appropriate name and value for the cookie, then make sure to save the configured cookies to initiate the request. 

The request will be sent to the Meshery API using the configured cookies for authentication.
</div>


### Endpoints

Each of the API endpoints are exposed through [server.go](https://github.com/meshery/meshery/blob/master/server/router/server.go). Endpoints are grouped by function (e.g. `/api/mesh` or `/api/perf`).

Alternatively, [Remote Providers](/pr-preview/pr-21670/reference/extensibility/providers/) can extend Meshery's endpoints behind the `/api/extensions/` endpoint.

## Authorization

While Meshery only requires a valid token in order to allow clients to invoke its APIs, Remote Providers can optionally enforce key-based permissions.
