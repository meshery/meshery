---
layout: page
title: Service Meshes
permalink: service-meshes
---

Meshery supports the following service meshes. Meshery uses both abstraction APIs and service mesh-specific adapters to interface with and manage service meshes. You may review the full list of Meshery [adapters](adapters).

| Service Mesh      | Version       |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.version }} |
{% endif -%}
{% endfor %}
