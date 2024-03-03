---
layout: default
title: Tutorials
display_title: true
permalink: guides/tutorials
redirect_from: guides/tutorials/
---


{% assign tutorials = site.pages | where: "type", "tutorials" %}

{% assign items_grouped = tutorials | group_by: 'model' %}
<ul>
  {% for group in items_grouped %}
    <li>
      <h2>{{ group.name }}</h2>
      <ul>
        {% for item in group.items %}
          <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
        {% endfor %}
      </ul>
    </li>
  {% endfor %}
</ul>
