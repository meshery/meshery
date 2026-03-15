---
layout: default
title: Infrastructure Management Guides
permalink: guides/infrastructure-management/
language: en
category: infrastructure
list: exclude
abstract: Guides for understanding Meshery's ability to manage infrastructure.
---

Guides for understanding Meshery's ability to manage infrastructure.

{% assign sorted_pages = site.pages | sort: "name" %}

<ul class="section-title">
    {% for item in sorted_pages %}
    {% if item.type=="guides" and item.category=="infrastructure" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    {% if item.abstract != " " %}
        -  {{ item.abstract }}
    {% endif %}
    </li>
    {% endif %}
    {% endfor %}
</ul>
