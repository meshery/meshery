---
layout: page
title: Linkerd Adapter
name: Meshery Adapter for Linkerd
version: v2.5.0
port: 10001/tcp
project_status: stable
github_link: https://github.com/layer5io/meshery-linkerd
image: /docs/assets/img/service-meshes/linkerd.svg
---
# {{ page.name }}

| Service Mesh   | Adapter Status | Latest Supported Mesh Version |
| :------------: | :------------:   | :------------:              |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) | {{page.version}}  |

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the Linkerd service mesh. The SMI adapter for Linkerd can also be installed using Meshery.

### Features
1. Lifecycle management of Linkerd
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

The {{ page.name }} includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

Identify overhead involved in running Linkerd, various Linkerd configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery adapter for Linkerd]({{ page.github_link }}) will connect to Linkerd's Prometheus and Grafana instances running in the control plane.
