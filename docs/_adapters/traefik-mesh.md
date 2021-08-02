---
layout: default
title: Meshery Adapter for Traefik Mesh
name: Meshery Adapter for Traefik Mesh
mesh_name: Traefik Mesh
version: v1.0
port: 10006/tcp
project_status: stable
adapter_version: v0.5.2
github_link: https://github.com/meshery/meshery-traefik-mesh
image: /assets/img/service-meshes/traefik-mesh.svg
permalink: service-meshes/adapters/traefik-mesh
---

{% include adapter-status.html %}

## Lifecycle management

The {{page.name}} can install **{{page.version}}** of {{page.mesh_name}}. A number of sample applications can be installed using the {{page.name}}.

The {{ page.name }} is currently under construction ({{ page.project_status }} state), which means that the adapter is not functional and cannot be interacted with through the <a href="{{ site.baseurl }}installation#6-you-will-now-be-directed-to-the-meshery-ui"> Meshery UI </a>at the moment. Check back here to see updates.

Want to contribute? Check our [progress](page.github_link).

### Sample Applications

The {{ page.name }} includes some sample applications operations. Meshery can be used to deploy any of these sample applications.

- [Bookinfo]({{ site.baseurl }}/guides/deploying-sample-apps#bookinfo)
  - Follow this [tutorial workshop](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md) to set up and deploy the BookInfo sample app on Istio using Meshery.
- [Httpbin]({{ site.baseurl }}/guides/deploying-sample-apps#httpbin)
  - Httpbin is a simple HTTP request and response service.

## Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
