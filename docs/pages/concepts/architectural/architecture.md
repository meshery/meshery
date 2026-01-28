---
layout: default
title: Architecture
permalink: concepts/architecture
redirect_from: architecture
type: components
abstract: overview of different individual components of Meshery architecture and how they interact as a system.
language: en
list: include
---

## Components, their Purpose, and Languages

Meshery and its components are written using the following languages and technologies.

| Components                                                           | Languages and Technologies                                                        |
| :------------------------------------------------------------------- | :-------------------------------------------------------------------------------- |
| Meshery Server                                                       | Golang, gRPC, GraphQL, [SMP](https://smp-spec.io)                                 |
|   [Meshery Database](/concepts/architecture/database)                | Golang, SQLite                                                                   |
| Meshery UI                                                           | ReactJS, NextJS, BillboardJS                                                      |
| Meshery Provider UI                                                  | ReactJS, NextJS                                                                   |
| [Meshery Operator](/concepts/architecture/operator)                  | Golang                                                                            |
|   [MeshSync](/concepts/architecture/meshsync)                        | Golang                                                                            |
|   [Broker](/concepts/architecture/broker)                            | Golang, NATS                                                                      |
| [Meshery CLI](#meshery-cli)                                          | Golang                                                                            |
| --- [Extensions](/extensions) ---                                    |                                                                                   |
| [Meshery Adapters](/concepts/architecture/adapters)                  | Golang, gRPC, [CloudEvents](https://cloudevents.io/)                              |
| [Meshery Remote Providers](/extensibility/providers)                 | _any_ - must adhere to Meshery [Extension Points]({{site.baseurl}}/extensibility) |
| [Envoy WASM Filters](https://github.com/layer5io/wasm-filters)     | Rust and C++                                                                      |

## Deployments

Meshery deploys as a set of containers. Meshery's containers can be deployed to either Docker or Kubernetes. Meshery components connect to one another via gRPC requests. Meshery Server stores the location of the other components and connects with those components as needed. Typically, a connection from Meshery Server to Meshery Adapters is initiated from a client request (usually either `mesheryctl` or Meshery UI) to gather information from the Adapter or invoke an Adapter's operation.

### Adapters

In Meshery v0.6.0, Adapters will register with Meshery Server over HTTP POST. If Meshery Server is not available, Meshery Adapters will backoff and retry to connect to Meshery Server perpetually.

<a href="{{ site.baseurl }}/assets/img/architecture/meshery-architecture.webp" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/meshery-architecture.webp" width="50%" /></a>

_Figure: Meshery deploys inside or outside of a Kubernetes cluster_

#### Adapters and Capabilities Registry

Each Meshery Adapter delivers its own unique specific functionality. As such, at time of deployment, the Meshery Adapter will register its cloud native infrastructure-specific capabilities (its operations) with Meshery Server's capability registry.

<a href="{{ site.baseurl }}/assets/img/adapters/meshery-adapter-operation-registration.svg" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/adapters/meshery-adapter-operation-registration.svg" width="50%" /></a>

_Figure: Meshery Adapter Operation Registration_

### Clients

Meshery's REST API may be consumed by any number of clients. Clients need to present valid JWT token.

<a href="{{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.webp" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.webp" width="50%" /></a>

_Figure: Clients use Meshery's [REST API](/extensibility/api#rest), [GraphQL API](/extensibility/api#graphql), or a combination of both._

### Providers

As a point of extensibility, Meshery supports two types of [providers](/extensibility/providers): _Local_ and _Remote_.

<a href="{{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.webp" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.webp" width="50%" /></a>
<figure>
  <figcaption>Figure: Meshery Provider architecture</figcaption>
</figure>

## Object Model

This diagram outlines logical constructs within Meshery and their relationships.

<a href="{{ site.baseurl }}/assets/img/architecture/meshery_extension_points.svg" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/meshery_extension_points.svg" width="50%" /></a>
<figure>
  <figcaption>Figure: Meshery Object Model</figcaption>
</figure>

## Meshery Operator and MeshSync

Meshery Operator is the multi-cluster Kubernetes operator that manages MeshSync and Meshery Broker.

<a href="{{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg" width="50%" /></a>
<figure>
  <figcaption>Figure: Meshery Operator and MeshSync</figcaption>
</figure>

See the [**Operator**]({{ site.baseurl }}/concepts/architecture/operator) section for more information on the function of an operator and [**MeshSync**]({{ site.baseurl }}/concepts/architecture/meshsync) section for more information on the function of meshsync.

## Database

Meshery Server's database is responsible for collecting and centralizing the state of all elements under management, including infrastructure, application, and Meshery's own components. Meshery's database, while persisted to file, is treated as a cache.

<a href="{{ site.baseurl }}/assets/img/architecture/meshery-database.webp" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/meshery-database.webp" width="50%" /></a>
<figure>
  <figcaption>Figure: Meshery Docker Extension</figcaption>
</figure>

_See the [**Database**]({{ site.baseurl }}/concepts/architecture/database) section for more information on the function of the database._

## Meshery Docker Extension

Meshery's Docker extension provides a simple and flexible way to design and operate cloud native infrastructure on top of Kubernetes using Docker containers. The architecture of this extension is designed to be modular and extensible, with each component serving a specific purpose within the overall deployment process.

<a href="{{ site.baseurl }}/assets/img/architecture/meshery-docker-extension.svg" class="lightbox-image">
<img src="{{ site.baseurl }}/assets/img/architecture/meshery-docker-extension.svg" width="50%" /></a>
<figure>
  <figcaption>Figure: Meshery Docker Extension</figcaption>
</figure>

## Meshery CLI

The Command Line Interface ( also known as [mesheryctl](/guides/mesheryctl/working-with-mesheryctl) ) that is used to manage Meshery. Use `mesheryctl` to both manage the lifecycle of Meshery itself and to access and invoke any of Meshery's application and cloud native management functions.

### **Statefulness in Meshery components**

Some components within Meshery's architecture are concerned with persisting data while others are only
concerned with a long-lived configuration, while others have no state at all.

| Components        | Persistence  | Description                                                           |
| :---------------- | :----------- | :-------------------------------------------------------------------- |
| [mesheryctl](/guides/mesheryctl/working-with-mesheryctl)        | stateless    | command line interface that has a configuration file                  |
| [Meshery Adapters](/concepts/architecture/adapters)  | stateless    | interface with cloud native infrastructure on a transactional basis                |
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
| Meshery Server          | 80/tcp | Websocket                          |
| [Meshery Broker](/concepts/architecture/broker)           | 4222/tcp | Client communication with Meshery Server        |
| [Meshery Broker](/concepts/architecture/broker)            | 8222/tcp | HTTP management port for monitoring Meshery Broker. Available as of Meshery v0.5.0 |
| [Meshery Broker](/concepts/architecture/broker)            | 6222/tcp | Routing port for Broker clustering. Unused as of Meshery v0.6.0-rc-2             |
| [Meshery Broker](/concepts/architecture/broker)            | 7422/tcp | Incoming/outgoing leaf node connections. Unused as of Meshery v0.6.0-rc-2 |
| [Meshery Broker](/concepts/architecture/broker)            | 7522/tcp | Gateway to gateway communication. Unused as of Meshery v0.6.0-rc-2 |
| [Meshery Broker](/concepts/architecture/broker)            | 7777/tcp | used for Prometheus NATS Exporter. Unused as of Meshery v0.6.0-rc-2 |
| [Meshery Remote Providers](/extensibility/providers)      | 443/tcp    | e.g. Meshery Cloud                             |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" data-logo-for-dark="{{ adapter.white_image }}" data-logo-for-light="{{ adapter.image }}" id="logo-dark-light" loading="lazy"/> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }} | Communication with Meshery Server |
{% endif -%}
{% endfor -%}
| [Meshery Perf]({{ site.baseurl }}/guides/performance-management/managing-performance) | 10013/gRPC    | Performance Management|

See the [**Adapters**]({{ site.baseurl }}/concepts/architecture/adapters) section for more information on the function of an adapter.

### **Meshery Connections and their Actions**

| Connection Type | **Connect mesheryctl** | **Connect Meshery UI** | **Disconnect** | **Ad hoc Connectivity Test** | **Ongoing Connectivity Test** | **Synthetic Check** | **Deploy mesheryctl** | **Undeploy mesheryctl** | **Deploy Meshery UI** | **Undeploy Meshery UI** |
|---|---|---|---|---|---|---|---|---|---|---|---|
| Kubernetes clusters | `system start` | Upload kubeconfig | Click "X" on chip | On click of connection chip | Yes, via MeshSync | No | No | No | No | No |
| Grafana Servers | No | Enter IP/hostname into Meshery UI | Click "X" on chip | On click of connection chip | No | No | No | No | No | No |
| Prometheus Servers | No | Enter IP/hostname into Meshery UI | Click "X" on chip | On click of connection chip | Yes, when metrics are configured in a dashboard | Yes | No | No | No | No |
| [Meshery Adapters](/concepts/architecture/adapters) | `system check` | Server to Adapter on every UI refresh | Click "X on" chip | Server to Adapter every click on adapter chip in UI | Server to Adapter every 10 seconds | - | Yes, as listed in meshconfig contexts | Yes, as listed in meshconfig contexts | Toggle switch needed | Toggle switch needed |
| [Meshery Operator](/concepts/architecture/operator) | `system check` | Upon upload of kubeconfig | No | On click of connection chip in UI to Server to Kubernetes to Meshery Operator | No | - | `system start` | `system stop` | Upon upload of kubeconfig & Toggle of switch | Toggle of switch |
| [MeshSync](/concepts/architecture/meshsync) | `system check` | follows the lifecycle of Meshery Operator | No | On click of connection chip in UI to Server to Kubernetes to Meshery Operator to MeshSync | Managed by Meshery Operator | On click of connection chip | follows the lifecycle of Meshery Operator | follows the lifecycle of Meshery Operator | follows the lifecycle of Meshery Operator | follows the lifecycle of Meshery Operator |
| [Broker](/concepts/architecture/broker) | `system check` | follows the lifecycle of Meshery Operator | No | On click of connection chip in UI to Server to Brokers exposed service port | NATS Topic Subscription | On click of connection chip | follows the lifecycle of Meshery Operator | follows the lifecycle of Meshery Operator | follows the lifecycle of Meshery Operator | follows the lifecycle of Meshery Operator |

<br>

Please also see the [Troubleshooting Toolkit](https://docs.google.com/document/d/1q-aayRqx3QKIk2soTaTTTH-jmHcVXHwNYFsYkFawaME/edit#heading=h.ngupcd4j1pfm) and the [Meshery v0.7.0: Connection States (Kubnernetes) Design Review](https://discuss.meshery.io/t/meshery-v0-7-0-connection-states-kubnernetes-design-review/958)
