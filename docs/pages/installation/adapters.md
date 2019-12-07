---
layout: page
title: Adapters
permalink: installation/adapters
---

Meshery has adapters for managing the following service meshes.

| Platform      | Version       |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.version -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.version }} |
{% endif -%}
{% endfor %}
