---
layout: swagger
title: "REST API Reference"
permalink: reference/rest-apis
redirect_from: reference/rest-apis/
type: Reference
abstract: "Meshery REST API Documentation and Reference"
data: swagger
language: en
---

## Meshery's APIs

Each of Meshery's APIs are subject to the following authentication and authorization system. Meshery requires a valid token in order to allow clients to invoke its APIs.

<details>
  <summary>Authentication</summary>
  Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers. Type of authentication is determined by the selected [Provider](#providers). Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication. {% include alert.html type="dark" title="What are authentication tokens?" content="Meshery authentication tokens allow users or systems to authenticate with Meshery Server via either its two clients, <a href='/reference/mesheryctl'>Meshery CLI</a> and <a href='/extensibility/api#how-to-get-your-token'>Meshery UI</a>, or its two APIs: <a href='/reference/rest-apis'>REST</a> or <a href='/reference/graphql-apis'>GraphQL</a>. <p>Meshery's authentication token system provide secure access to Meshery's powerful management features.</p>" %}
</details>

### How to get your token

There are two ways to get your authentication token:

<details>
  <summary>Meshery UI</summary>

Using Meshery UI, you can get a copy of your authentication token by following these steps:
<br/>

1. Log into Meshery by selecting your identity provider of choice (typically found at <code style="
       color: inherit;
       padding: 0.2em 0.4em;
       margin: 0;
       font-size: 85%;
       word-break: normal;
       background-color: var(--color-primary-dark);
       border-radius: 0.25rem;
       ">http:\\{meshery-server}:9081/provider</code>)
   <br/>

2. Navigate to your user's avatar in the upper righthand corner and select "Get Token" from the dropdown of profile section:

<a href="{{ site.baseurl }}/assets/img/token/MesheryTokenUI.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/token/MesheryTokenUI.png" width="680" height="400"/></a>
<br/>

</details>

<details>
  <summary>Meshery CLI</summary>
  <br />
Using <a href='/reference/mesheryctl'>Meshery CLI</a>, you can get a copy of your authentication token by executing this command:
  <br/>
  <br/>
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">mesheryctl system login</div></div>
  </pre>
  <br />
  <br />
  In order to use this command, you must have a web browser available on your system (this command cannot be executed on a headless system).

</details>

### How to access Meshery's REST API

<details>
  <summary>Example using curl</summary>
  <br />
Using curl, you can access Meshery's REST API by executing this command:
  <br/>
  <br/>
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">curl --location 'localhost:9081/api/&lt;endpoint&gt;' \
--header 'meshery-token: &lt;yourToken&gt;\
--header 'Cookie: meshery-provider=Meshery; meshery.layer5.io_ref=/;token=&lt;yourToken&gt;
</div>
</div>
  </pre>
  <br />
  <br />

</details>

### Self-generated documentation based on Meshery's OpenAPI specification for it's REST API.

Meshery's REST API can be explored interactively using the Swagger UI Playground.

## Endpoints
