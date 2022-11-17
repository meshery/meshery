---
layout: default
title: Citrix Service Mesh (CPX)
name: Meshery Adapter for Citrix Service Mesh
mesh_name: Citrix
earliest_version: "1.0"
port: 10008/gRPC
project_status: beta
github_link: https://github.com/meshery/meshery-cpx
image: /assets/img/service-meshes/citrix.svg
---

{% include adapter-status.html %}

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of the {{page.mesh_name}} service mesh.

### Install {{ page.mesh_name }}

##### Choose the Meshery Adapter for Citrix

<a href="{{ site.baseurl }}/assets/img/adapters/citrix/citrix-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/citrix/citrix-adapter.png">
</a>

##### Click on (+) and choose the {{page.earliest_version}} of the {{page.mesh_name}} service mesh.

<a href="{{ site.baseurl }}/assets/img/adapters/citrix/citrix-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/citrix/citrix-install.png">
</a>

## Features

1. Lifecycle Management of {{page.mesh_name}}
2. Lifecycle Management of Sample Applications

### Sample Applications

The {{page.name}} includes a handful of sample applications. Some of these applications are from other service meshes and some of these sample applications are general-purpose examples. Use Meshery to deploy any of these sample applications.

- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)
  - Follow this [tutorial workshop](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md) to set up and deploy the BookInfo sample app on Istio using Meshery.
- [Httpbin]({{site.baseurl}}/guides/sample-apps#httpbin)
  - Httpbin is a simple HTTP request and response service.
- [Online Boutique]({{site.baseurl}}/guides/sample-apps#online-boutique)
  - Online Boutique Application is a web-based, e-commerce demo application from the Google Cloud Platform.

### Suggested Topics

- Examine [Meshery's architecture]({{site.baseurl}}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{site.baseurl}}/architecture/adapters).
