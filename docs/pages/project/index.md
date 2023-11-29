---
layout: default
title: Contributing and Community
type: project
category: contributing
abstract: Information about contributing to the Meshery project and participating the Meshery community.
permalink: project
redirect_from: project/
language: en
list: exclude
---

<div class="wrapper">
  <a href="/project/community"><div class="overview">Community</div></a>
</div>

{% assign sorted_pages = site.pages | sort: "name" | alphabetical %}

<ul>
  {% for item in sorted_pages %}
  {% if item.type=="project" and item.category!="contributing" and item.list=="include" and item.language =="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>

<div class="wrapper"> 
  <a href="/project/contributing"><div class="overview">Contributing</div></a>
</div>

{% assign sorted_index = site.pages | sort: "name" | alphabetical %}

<ul>
  {% for item in sorted_pages %}
  {% if item.type=="project" and item.category=="contributing" and item.list=="include" and item.language =="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
