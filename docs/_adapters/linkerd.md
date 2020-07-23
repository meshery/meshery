---
layout: page
title: Linkerd Adapter
name: Linkerd
version: v2.5.0
port: 10001/tcp
project_status: stable
image: /docs/assets/img/service-meshes/linkerd.svg
---
# Meshery Adapter for {{ page.name }}

| Adapter Status |
| :------------: |
| [{{ page.project_status }}]({{ page.github_link }})|

## Sample Applications

The Meshery adapter for ({{ page.title }}) includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

## Features
1. Lifecycle management of Linkerd
1. Lifecycle management of sample applications
1. Performance testing

Identify overhead involved in running Linkerd, various Linkerd configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [Meshery adapter for Linkerd]({{ page.github_link }}) will connect to Linkerd's Prometheus and Grafana instances running in the control plane.
