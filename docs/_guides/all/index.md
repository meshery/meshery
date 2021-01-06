---
layout: default
title: Operation Guides
language: en
exclude: true
---

Guides to using Meshery's various features and components. 

{% assign sorted_guides = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_guides %}
    {% if item.type=="Guides" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a></li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=Guides %}

{:toc}

<!-- {% comment %}
#
#  Change date order by adding '| reversed'
#  To sort by title or other variables use {% assign sorted_posts = category[1] | sort: 'title' %}
#
{% endcomment %}
