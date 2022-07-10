---
layout: default
title: Reference
permalink: reference
redirect_from: reference/
language: en
list: exclude
---
References for using Meshery's various features and components.

{% assign sorted_reference = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_reference %}
    {% if item.type=="Reference" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=reference %}

{:toc}