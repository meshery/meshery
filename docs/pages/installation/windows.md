---
layout: default
title: Install Meshery CLI on Windows
permalink: installation/windows
type: installation
category: mesheryctl
redirect_from:
- installation/platforms/windows
display-title: "true"
language: en
list: include
image: /assets/img/platforms/wsl2.png
---


On Windows systems, `mesheryctl` can be installed via Scoop or can be [downloaded directly](https://github.com/meshery/meshery/releases/latest).

{% include mesheryctl/installation-scoop.md %}

## Install `mesheryctl` as a direct download

Follow the [installation steps]({{ site.baseurl }}/installation#windows) to install the mesheryctl CLI. Then, execute:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">./mesheryctl system start</div></div>
</pre>

Optionally, move the mesheryctl binary to a directory in your PATH.


<!-- Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host ./mesheryctl system start</div></div>
  </pre>

Type `yes` when prompted to choose to configure a file. To get started, choose Docker as your platform to deploy Meshery. -->

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "name" %}

<ul>
  {% for item in sorted_guides %}
  {% if item.type=="guides" and item.category=="mesheryctl" and item.list!="exclude" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
    <li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="mesheryctl" %}

{:toc}