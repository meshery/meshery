---
title: Meshery Adapter for App Mesh
name: Meshery Adapter for App Mesh
component: App Mesh
earliest_version: v1.4.1
adapter_version: v0.6.4
port: 10005/gRPC
project_status: beta
github_link: https://github.com/meshery/meshery-app-mesh
image: /extensions/adapters/app-mesh/images/app-mesh.svg   
white_image: /extensions/adapters/app-mesh/images/app-mesh-white.svg
aliases: 
- /service-meshes/adapters/app-mesh
- /extensibility/adapters/app-mesh
---

### Features

1. Lifecycle management of App Mesh
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

- [Emojivoto](/guides/infrastructure-management/sample-apps#emojivoto)

  - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo](/guides/infrastructure-management/sample-apps#bookinfo)

- [Httpbin](/guides/infrastructure-management/sample-apps#httpbin)

  - Httpbin is a simple HTTP request and response service.

### Performance Management

Identify overhead involved in running the App Mesh, various App Mesh configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery Adapter for App Mesh](https://github.com/meshery-extensions/meshery-app-mesh) will connect to App Mesh’s Prometheus and Grafana instances running in the control plane.