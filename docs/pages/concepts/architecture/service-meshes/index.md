---
layout: default
title: Service Meshes
permalink: service-meshes
type: service-mesh
---
## Support for Service Meshes
As the multi-mesh manager, Meshery offers support for more adapters than any other project or product in the world.

| Service Mesh  | Service Mesh Version  |
| :------------ | :------------: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.version }} |
{% endif -%}
{% endfor %}

Meshery supports the following service meshes. It uses both abstraction APIs and service mesh-specific adapters to interface with and manage service meshes. Review the full list of Meshery [adapters](/docs/architecture/adapters).
