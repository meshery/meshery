---
layout: page
title: Contributing to Meshery Adapters
permalink: project/contributing-adapters
description: How to contribute to Meshery Adapters
language: en
type: project
category: contributing
---

**See [Extensibility: Service Mesh Adapters]({{site.baseurl}}/extensibility/adapters)**

Meshery Adapters are an Extension Point in Meshery's architecture. Their design and the process of creating a new adapter or contributing to an existing adapter is documented in [Extensibility: Service Mesh Adapters]({{site.baseurl}}/extensibility/adapters).

# Suggested Reading

{% assign sorted_reading = site.pages | sort: page.title | reverse %}

<ul>
  {% for item in sorted_reading %}
  {% if item.type=="project" and item.category=="contributing" and item.list!="exclude" and item.url!=page.url %}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
