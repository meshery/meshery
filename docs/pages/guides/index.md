---
layout: guide
title: Operation Guides
permalink: guides
---
<ul>
    {% for item in site.pages %}
    {% if item.type=="guide" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

    
{% include toc.html page=guide %}

{:toc}