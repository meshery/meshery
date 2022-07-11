---
layout: default
title: Meshery Adapter for NGINX Service Mesh
name: Meshery Adapter for NGINX Service Mesh
mesh_name: NGINX Service Mesh
earliest_version: v1.2.0
port: 10010/gRPC
project_status: stable
github_link: https://github.com/meshery/meshery-nginx-sm
image: /assets/img/service-meshes/nginx-sm.svg
permalink: service-meshes/adapters/nginx-sm
---

{% include adapter-status.html %}

The {{ page.name }} is currently under construction ({{ page.project_status }} state). Want to contribute? Check our [progress]({{page.github_link}}).

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of {{page.mesh_name}} service mesh. A number of sample applications can be installed using the {{page.name}}.

### Features

1. Lifecycle management of {{page.mesh_name}}
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

The {{ page.name }} includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

- [Emojivoto]({{site.baseurl}}/guides/sample-apps#emojivoto)

  - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)

  - Follow this [tutorial workshop](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md) to set up and deploy the BookInfo sample app on Istio using Meshery.

- [Httpbin]({{site.baseurl}}/guides/sample-apps#httpbin)

  - Httpbin is a simple HTTP request and response service.

- [NGINX Servce Mesh Books](https://github.com/BuoyantIO/booksapp)
  - Application that helps you manage your bookshelf.

Identify overhead involved in running {{page.mesh_name}}, various {{page.mesh_name}} configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [{{page.name}}]({{ page.github_link }}) will connect to NGINX Service Mesh's Prometheus and Grafana instances running in the control plane.

