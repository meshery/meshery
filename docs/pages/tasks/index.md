---
layout: default
title: Tasks
permalink: tasks
redirect_from: tasks/
language: en
---

These sections provide user guides to the functionality Meshery offers.
## Meshery Functionality

{% assign sorted_functionality = site.pages | sort: "functionality" %}

<ul>
    {% for item in sorted_functionality %}
    {% if item.type=="tasks" and item.list!="exclude" and item.language !="es"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

<!-- {% include toc.html page=functionality %} -->

<!-- {:toc} -->
