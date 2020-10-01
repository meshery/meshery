---
layout: page
title: Meshery Architecture
permalink: architecture
---

# Architecture

<a href="{{site.baseurl}}/assets/img/architecture/meshery-architecture-diagram.svg"><img src="{{site.baseurl}}/assets/img/architecture/meshery-architecture-diagram.svg" /></a>

## Clients

<a href="{{site.baseurl}}/assets/img/architecture/meshery-client-architecture.svg"><img src="{{site.baseurl}}/assets/img/architecture/meshery-client-architecture.svg" /></a>

## Providers

<a href="{{site.baseurl}}/assets/img/architecture/meshery-provider-architecture.svg"><img src="{{site.baseurl}}/assets/img/architecture/meshery-provider-architecture.svg" /></a>

## Network Ports

Meshery uses the following list of network ports to interface with its various components:

| Adapter                                        | Port             |
| :--------------------------------------------- | :--------------- |
| Meshery REST API                               | 9081/tcp      
| Learn Layer5 Application                       | 10011 |
| {% assign adaptersSortedByPort = site.adapters | sort: 'port' -%} |

{% for adapter in adaptersSortedByPort -%}
{% if adapter.port -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [Adapters](/docs/architecture/adapters) section for more information on the function of an adapter.

## Statefulness in Meshery components

Some components within Meshery's architecture are concerned with persisting data while others are only
concerned with a long-lived configuration, while others have no state at all.

| Components        | Persistence  | Description                                                           |
| :---------------- | :----------- | :-------------------------------------------------------------------- |
| mesheryctl        | stateless    | command line interface that has a configuration file                  |
| Meshery Adapters  | stateless    | interface with service meshes on a transactional basis                |
| Meshery Server    | caches state | application cache is stored in user's `$HOME/.meshery/` folder        |
| Meshery Providers | stateful     | location of persistent user preferences, environment, tests and so on |
