---
layout: page
title: Linkerd Adapter
name: Meshery Adapter for Linkerd
mesh_name: Linkerd
version: v2.5.0
port: 10001/tcp
project_status: stable
github_link: https://github.com/layer5io/meshery-linkerd
image: /docs/assets/img/service-meshes/linkerd.svg
---
# {{ page.name }}

| Service Mesh   | Adapter Status | Latest Supported Mesh Version |
| :------------: | :------------:   | :------------:              |
| <img src="{{ page.image }}" style="width:20px" /> {{ page.mesh_name }} | [{{ page.project_status }}]({{ page.github_link }}) | {{page.version}}  |

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of {{page.mesh_name}}. A number of sample applications can be installed using the {{page.name}}.

### Install {{ page.mesh_name }}

##### **Choose the Meshery adapter for {{ page.mesh_name }}**

<a href="#linkerd-adapter">
  <img style="width:500px;" src="/docs/assets/img/adapters/linkerd/linkerd-adapter.png" />
</a>
<a href="#" class="lightbox" id="linkerd-adapter">
  <span style="background-image: url('/docs/assets/img/adapters/linkerd/linkerd-adapter.png')"></span>
</a>

##### **Click on (+) and choose the `{{page.version}}` of the {{page.mesh_name}} service mesh.**

<a href="#linkerd-install">
  <img style="width:500px;" src="/docs/assets/img/adapters/linkerd/linkerd-install.png" />
</a>
<a href="#" class="lightbox" id="linkerd-install">
  <span style="background-image: url('/docs/assets/img/adapters/linkerd/linkerd-install.png')"></span>
</a>

### Features
1. Lifecycle management of {{page.mesh_name}}
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

The {{ page.name }} includes the ability to deploy a variety of sample applications. Use Meshery to deploy any of these sample applications:

- [Emojivoto](/docs/guides/sample-apps#emoji.voto)
    - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo](/docs/guides/sample-apps#bookinfo) 
    - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.

- [Linkerd Books](/docs/guides/sample-apps#linkerdbooks)
    - A sample application built for demonstrating  manage your bookshelf.

- [HTTPbin](/docs/guides/sample-apps#httpbin)
    - A simple HTTP Request & Response Service.

Identify overhead involved in running {{page.mesh_name}}, various {{page.mesh_name}} configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery adapter for Linkerd]({{ page.github_link }}) will connect to Linkerd's Prometheus and Grafana instances running in the control plane.
