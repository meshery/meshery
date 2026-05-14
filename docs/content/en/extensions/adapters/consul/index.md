---
title: Meshery Adapter for Consul
name: Meshery Adapter for Consul
component: Consul
earliest_version: v1.8.4
adapter_version: v0.8.3
port: 10002/gRPC
project_status: stable
lab: consul-meshery-adapter
github_link: https://github.com/meshery/meshery-consul
image: /extensions/adapters/consul/images/consul.svg
white_image: /extensions/adapters/consul/images/consul-white.svg
aliases: 
- /service-meshes/adapters/consul
- /extensibility/adapters/consul
---

### Features

1. Lifecycle management of Consul
1. Lifecycle management of sample applications
1. Performance management of Consul and it workloads
   - Prometheus and Grafana integration
1. Configuration management and best practices of Consul
1. Custom configuration

### Sample Applications

Meshery supports the deployment of a variety of sample applications on Meshery Adapter for Consul. Use Meshery to deploy any of these sample applications.

- [httpbin](/guides/infrastructure-management/sample-apps#httpbin)
  - Httpbin is a simple HTTP request and response service.
- [Bookinfo](/guides/infrastructure-management/sample-apps#bookinfo)
  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- [Image Hub](/guides/infrastructure-management/sample-apps#imagehub)
  - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.

### Performance management of Consul and it workloads

#### Prometheus and Grafana integration

The Meshery Adapter for Consul will connect to Meshery Adapter for Consul's Prometheus and Grafana instances running in the control plane (typically found in a separate namespace) or other instances to which Meshery has network reachability.

### Architecture

[![Consul Service Mesh Architecture](/extensions/adapters/consul/images/service-mesh-architecture-consul.png)](/extensions/adapters/consul/images/service-mesh-architecture-consul.png)

### Suggested Topics

- Examine [Meshery's architecture](/concepts/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters](/concepts/architecture/adapters).
