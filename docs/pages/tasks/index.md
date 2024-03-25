---
layout: default
title: Tasks
permalink: tasks
# redirect_from: tasks/
language: en
list: exclude
published: false
abstract: Common tasks for using Meshery's various features and components.
---

These sections provide user guides to the functionality Meshery offers.

## Meshery Functionality

{% assign sorted = site.pages | sort: "functionality" %}

<ul>
    {% for item in sorted %}
    {% if item.type=="tasks" and item.list!="exclude" and item.language !="es"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>