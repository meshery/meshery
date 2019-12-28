---
layout: page
title: Architecture
permalink: architecture
---

# Architecture

<iframe class="container" src="https://docs.google.com/presentation/d/e/2PACX-1vSj6eYr6AgZ4mBgOL_Gv9T4WyLBFkPv49asNtdw1_Gn_xCsk37QRhOjdBRB-3Jp1ehneFmm2dpgFie-/embed?start=false&loop=false&delayms=3000#slide=id.g55c4016581_0_0" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>

# Value provided by a Provider
Meshery Providers offer different things, but each must at least offer some form of authentication (including no authentication).

For the Provider, Meshery Cloud, performance test reports area anonymously sent to Meshery Cloud for statistical analysis and sharing back to the community. Data points like the average overhead of service meshes as seen in various environments are shared through this mechanism.

<iframe class="container" src="https://docs.google.com/presentation/d/e/2PACX-1vSj6eYr6AgZ4mBgOL_Gv9T4WyLBFkPv49asNtdw1_Gn_xCsk37QRhOjdBRB-3Jp1ehneFmm2dpgFie-/embed?start=false&loop=false&delayms=3000#slide=id.g4f68f671f0_0_0" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>


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

See the [Adapters](installation/adapters) section for more information on the function of an adapter.
