---
layout: page
title: Citrix Service Mesh (CPX) Adapter
name: Meshery Adapter for Citrix Service Mesh
mesh_name: Citrix Service Mesh
version: "1.0"
port: 10008/tcp
project_status: beta
github_link: https://github.com/layer5io/meshery-cpx
image: /docs/assets/img/service-meshes/citrix.svg
---

# {{ page.name }}

|  Service Mesh  |                   Adapter Status                    | Latest Supported Mesh Version |
| :------------: | :-------------------------------------------------: | :---------------------------: |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) |       {{page.version}}        |

## Features

1. Lifecycle Management of {{page.mesh_name}}
2. Lifecycle Management of Sample Applications

### Lifecycle management
The {{page.name}} can install **{{page.version}}** of the {{page.mesh_name}}. 

### Sample Applications

The {{ page.name }} includes a handful of sample applications. Some of these applications are from other service meshes and some of these sample applications are general-purpose examples. Use Meshery to deploy any of these sample applications.

- [hipster](https://github.com/GoogleCloudPlatform/microservices-demo)

  - Hipster Shop Application is a web-based, e-commerce demo application from the Google Cloud Platform.

- [httpbin](https://httpbin.org)
  - Httpbin is a simple HTTP request and response service.

- [Istio BookInfo Application](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md#what-is-the-bookinfo-application)
  - This application is a polyglot composition of microservices that are written in different languages. The application displays information about a book, similar to a single catalog entry of an online book store. Displayed on the page is a description of the book, book details (ISBN, number of pages, and so on), and a few book reviews.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
