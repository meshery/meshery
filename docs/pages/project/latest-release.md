---
layout: blank
title: Latest Release
permalink: project/releases/latest
abstract: Return only the latest release of Meshery version information
# redirect_from: project/releases/latest/
language: en
type: project
category: none
list: exclude
---
{% assign sorted_release = site.releases | sort: 'date' | reverse %}
{%- for release in sorted_release limit:1 offset:0 -%}
  {{ release.tags }}
{% endfor %}
