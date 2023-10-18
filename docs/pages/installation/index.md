---
layout: default
title: Quick Start Guide
type: installation
permalink: installation
redirect_from: installation/
language: en
list: exclude
---

## Supported Platforms

Installation procedures for deploying Meshery.

{% assign sorted_index = site.pages | sort: "name" | alphabetical %}

<ul>
    {% for item in sorted_index %}
    {% if item.type=="installation" and item.list=="include" and item.language == "en" -%}
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

{% include suggested-reading.html diffName="false" isDiffTag="true" diffTag=tag %}

{% include related-discussions.html tag="meshery" %}
