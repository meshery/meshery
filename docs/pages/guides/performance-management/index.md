---
layout: default
title: Performance Management Guides
permalink: guides/performance-management
language: en
list: exclude
abstract: Guides for understanding Meshery's ability to characterize and analyze performance of your infrastructure and applications.
---

Guides for understanding Meshery's ability to characterize and analyze performance of your infrastructure and applications.

{% assign sorted_pages = site.pages | sort: "name" %}

<ul class="section-title">
    {% for item in sorted_pages %}
    {% if item.type=="guides" and item.category=="performance" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    {% if item.abstract != " " %}
        -  {{ item.abstract }}
    {% endif %}
    </li>
    {% endif %}
    {% endfor %}
</ul>
