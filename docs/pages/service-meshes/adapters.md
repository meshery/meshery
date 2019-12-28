---
layout: page
title: Adapters
permalink: service-meshes/adapters
---

Meshery has adapters for managing the following service meshes.

| Platform      | Status        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
