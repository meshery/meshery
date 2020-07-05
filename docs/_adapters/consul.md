---
layout: page
title: Consul Adapter
name: Consul
version: v1.5
port: 10002/tcp
project_status: stable
image: /docs/assets/img/service-meshes/consul.svg
---
# Meshery Adapter for {{ page.name }}

| Adapter Status |
| :------------: |
| [{{ page.project_status }}]({{ page.github_link }})|

## Features

1. Lifecycle management of Consul
1. Lifecycle management of sample applications
1. Performance management of Consul and it workloads
    1. Prometheus and Grafana integration
1. Configuration management and best practices of Consul
1. Custom service mesh configuration

## Sample Applications

The Meshery adapter for {{ page.name }} includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

- [httpbin](https://httpbin.org)
    - Httpbin is a simple HTTP request and response service.
- [Istio BookInfo](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md#what-is-the-bookinfo-application)
    - This application is a polyglot composition of microservices are written in different languages and sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- [Image Hub](https://github.com/layer5io/image-hub)
    - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.
    
![Layer5 Image Hub on HashiCorp Consul](/docs/service-meshes/adapters/consul/layer5-image-hub-on-hashicorp-consul.png)


## Performance management of Consul and it workloads

### Prometheus and Grafana integration

The Meshery adapter for {{ page.name }} will connect to {{ page.name }}'s Prometheus and Grafana instances running in the control plane (typically found in a separate namespace) or other instances to which Meshery has network reachability.

## Architecture

![Consul Service Mesh Archicture](/docs/service-meshes/adapters/consul/service-mesh-architecture-consul.png)