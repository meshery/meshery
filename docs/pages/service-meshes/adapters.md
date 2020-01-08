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

### Running more than one instance of the same Meshery adapter
The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the `meshery-istio` adapter. To do so, modify ~/.meshery/meshery.yaml to include multiple copies of the given adapter.