---
layout: page
title: Funcionalidad
permalink: es/functionality
type: functionality
language: es
list: exclude
---

Estas secciones proveen guías de usuario sobre las funcionalidades que ofrece Meshery.

{% assign sorted_functionality = site.pages | sort: "functionality" %}

<ul>
    {% for item in sorted_functionality %}
      {% if item.type=="functionality" and item.language=="es" and item.list!="exclude" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=functionality %}

{:toc}

