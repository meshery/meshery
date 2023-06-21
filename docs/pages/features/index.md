---
layout: default
title: Features
permalink: features
redirect_from: features/
language: en
list: exclude
---

Features refer to the functionalities or capabilities that Meshery provides. Features are the specific tools and abilities that make Meshery a useful cloud native management tool for developers working with multiple clusters and service meshes.


{% assign sorted_pages = site.pages | sort: "name" %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="feature" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>
