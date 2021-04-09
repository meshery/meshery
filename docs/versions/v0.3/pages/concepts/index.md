---
layout: default
title: Concepts
permalink: /v0.3/concepts
language: en
list: exclude
---
Concepts for understanding Meshery's various features and components.

{% assign sorted_pages = site.pages | sort: "name" | reverse %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="concepts" and item.url contains '/v0.3/'-%}
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