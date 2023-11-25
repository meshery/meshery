---
layout: default
title: Scoop
permalink: installation/windows/scoop
type: installation
category: mesheryctl
redirect_from:
- installation/platforms/scoop
display-title: "false"
language: en
list: include
image: /assets/img/platforms/scoop.png
---
# Install Meshery CLI with Scoop

{% include mesheryctl/installation-scoop.md %}

# Related Reading

## Mesherctl Guides

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

