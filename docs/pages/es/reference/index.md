---
layout: default
title: Referencia
permalink: es/reference
language: es
list: exclude
---

Referencias para utilizar las diversas funciones y componentes de Meshery.

{% assign sorted_reference = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_reference %}
    {% if item.type=="Reference" and item.lang == "es" and item.list!="exclude" -%}

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
