---
layout: default
title: Install Meshery CLI on Linux or Mac
permalink: installation/linux-mac
type: installation
category: mesheryctl
redirect_from:
- installation/platforms/linux-mac
display-title: "false"
language: en
list: include
image: /assets/img/platforms/linux_mac.png 
abstract: Install Meshery CLI on Linux or Mac
---

# Overview

To set up and run Meshery on Linux or macOS, you will need to install `mesheryctl`. `mesheryctl` is the command line interface (CLI) for Meshery. It is used to install, manage, and operate one or more Meshery deployments. `mesheryctl` can be installed via `bash` is also available [directly](https://github.com/meshery/meshery/releases/latest) or through [Homebrew]({{site.baseurl}}/installation/linux-mac/brew) or [Scoop]({{site.baseurl}}/installation/windows/scoop).

# Brew

{% include mesheryctl/installation-brew.md %}

# Bash

{% include mesheryctl/installation-bash.md %}

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

{% if page.suggested-reading != false and page.title and page.type and page.category and page.url %}
{% include_cached suggested-reading.html  title=page.title type=page.type category=page.category url=page.url language="en" %}
{% endif %}

{% include related-discussions.html tag="mesheryctl" %}

{:toc}

<!-- 
1. You can either use **Bash** or **Brew** to install <a href="/guides/mesheryctl">mesheryctl</a> ( Meshery command line interface ).
2. To run **Meshery**, execute the following command.

   <pre class="codeblock-pre"><div class="codeblock">
   <div class="clipboardjs">mesheryctl system start</div></div>
   </pre>

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way

<pre class="codeblock-pre"><div class="codeblock">
<div class="clipboardjs">
 $ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start

</div></div>
</pre>
-->