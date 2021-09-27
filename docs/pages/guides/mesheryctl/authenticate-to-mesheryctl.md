---
layout: default
title: Authenticating to `mesheryctl` 
permalink: guides/mesheryctl/authenticate-to-mesheryctl
language: en
type: Guides
category: mesheryctl
---

To authenticate to mesheryctl you have to use the command `mesheryctl system login`. On execution of this command you get a an option to choose a provider, you would have to choose `Meshery` as your provider to authenticate yourself. After doing so, a new window opens up with ways to login. 

## Get your Token

##### There are 2 ways to get your token:

- Get your token through **Meshery Dashboard**, from the `Get Token` option.

    _Downloading the token_

    <a href="{{ site.baseurl }}/assets/img/token/token.png"><img alt="Meshery Dashboard" src="{{ site.baseurl }}/assets/img/token/token.png" /></a>
    <br/>
    <br/>

- Get your token through **Meshery CLI**.
    <br/>
   > To get the token through `mesheryctl` you would have to use the following command and the path to token for authenticating to Meshery API (default "auth.json").
    <br/>
    <pre class="codeblock-pre">
        <div class="codeblock">
        mesheryctl system config --token [path-to-file]
        </div>
    </pre>


## Need for Authentication

The basic need for authentication to `Meshery` provider is:
- Persistent Sessions
- Save environment setup
- Retrieve performance test results 


# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "type" | reverse %}

<ul>
  <li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></lI>
  {% for item in sorted_guides %}
  {% if item.type=="Guides" and item.category=="mesheryctl" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>