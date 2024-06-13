---
layout: default
title: Contributing and Community
type: project
category: contributing
abstract: Information about contributing to the Meshery project and participating the Meshery community.
permalink: project
suggested-reading: false
# redirect_from: project/
language: en
list: exclude
---
<!-- PROJECT 
{% assign project = site.pages | sort: "name" | alphabetical %}
<ul>
  {% for item in project %}
  {% if item.type=="project" and item.category!="contributing" and item.list=="include" and  item.list!="exclude" and item.language =="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
-->
<!-- CONTRIBUTING 
<details>
  <summary>
    <p style="display:inline">
      <a href="{{ site.baseurl }}/project/contributing" class="text-black">Contributing</a>
    </p>
  </summary>
  <ul class="section-title">
    {% assign contributing = site.pages | where: "category","contributing" %}
      {% for item in contributing %}
      {% if item.category=="contributing" and item.language=="en" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        {% if item.abstract != " " %}
          - {{ item.abstract }}
        {% endif %}
        </li>
        {% endif %}
      {% endfor %}
  </ul>
</details>
-->

{% assign sorted_pages = site.pages | sort: "name" | alphabetical %}

<div class="wrapper">
  <a href="/project/community"><div class="overview">Community</div></a>
</div>

<ul>
  {% for item in sorted_pages %}
  {% if item.type=="project" and item.category=="community" and item.list=="include" and item.language =="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>

<div class="wrapper">
  <a href="/project/community"><div class="overview">Project</div></a>
</div>

<ul>
  {% for item in sorted_pages %}
  {% if item.type=="project" or item.type=="installation" %}
    {% if item.category=="project" and item.list=="include" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      </li>
      {% endif %}
    {% endif %}
  {% endfor %}
</ul>

<div class="wrapper"> 
  <a href="/project/contributing"><div class="overview">Contributing Guides</div></a>
</div>

<ul>
  {% for item in sorted_pages %}
  {% if item.type=="project" and item.category=="contributing" and item.list=="include" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
