---
layout: default
title: Authenticating with Meshery via CLI
permalink: guides/mesheryctl/authenticate-with-meshery-via-cli
language: en
type: Guides
category: mesheryctl
---

To authenticate with Meshery through `mesheryctl` you will use the command `mesheryctl system login`. Upon execution of this command, select your Provider of choice, then authenticate to your chosen Provider.

## Get your Token

You can retrieve your authentication token from either of Meshery's two clients: the CLI or the UI.

- Get your token through [Meshery UI](/concepts/architecture/ui), from the `Get Token` option.

  _Downloading the token_

  <a href="{{ site.baseurl }}/assets/img/token/token.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/token/token.png" /></a>
  <br/>
  <br/>

- Get your token through **Meshery CLI**.
  <br/>
  To get the token through `mesheryctl` you would have to use the following command and the path to token for authenticating to Meshery API (default "auth.json").
  <br/>
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">
  mesheryctl system config [aks|eks|gke|minikube] --token [path-to-file]
  
  </div></div>
  </pre>
  <br />

**_The need for authentication to `Meshery` [provider](http://localhost:9081/provider) is to save your environment setup while also having persistent/steady sessions and to be able to retrieve performance test results._**

<br/>
# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% capture tag %}

<li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></li>

{% endcapture %}

{% include suggested-reading.html diffName="true" isDiffTag="true" diffTag=tag %}
