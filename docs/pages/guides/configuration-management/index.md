---
layout: default
title: Configuration Management Guides
permalink: guides/configuration-management/
language: en
list: exclude
abstract: Guides for understanding Meshery's ability to configure infrastructure and applications
---

Guides for understanding Meshery's ability to configure infrastructure and applications.

{% assign sorted_pages = site.pages | sort: "name" %}

<ul class="section-title">
    {% for item in sorted_pages %}
    {% if item.type=="guides" and item.category=="configuration" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    {% if item.abstract != " " %}
        -  {{ item.abstract }}
    {% endif %}
    </li>
    {% endif %}
    {% endfor %}
</ul>
