---
layout: default
title: Quick Start Guide
type: installation
permalink: installation
redirect_from: installation/
language: en
list: exclude
---

Installation procedures for deploying Meshery.

{% assign sorted_index = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_index %}
    {% if item.type=="installation" and item.language == "en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

{% include toc.html page=reference %}

{:toc}

{% include related-discussions.html tag="meshery" %}
