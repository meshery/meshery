---
layout: default
title: Architecture
permalink: concepts/architecture
redirect_from: architecture
type: concepts
abstract: overview of different individual components of Meshery architecture and how they interact as a system.
language: en
list: include
---

## Architectural Components and Their Languages

Meshery and its components are written using the following languages and technologies.

| Components                                                           | Languages and Technologies                                                        |
| :------------------------------------------------------------------- | :-------------------------------------------------------------------------------- |
| Meshery Server                                                       | Golang, gRPC, GraphQL, [SMP](https://smp-spec.io), [SMI](https://smi-spec.io)     |
| [Meshery Adapters](/concepts/architecture/adapters)                  | Golang, gRPC, CloudEvents                                                         |
| Meshery WASM Filters                                                 | Rust and C++                                                                      |
| Meshery UI                                                           | ReactJS, NextJS, BillboardJS                                                      |
| Meshery Provider UI                                                  | ReactJS, NextJS                                                                   |
| [Meshery Remote Providers](/extensibility/providers)                 | _any_ - must adhere to Meshery [Extension Points]({{site.baseurl}}/extensibility) |
| [Meshery Operator](/concepts/architecture/operator)                  | Golang                                                                            |
| &nbsp;&nbsp;&nbsp;&nbsp; [MeshSync](/concepts/architecture/meshsync) | Golang                                                                            |
| &nbsp;&nbsp;&nbsp;&nbsp; [Broker](/concepts/architecture/broker)     | Golang, NATS                                                                      |
| [Meshery Database](/concepts/architecture/database)                  | Golang, SQLlite                                                                   |

## Deployments

Meshery deploys as a set of containers. Meshery's containers can be deployed to either Docker or Kubernetes. Meshery components connect to one another via gRPC requests. Meshery Server stores the location of the other components and connects with those components as needed. Typically, a connection from Meshery Server to Meshery Aapters is initiated from a client request (usually either `mesheryctl` or Meshery UI) to gather information from the Adapter or invoke an Adapter's operation.

### Adapters

In Meshery v0.6.0, Adapters will register with Meshery Server over HTTP POST. If Meshery Server is not available, Meshery Adapters will backoff and retry to connect to Meshery Server perpetually.

[![Meshery architecture]({{ site.baseurl }}/assets/img/architecture/meshery-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/meshery-architecture.svg)

_Figure: Meshery deploys inside or outside of a Kubernetes cluster_

#### Adapters and Capabilities Registry

Each Meshery Adapter delivers its own unique service mesh-specific functionality. As such, at time of deployment, the Meshery Adapter will register its service mesh-specific capabilities (its operations) with Meshery Server's capability registry.

[![Meshery Adapter Operation Registration]({{ site.baseurl }}/assets/img/adapters/meshery-adapter-operation-registration.svg
)]({{ site.baseurl }}/assets/img/adapters/meshery-adapter-operation-registration.svg)

_Figure: Service Mesh Adapter Operation Registration_

### Clients

Meshery's REST API may be consumed by any number of clients. Clients need to present valid JWT token.

[![Client architecture]({{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.svg)

_Figure: Clients use Meshery's [REST API](extensibility/api#rest), [GraphQL API](extensibility/api#graphql), or a combination of both._

### Providers

As a point of extension, Meshery supports two types of providers: _Local_ and _Remote_.

[![Provider architecture]({{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.svg)

## Object Model

This diagram outlines logical constructs within Meshery and their relationships.

[![Object Model]({{ site.baseurl }}/assets/img/architecture/meshery_logical_object_model.svg)]({{ site.baseurl }}/assets/img/architecture/meshery_logical_object_model.svg)

## Meshery Operator and MeshSync

Meshery Operator is the multi-service mesh operator (a Kubernetes custom controller) that manages MeshSync and it's messaging broker.

[![Meshery Operator and MeshSync]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg)

See the [**Operator**]({{ site.baseurl }}/concepts/architecture/operator) section for more information on the function of an operator and [**MeshSync**]({{ site.baseurl }}/concepts/architecture/meshsync) section for more information on the function of meshsync.

## Database

Meshery Server's database is responsible for collecting and centralizing the state of all elements under management, including infrastructure, application, and Meshery's own components. Meshery's database, while persisted to file, is treated as a cache.

[![Meshery Database]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)]({{ site.baseurl }}/concepts/architecture/database)

_See the [**Database**]({{ site.baseurl }}/concepts/architecture/database) section for more information on the function of the database._

### **Network Ports**

Meshery uses the following list of network ports to interface with its various components:

| Component                |                            Port                            |
| :----------------------- | :--------------------------------------------------------: |
| Meshery REST API         |                          9081/tcp                          |
| Meshery GraphQL          |                          9081/tcp                          |
| Meshery Broker           | 4222/tcp, 6222/tcp, 8222/tcp, 7777/tcp, 7422/tcp, 7522/tcp |
| Learn Layer5 Application |                         10011/tcp                          |
| Meshery Adapters         |                         10000+/tcp                         |
| Meshery Remote Providers |                          443/tcp                           |

### **Adapter Ports**

| Service Mesh | Port |
| :----------- | ---: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [**Adapters**]({{ site.baseurl }}/concepts/architecture/adapters) section for more information on the function of an adapter.

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
