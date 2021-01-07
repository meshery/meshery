---
layout: default
title: Meshery Architecture
permalink: concepts/architecture
redirect_from: architecture
type: concepts
abstract: overview of different individual components of Meshery architecture and how they interact as a system.
language: en
list: include
---

## Architectural Components and Their Languages

Meshery and its components are written using the following languages and technologies.

| Components               | Languages and Technologies                                                        |
| :----------------------- | :-------------------------------------------------------------------------------- |
| Meshery Server           | Golang, gRPC, GraphQL, SQLlite                                                    |
| Meshery Adapters         | Golang, gRPC                                                                      |
| Meshery WASM Filters     | Rust and C++                                                                      |
| Meshery UI               | ReactJS, NextJS, BillboardJS                                                      |
| Meshery Provider UI      | ReactJS, NextJS                                                                   |
| Meshery Remote Providers | _any_ - must adhere to Meshery [Extension Points]({{site.baseurl}}/extensibility}}) |
| Meshery Operator         | Golang, NATS                                                                      |
| MeshSync                 | Golang                                                                            |

## Deployments

Meshery deploys as a set of containers. Meshery's containers can be deployed to either Docker or Kubernetes.

[![Meshery architecture]({{ site.baseurl }}/assets/img/architecture/meshery-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/meshery-architecture.svg)

## Clients

Meshery's REST API may be consumed by any number of clients. Clients need to present valid JWT token.

[![Client architecture]({{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.svg)

## Providers

As a point of extension, Meshery supports two types of providers: _Local_ and _Remote_.

[![Provider architecture]({{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.svg)

## Object Model

This diagram outlines logical constructs within Meshery and their relationships.

[![Object Model]({{ site.baseurl }}/assets/img/architecture/meshery_logical_object_model.svg)]({{ site.baseurl }}/assets/img/architecture/meshery_logical_object_model.svg)

## Meshery Operator and MeshSync

[![Meshery Operator and MeshSync]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg)

### **Network Ports**

Meshery uses the following list of network ports to interface with its various components:

| Component                |    Port    |
| :----------------------- | :--------: |
| Meshery REST API         |  9081/tcp  |
| Meshery GraphQL          |  9081/tcp  |
| Meshery Operator NATS    |  4222/tcp  |
| Learn Layer5 Application | 10011/tcp  |
| Meshery Adapters         | 10000+/tcp |
| Meshery Remote Providers |  443/tcp   |

### **Adapter Ports**

| Service Mesh | Port |
| :----------- | ---: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [**Adapters**]({{ site.baseurl }}/architecture/adapters) section for more information on the function of an adapter.

### **Statefulness in Meshery components**

Some components within Meshery's architecture are concerned with persisting data while others are only
concerned with a long-lived configuration, while others have no state at all.

| Components        | Persistence  | Description                                                           |
| :---------------- | :----------- | :-------------------------------------------------------------------- |
| mesheryctl        | stateless    | command line interface that has a configuration file                  |
| Meshery Adapters  | stateless    | interface with service meshes on a transactional basis                |
| Meshery Server    | caches state | application cache is stored in user's `$HOME/.meshery/` folder        |
| Meshery Providers | stateful     | location of persistent user preferences, environment, tests and so on |
| Meshery Operator  | stateless    | operator of Meshery custom controllers, notably MeshSync              |
| MeshSync          | stateless    | Kubernetes custom controller, continuously running discovery          |
