---
layout: default
title: Documentation in Espanol
permalink: es
---
{% assign sorted_pages = site.pages | sort: "name" %}

## Installation

<ul>
    {% for item in sorted_pages %}
    {% if item.language=="espanol" and item.type=="installation" %}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

## Contribuir

<ul>
{% for item in sorted_pages %}
    {% if item.language=="espanol" and item.type=="project" %}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
    {% endif %}
{% endfor %}
<ul>

{% include toc.html page=espanol %}