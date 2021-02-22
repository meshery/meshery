---
layout: default
title: Referans
permalink: tr/reference
language: tr
lang: tr
categories: tr
list: exclude
---
References for using Meshery's various features and components.

{% assign sorted_reference = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_reference %}
    {% if item.type=="Reference"  and item.lang == "tr"  -%}
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