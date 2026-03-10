---
title: Meshery Adapter for NGINX Service Mesh
name: Meshery Adapter for NGINX Service Mesh
component: NGINX Service Mesh
earliest_version: v1.2.0
adapter_version: v0.8.1
port: 10010/gRPC
project_status: stable
github_link: https://github.com/meshery/meshery-nginx-sm
image: /extensions/adapters/nginx-sm/images/nginx-sm.svg
white_image: /extensions/adapters/nginx-sm/images/nginx-sm-white.svg
aliases: 
- /service-meshes/adapters/nginx-sm
- /extensibility/adapters/nginx-sm
---

The Meshery Adapter for NGINX Service Mesh is currently under construction (stable state). Want to contribute? Check our [progress](https://github.com/meshery/meshery-nginx-sm).

### Features

1. Lifecycle management of NGINX Service Mesh
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

The Meshery Adapter for NGINX Service Mesh includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

- [Emojivoto](/guides/infrastructure-management/sample-apps)

  - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo](/guides/infrastructure-management/sample-apps)

- [Httpbin](/guides/infrastructure-management/sample-apps)

  - Httpbin is a simple HTTP request and response service.

- [NGINX Books](https://github.com/BuoyantIO/booksapp)
  - Application that helps you manage your bookshelf.

Identify overhead involved in running NGINX Service Mesh, various NGINX Service Mesh configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration


