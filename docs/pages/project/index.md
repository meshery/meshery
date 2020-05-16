---
layout: page
title: Project
permalink: project
---
{% assign sorted_project = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_project %}
    {% if item.type=="project" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

    
{% include toc.html page=project %}

{:toc}