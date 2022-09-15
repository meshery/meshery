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

## Components, their Purpose, and Languages

Meshery and its components are written using the following languages and technologies.

| Components                                                           | Languages and Technologies                                                        |
| :------------------------------------------------------------------- | :-------------------------------------------------------------------------------- |
| Meshery Server                                                       | Golang, gRPC, GraphQL, [SMP](https://smp-spec.io)                                 |
| [Meshery Adapters](/concepts/architecture/adapters)                  | Golang, gRPC, [CloudEvents](https://cloudevents.io/), [SMI](https://smi-spec.io), [OAM](https://oam.dev)  |
| [Meshery WASM Filters](https://github.com/layer5io/wasm-filters)     | Rust and C++                                                                      |
| Meshery UI                                                           | ReactJS, NextJS, BillboardJS                                                      |
| Meshery Provider UI                                                  | ReactJS, NextJS                                                                   |
| [Meshery Remote Providers](/extensibility/providers)                 | _any_ - must adhere to Meshery [Extension Points]({{site.baseurl}}/extensibility) |
| [Meshery Operator](/concepts/architecture/operator)                  | Golang                                                                            |
| &nbsp;&nbsp;&nbsp;&nbsp; [MeshSync](/concepts/architecture/meshsync) | Golang                                                                            |
| &nbsp;&nbsp;&nbsp;&nbsp; [Broker](/concepts/architecture/broker)     | Golang, NATS                                                                      |
| [Meshery Database](/concepts/architecture/database)                  | Golang, SQLlite                                                                   |

## Deployments

Meshery deploys as a set of containers. Meshery's containers can be deployed to either Docker or Kubernetes. Meshery components connect to one another via gRPC requests. Meshery Server stores the location of the other components and connects with those components as needed. Typically, a connection from Meshery Server to Meshery Adapters is initiated from a client request (usually either `mesheryctl` or Meshery UI) to gather information from the Adapter or invoke an Adapter's operation.

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


### **Statefulness in Meshery components**

Some components within Meshery's architecture are concerned with persisting data while others are only
concerned with a long-lived configuration, while others have no state at all.

| Components        | Persistence  | Description                                                           |
| :---------------- | :----------- | :-------------------------------------------------------------------- |
| [mesheryctl]((/guides/mesheryctl/working-with-mesheryctl))        | stateless    | command line interface that has a configuration file                  |
| [Meshery Adapters](/concepts/architecture/adapters)  | stateless    | interface with service meshes on a transactional basis                |
| Meshery Server    | caches state | application cache is stored in `$HOME/.meshery/` folder               |
| [Meshery Providers](/extensibility/providers) | stateful     | location of persistent user preferences, environment, tests and so on |
| [Meshery Operator](/concepts/architecture/operator)  | stateless    | operator of Meshery custom controllers, notably MeshSync              |
| [MeshSync](/concepts/architecture/meshsync)          | stateless    | Kubernetes custom controller, continuously running discovery          |

### **Network Ports**

Meshery uses the following list of network ports to interface with its various components:

{% for adapter in site.adapters -%}
{% if adapter.port -%}
{% capture adapter-ports %}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }}/gRPC | Communication with Meshery Server |
{% endcapture %}
{% endif -%}
{% endfor %}

| Component                |   Port   | Purpose                                         |
| :----------------------- | :------: | :-----------------------------------------------|
| Meshery Server          | 9081/tcp | UI, REST and GraphQL APIs                           |
| [Meshery Broker](/concepts/architecture/broker)           | 4222/tcp | Client communication with Meshery Server        |
| [Meshery Broker](/concepts/architecture/broker)            | 8222/tcp | HTTP management port for monitoring Meshery Broker. Available as of Meshery v0.5.0 |
| [Meshery Broker](/concepts/architecture/broker)            | 6222/tcp | Routing port for Broker clustering. Unused as of Meshery v0.6.0-rc-2             |
| [Meshery Broker](/concepts/architecture/broker)            | 7422/tcp | Incoming/outgoing leaf node connections. Unused as of Meshery v0.6.0-rc-2 |
| [Meshery Broker](/concepts/architecture/broker)            | 7522/tcp | Gateway to gateway communication. Unused as of Meshery v0.6.0-rc-2 |
| [Meshery Broker](/concepts/architecture/broker)            | 7777/tcp | used for Prometheus NATS Exporter. Unused as of Meshery v0.6.0-rc-2 |
| Learn Layer5 Application | 10011/tcp  | SMI conformance testing                        |
| [Meshery Remote Providers]((/extensibility/providers)) | 443/tcp    | e.g. Meshery Cloud                             |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }} | Communication with Meshery Server |
{% endif -%}
{% endfor %}
| [Meshery Perf](https://docs.meshery.io/functionality/performance-management) | 10013/gRPC    | Performance Management                            |

See the [**Adapters**]({{ site.baseurl }}/concepts/architecture/adapters) section for more information on the function of an adapter.
