---
layout: page
title: Adapters
permalink: installation/adapters
---

Meshery has adapters for managing the following service meshes.

| Platform      | Version       |
| -------------:|:------------- |
{% for adapter in site.adapters -%}
| [{{ adapter.name}} ]({{ adapter.url }}) | {{ adapter.version }} |
{% endfor -%}
