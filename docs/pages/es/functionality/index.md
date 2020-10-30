---
layout: page
title: Functionality
permalink: functionality
---

# Meshery Functionality

Estas secciones proveen guías de usuario sobre las funcionalidades que ofrece Meshery.

{% assign sorted_functionality = site.pages | sort: "functionality" %}

<ul>
    {% for item in sorted_functionality %}
    {% if item.type=="functionality" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=functionality %}

{:toc}
