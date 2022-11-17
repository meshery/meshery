---
layout: default
title: Malla de Servicios
permalink: es/service-meshes
type: service-mesh
list: exclude
language: es
# list: exclude
---

Como administrador de múltiples mallas, Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo.

## Soporte para mallas de servicio

| Service Mesh | Versión de Service Mesh |
| :----------- | :---------------------: |

{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.version }} |
{% endif -%}
{% endfor %}

Meshery admite las siguientes mallas de servicio. Utiliza tanto API de abstracción como adaptadores específicos de malla de servicio para interactuar con las mallas de servicio y administrarlas. Revise la lista completa de [adaptadores]({{ site.baseurl }}/es/architecture/adapters) Meshery .
