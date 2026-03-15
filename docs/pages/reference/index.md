---
layout: default
title: Reference
permalink: reference
redirect_from: reference/
language: en
list: exclude
abstract: References for using Meshery's various features and components.
---

References for using Meshery's various features and components.

{% assign sorted_reference = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_reference %}
    {% if item.type=="Reference" and item.list!="exclude" and item.language !="es" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>
{% include related-discussions.html tag="meshery" %}

