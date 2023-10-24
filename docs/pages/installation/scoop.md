---
layout: default
title: Install Meshery CLI with Scoop
permalink: installation/windows/scoop
type: installation
category: mesheryctl
display-title: "true"
language: en
list: include
image: /assets/img/platforms/scoop.png
---

{% include mesheryctl/installation-scoop.md %}

# Related Reading

## Mesherctl Guides

Guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "name" %}

<ul>
  {% for item in sorted_guides %}
  {% if item.type=="Guides" and item.category=="mesheryctl" and item.list!="exclude" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
    <li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>

