---
layout: default
title: Kavramlar
permalink: tr/concepts
language: tr
lang: tr
categories: tr
type: concepts
list: include
# list: exclude
---
Concepts for understanding Meshery's various features and components.

{% assign sorted_pages = site.pages | sort: "name" | reverse %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="concepts"  and item.language=="tr"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=reference %}

{:toc}