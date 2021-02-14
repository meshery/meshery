---
layout: default
title: Servis Ağları
permalink: tr/service-meshes
type: service-mesh
list: exclude
language: tr
lang: tr
categories: tr
---


As the multi-mesh manager, Meshery offers support for more adapters than any other project or product in the world.
## Support for Service Meshes

| Service Mesh  | Service Mesh Version  |
| :------------ | :------------: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.version }} |
{% endif -%}
{% endfor %}

Meshery supports the following service meshes. It uses both abstraction APIs and service mesh-specific adapters to interface with and manage service meshes. Review the full list of Meshery [adapters]({{ site.baseurl }}/tr/architecture/adapters).
