---
layout: page
title: Meshery Architecture
permalink: architecture
---

# Architecture

<a href="{{site.baseurl}}/assets/images/meshery-architecture.svg"><img src="{{site.baseurl}}/assets/images/meshery-architecture.svg" /></a>

## Network Ports 
Meshery uses the following list of network ports to interface with its various components:

| Adapter       | Port          |
| :------------ | :------------ |
| Meshery REST API | 9081/tcp |
{% assign adaptersSortedByPort = site.adapters | sort: 'port' -%}
{% for adapter in adaptersSortedByPort -%}
{% if adapter.port -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [Adapters](service-meshes/adapters) section for more information on the function of an adapter.

