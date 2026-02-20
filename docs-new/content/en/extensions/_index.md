---
title: Extensions
description: Extensions are additional plugins or add-ons that provide extra functionalities beyond the core features of the Meshery.
aliases: 
- extensibility/extensions
display_title: true
---

Extensions in Meshery are additional plugins or add-ons that provide extra functionalities beyond the core features of the platform. These extensions can be used to customize, extend and integrate with other tools and services.

{% assign sorted_pages = site.pages | sort: "name" %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="extensions" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

