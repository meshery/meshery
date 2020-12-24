---
layout: default
title: Releases
permalink: project/releases
---
{% for collection in site.releases reversed %}
{{ collection.content | markdownify  }}
{% endfor %}
