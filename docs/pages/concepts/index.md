---
layout: default
title: Concepts
permalink: concepts
# redirect_from: concepts/
language: en
list: exclude
abstract: Concepts for understanding Meshery's various features and components.
---

Concepts for understanding Meshery's various features and components.

{% assign sorted_pages = site.pages | sort: "name" %}

<h2>Architectural Components</h2>
<ul>
    {% for item in sorted_pages %}
    {% if item.type=="components" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

<h2>Logical Concepts</h2>
<ul>
    {% for item in sorted_pages %}
    {% if item.type=="concepts" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>