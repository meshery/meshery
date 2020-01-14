---
layout: page
title: Meshery Architecture
permalink: architecture
---

# Architecture

<iframe class="container" src="https://layer5.io/assets/images/meshery/meshery-architecture.svg" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>

## Network Ports 
Meshery uses the following list of network ports to interface with its various components:

| Adapter       | Port          |
| :------------ | :------------ |
| Meshery web-based UI | 9081/tcp |
{% assign adaptersSortedByPort = site.adapters | sort: 'port' -%}
{% for adapter in adaptersSortedByPort -%}
{% if adapter.port -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [Adapters](service-meshes/adapters) section for more information on the function of an adapter.

