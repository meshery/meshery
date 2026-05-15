---
title: Meshery Adapter for Linkerd
name: Meshery Adapter for Linkerd
component: Linkerd
earliest_version: v2.10.2
adapter_version: v0.8.2
port: 10001/gRPC
project_status: stable
lab: linkerd-meshery-adapter
github_link: https://github.com/meshery/meshery-linkerd
image: /extensions/adapters/linkerd/images/linkerd.svg
white_image: /extensions/adapters/linkerd/images/linkerd-white.svg
aliases: 
- /service-meshes/adapters/linkerd
- /extensibility/adapters/linkerd
---

### Features

1. Lifecycle management of Linkerd
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

The Meshery Adapter for Linkerd includes the ability to deploy a variety of sample applications. Use Meshery to deploy any of these sample applications:

- [Emojivoto](/guides/infrastructure-management/sample-apps)

  - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo](/guides/infrastructure-management/sample-apps)

  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.

- [Linkerd Books](/guides/infrastructure-management/sample-apps)

  - A sample application built for demonstrating manage your bookshelf.

- [HTTPbin](/guides/infrastructure-management/sample-apps)
  - A simple HTTP Request & Response Service.

Identify overhead involved in running Linkerd, various Linkerd configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery Adapter for Linkerd](https://github.com/meshery/meshery-linkerd) will connect to Linkerd's Prometheus and Grafana instances running in the control plane.
