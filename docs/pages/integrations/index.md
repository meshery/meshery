---
layout: default
title: Integrations
permalink: /integrations
display-title: "false"
display-toc: "false"
language: en
---

{% assign sorted_pages = site.pages | where: "type", "integration" | sort: "name" | alphabetical %}
<ul>
{% for item in sorted_pages %}
    <!-- {% if item.type=="extensions" and item.list!="exclude" and item.language!="es" -%} -->
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    <!-- {% endif %} -->
{% endfor %}
</ul>