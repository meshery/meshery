---
layout: default
title: Releases
permalink: project/releases
---
{% assign sorted_release = site.releases | sort: 'date' | reverse %}
{% for collection in sorted_release %}
{{ collection.content | markdownify  }}
{% endfor %}
