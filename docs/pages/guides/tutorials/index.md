---
layout: default
title: Tutorials
display_title: true
type: guides
category: tutorials
language: en
list: exclude
permalink: guides/tutorials
redirect_from: guides/tutorials/
---


{% assign tutorials = site.pages | where: "category", "guides" %}

{% assign items_grouped = tutorials | group_by: 'model' %}
<ul>
  {% for group in items_grouped %}
    <h2>{{ group.name | capitalize }}</h2>
      <ul>
        {% for item in group.items %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endfor %}
      </ul>
  {% endfor %}
</ul>
