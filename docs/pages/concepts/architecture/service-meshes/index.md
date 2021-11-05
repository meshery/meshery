---
layout: default
title: Service Meshes
permalink: service-meshes
redirect_from:
  - service-meshes/adapters/
  - service-meshes/
type: service-mesh
list: exclude
---

As the multi-mesh manager, Meshery offers support for more adapters than any other project or product in the world. Learn more about [Meshery Adapters]({{ site.baseurl }}/concepts/architecture/adapters) in the Architecture section.

## Support for Service Meshes

| Service Mesh | Earliest Version |
| :----------: | :--------------: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.earliest_version }} |
{% endif -%}
{% endfor %}

Meshery supports the following service meshes. It uses both abstraction APIs and service mesh-specific adapters to interface with and manage service meshes. Review the full list of Meshery [adapters]({{ site.baseurl }}/architecture/adapters).
