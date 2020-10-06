---
layout: page
title: Traefik Mesh Adapter
name: Meshery Adapter for Traefik Mesh
mesh_name: Traefik Mesh
version: v1.0
port: 10006/tcp
project_status: alpha
github_link: https://github.com/layer5io/meshery-maesh
image: /docs/assets/img/service-meshes/traefik.svg
---
# {{ page.name }}

| Service Mesh   | Adapter Status | Latest Supported Mesh Version |
| :------------: | :------------:   | :------------:              |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) | {{page.version}}  |

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of {{page.mesh_name}}. A number of sample applications can be installed using the {{page.name}}.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).

### Sample Applications

The {{ page.name }} includes some sample applications operations. Meshery can be used to deploy any of these sample applications.  

- [Istio BookInfo](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md#what-is-the-bookinfo-application)
    - This application is a polyglot composition of microservices are written in different languages and sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- [httpbin](https://httpbin.org)
    - This is a simple HTTP Request & Response Service.
