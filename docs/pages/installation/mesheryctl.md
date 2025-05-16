---
layout: default
title: Install mesheryctl
permalink: installation/mesheryctl
type: installation
category: mesheryctl
redirect_from:
- installation/mesheryctl/
- installation/platforms/mesheryctl
display-title: "true"
language: en
list: exclude
suggested-reading: false
# image: /assets/img/platforms/brew.png
abstract: Install Meshery CLI
---

Meshery's command line client is `mesheryctl` and is the recommended tool for configuring and deploying one or more Meshery deployments. To install `mesheryctl` on your system, you may choose from any of the following supported methods.

`mesheryctl` can be installed via [bash]({{site.baseurl}}/installation/linux-mac/bash), [Homebrew]({{site.baseurl}}/installation/linux-mac/brew), [Scoop]({{site.baseurl}}/installation/windows/scoop) or [directly downloaded](https://github.com/meshery/meshery/releases/latest).

# Install Meshery CLI with Brew

{% include mesheryctl/installation-brew.md %}

# Install Meshery CLI with Bash

{% include mesheryctl/installation-bash.md %}

# Install Meshery CLI with Scoop

{% include mesheryctl/installation-scoop.md %}

Continue deploying Meshery onto one of the [Supported Platforms]({{ site.baseurl }}/installation).

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
    <li><a href="{{ site.baseurl }}/installation/upgrades#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>

{% include related-discussions.html tag="mesheryctl" %}

{:toc}
