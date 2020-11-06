---
layout: page
title: Consul Adapter
name: Meshery Adapter for Consul
mesh_name: Consul
version: v1.8.2
port: 10002/tcp
project_status: stable
github_link: https://github.com/layer5io/meshery-consul
image: /docs/assets/img/service-meshes/consul.svg
---
# {{ page.name }}

| Service Mesh   | Adapter Status | Latest Supported Mesh Version |
| :------------: | :------------:   | :------------:              |
| <img src="{{ page.image }}" style="width:20px" /> {{ page.mesh_name }} | [{{ page.project_status }}]({{ page.github_link }}) | {{page.version}}  |


### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the {{page.mesh_name}} service mesh. 

### Install {{ page.mesh_name }}

##### **Choose the Meshery Adapter for {{page.mesh_name}}**

<a href="#consul-adapter">
  <img style="width:500px;" src="/docs/assets/img/adapters/consul/consul-adapter.png" />
</a>
<a href="#" class="lightbox" id="consul-adapter">
  <span style="background-image: url('/docs/assets/img/adapters/consul/consul-adapter.png')"></span>
</a>

##### **Click on (+) and choose the `{{page.version}}` of the {{page.mesh_name}} service mesh.**

<a href="#consul-install">
  <img style="width:500px;" src="/docs/assets/img/adapters/consul/consul-install.png" />
</a>
<a href="#" class="lightbox" id="consul-install">
  <span style="background-image: url('/docs/assets/img/adapters/consul/consul-install.png')"></span>
</a>

A number of [sample applications](#sample-applications) for {{page.mesh_name}} can also be installed using Meshery.

### Features

1. Lifecycle management of {{page.mesh_name}}
1. Lifecycle management of sample applications
1. Performance management of {{page.mesh_name}} and it workloads
    - Prometheus and Grafana integration
1. Configuration management and best practices of {{page.mesh_name}}
1. Custom service mesh configuration

### Sample Applications

Meshery supports the deployment of a variety of sample applications on {{ page.name }}. Use Meshery to deploy any of these sample applications.

- [httpbin](/docs/guides/sample-apps#httpbin)
    - Httpbin is a simple HTTP request and response service.
- [Bookinfo](/docs/guides/sample-apps#bookinfo) 
    - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- [Image Hub](/docs/guides/sample-apps/imagehub)
    - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.

![Layer5 Image Hub on HashiCorp Consul](/docs/service-meshes/adapters/consul/layer5-image-hub-on-hashicorp-consul.png)

### Performance management of Consul and it workloads

#### Prometheus and Grafana integration

The {{ page.name }} will connect to {{ page.name }}'s Prometheus and Grafana instances running in the control plane (typically found in a separate namespace) or other instances to which Meshery has network reachability.

### Architecture

![Consul Service Mesh Archicture](/docs/service-meshes/adapters/consul/service-mesh-architecture-consul.png)

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
