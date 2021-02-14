---
layout: default
title: Conceptos
permalink: es/concepts
language: es
lang: es
categories: es
list: exclude
---
Conceptos para la comprensión de las diversas características y componentes de Meshery.

{% assign sorted_pages = site.pages | sort: "name" | reverse %}

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

{% include toc.html page=reference %}

{:toc}