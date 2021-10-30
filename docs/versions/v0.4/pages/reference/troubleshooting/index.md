---
layout: default
title: Troubleshooting Guide
permalink: /v0.4/guides/troubleshooting
redirect_from: guides/troubleshooting/
language: en
---

Troubleshooting guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_guides %}
    {% if item.type=="Guides" and item.category=="troubleshooting" and item.list!="exclude" and item.url contains '/v0.4/'  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

#### See Also

- [Error Code Reference](/reference/error-codes)
<!-- {:toc} -->