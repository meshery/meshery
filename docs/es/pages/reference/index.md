---
layout: default
title: Referencia
permalink: es/reference
language: es
lang: es
categories: es
list: exclude
---
Referencias para utilizar las diversas funciones y componentes de Meshery.

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