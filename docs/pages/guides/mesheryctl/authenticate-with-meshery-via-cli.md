---
layout: enhanced
title: Authenticating with Meshery via CLI
permalink: guides/mesheryctl/authenticate-with-meshery-via-cli
language: en
type: guides
category: mesheryctl
list: include
abstract: Get your authentication token from Meshery CLI.
---

To authenticate with Meshery through `mesheryctl` you will use the command `mesheryctl system login`. Upon execution of this command, select your Provider of choice, then authenticate to your chosen Provider.

## Get your Token

You can retrieve your authentication token from either of Meshery's two clients: the CLI or the UI.

- Get your token through [Meshery UI](/extensibility/api#how-to-get-your-token), from the `Get Token` option.

  _Downloading the token_

  <a href="{{ site.baseurl }}/assets/img/token/MesheryTokenUI.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/token/MesheryTokenUI.png" /></a>
  <br/>
  <br/>

- Get your token through **Meshery CLI**.
  <br/>
  To get the token through `mesheryctl` you would have to use the following command and the path to token for authenticating to Meshery API (default "auth.json").
  <br/>
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs"> mesheryctl system login</div></div>
  </pre>
  <br />

**_The need for authentication to `Meshery` [provider](https://docs.meshery.io/extensibility/providers) is to save your environment setup while also having persistent/steady sessions and to be able to retrieve performance test results._**

<br/>
# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% capture tag %}

<li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></li>

{% endcapture %}

{% if page.suggested-reading != false and page.title and page.type and page.category and page.url %}
{% include_cached suggested-reading.html  title=page.title type=page.type category=page.category url=page.url language="en" %}
{% endif %}
{% include related-discussions.html tag="mesheryctl" %}

