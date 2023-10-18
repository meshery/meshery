---
layout: default
title: Overview
permalink: overview
redirect_from: overview/
language: en
display-title: "false"
list: exclude
---

# Meshery Overview

## Meshery as a project and its community

{% assign sorted_pages = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="project" and item.language=="en" and item.list != "exclude" %}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.description != " " %}
        -  {{ item.description }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>
