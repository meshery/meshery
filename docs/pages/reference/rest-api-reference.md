---
layout: swagger
title: "REST API Reference"
permalink: reference/rest-apis
redirect_from: reference/rest-apis/
type: Reference
abstract: 'Meshery REST API Documentation and Reference'
data: swagger
language: en
---
## Meshery's APIs

Each of Meshery's APIs are subject to the following authentication and authorization system.
Currently, Meshery requires a valid token in order to allow clients to invoke its APIs.

{% include alert.html type="dark" title="What are authentication tokens?" content="Meshery authentication tokens allow users or systems to authenticate with Meshery Server via either its two clients, <a href='/reference/mesheryctl'>Meshery >CLI</a> and <a href='/extensibility/api#how-to-get-your-token'>UI</a>, or its two APIs: <a href='/reference/rest-apis'>REST</a> or <a href='/reference/graphql-apis'>GraphQL</a>. <p>Meshery's authentication token system provide secure access to Meshery's powerful management features.</p>" %}

<details>
  <summary>Authentication</summary>
  Requests to any of the API endpoints must be authenticated and include a valid JWT access token in the HTTP headers. Type of authentication is determined by the selected [Provider](#providers). Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication.
</details>

### How to get your token

There are two ways to get your authentication token:

1. Meshery UI
2. Meshery CLI

<details>
  <summary>Meshery UI</summary>

Using Meshery UI, you can get a copy of your authentication token by following these steps:

1. Log into Meshery by selecting your identity provider of choice (typically found at <pre style="background-color:#343a40;">http:\\{meshery-server}:9081/provider</pre>)
2. Navigate to your user's avatar in the upper lefthand corner and select "Get Token" from the dropdown of profile section.

</details>

<details>
  <summary>Meshery CLI</summary>
Using Meshery CLI, you can get a copy of your authentication token by executing this command:

- Get your token through [Meshery UI](/extensibility/api#how-to-get-your-token), from the `Get Token` option.

  _Downloading the token_

  <a href="{{ site.baseurl }}/assets/img/token/MesheryTokenUI.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/token/MesheryTokenUI.png" width="680" height="400"/></a>
  <br/>
  <br/>

- Get your token through **Meshery CLI**.
  <br/>
  To get the token through `mesheryctl` you would have to use the following command. (default "auth.json").
  <br/>
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">mesheryctl system login</div></div>
  </pre>
  <br />

</details>


### Self-generated documentation based on Meshery's OpenAPI specification for it's REST API.  
Meshery's REST API can be explored interactively using the Swagger UI Playground.

## Endpoints
