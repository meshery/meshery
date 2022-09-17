---
layout: default
title: Conceptos
permalink: es/concepts
redirect_from: es/concepts/
language: es
list: exclude
---

Conceptos para comprender las diversas características y componentes de Meshery.

{% assign sorted_pages = site.pages | sort: "name" %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="concepts" and item.language=="es" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>
