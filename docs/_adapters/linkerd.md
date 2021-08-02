---
layout: default
title: Meshery Adapter for Linkerd
name: Meshery Adapter for Linkerd
mesh_name: Linkerd
version: v2.5.0
port: 10001/tcp
project_status: stable
adapter_version: v0.5.2
lab: linkerd-meshery-adapter
github_link: https://github.com/meshery/meshery-linkerd
image: /assets/img/service-meshes/linkerd.svg
permalink: service-meshes/adapters/linkerd
---
{% include adapter-status.html %}

{% include adapter-labs.html %}

### Features
1. Lifecycle management of {{page.mesh_name}}
1. Lifecycle management of sample applications
1. Performance testing


## Lifecycle management

The {{page.name}} can install **{{page.version}}** of {{page.mesh_name}}. A number of sample applications can be installed using the {{page.name}}.

### Install {{ page.mesh_name }}

Note: Linkerd's control plane will be deployed to the `linkerd` namespace. Linkerd does not support deployments of its control plane into namespaces under a different name.
##### Choose the Meshery Adapter for {{ page.mesh_name }}

<a href="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-adapter.png" />
</a>


##### Click on (+) and choose the {{page.version}} of the {{page.mesh_name}} service mesh.

<a href="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-install.png" />
</a>

### Sample Applications

The {{ page.name }} includes the ability to deploy a variety of sample applications. Use Meshery to deploy any of these sample applications:

- [Emojivoto]({{ site.baseurl }}/guides/deploying-sample-apps#emoji.voto)
    - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo]({{ site.baseurl }}/guides/deploying-sample-apps#bookinfo) 
    - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.

- [Linkerd Books]({{ site.baseurl }}/guides/deploying-sample-apps#linkerdbooks)
    - A sample application built for demonstrating  manage your bookshelf.

- [HTTPbin]({{ site.baseurl }}/guides/deploying-sample-apps#httpbin)
    - A simple HTTP Request & Response Service.

Identify overhead involved in running {{page.mesh_name}}, various {{page.mesh_name}} configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery Adapter for Linkerd]({{ page.github_link }}) will connect to Linkerd's Prometheus and Grafana instances running in the control plane.
