---
layout: page
title: Citrix Service Mesh (CPX) Adapter
name: Citrix Service Mesh
version: "1.0"
port: 10008/tcp
project_status: beta
github_link: https://github.com/layer5io/meshery-cpx
image: /docs/assets/img/service-meshes/cpx.png
---

# {{ page.name }}

|  Service Mesh  |                   Adapter Status                    | Latest Supported Mesh Version |
| :------------: | :-------------------------------------------------: | :---------------------------: |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) |       {{page.version}}        |

### Features

1. Lifecycle Management of Citrix
2. Lifecycle Management of Sample Applications

### Sample Applications

The Meshery adapter for {{ page.name }} includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

- [hipster](https://github.com/GoogleCloudPlatform/microservices-demo)

  - Hipster Shop Application is a web-based, e-commerce demo application from the Google Cloud Platform.

- [httpbin](https://httpbin.org)
  - Httpbin is a simple HTTP request and response service.
- [Istio BookInfo Application](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md#what-is-the-bookinfo-application)
  - This application is a polyglot composition of microservices that are written in different languages. The application displays information about a book, similar to a single catalog entry of an online book store. Displayed on the page is a description of the book, book details (ISBN, number of pages, and so on), and a few book reviews.

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the {{page.name}} service mesh. The SMI adapter for Kuma can also be installed using Meshery.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
