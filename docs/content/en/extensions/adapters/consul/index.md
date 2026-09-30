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
image: extensions/adapters/consul/images/consul.svg
white_image: extensions/adapters/consul/images/consul-white.svg
aliases: 
- /service-meshes/adapters/consul
- /extensibility/adapters/consul
---

{{< alert title="How adapters install and configure infrastructure" >}}
Adapters share a common Day 1 / Day 2 model: they install infrastructure with the project's own Helm charts, Kubernetes manifests, and/or CLI (with fallback), and they apply ongoing configuration by deploying Meshery Designs through Meshery Server's registry of registrants. That behavior, and related FAQs, is documented on the [Adapters]({{< ref "concepts/architecture/adapters.md" >}}) page. This page covers only what is specific to Consul.
{{< /alert >}}

### Features

1. Lifecycle management of Consul
1. Lifecycle management of sample applications
1. Performance management of Consul and its workloads
   - Prometheus and Grafana integration
1. Configuration management and best practices of Consul
1. Custom configuration

### Lifecycle management of Consul

The Meshery Adapter for Consul can install **v1.8.4** and later of Consul on a connected Kubernetes cluster.

This adapter installs the Consul control plane with the official HashiCorp Consul Helm chart (`hashicorp/consul`). Sample applications and custom configuration are applied as Kubernetes manifests. A Helm-rendered Consul demo snapshot is also available as manifests under the adapter's config templates.

To install Consul from Meshery's UI, choose the Meshery Adapter for Consul, click **(+)**, and select a supported Consul version.

See [Adapters]({{< ref "concepts/architecture/adapters.md" >}}) for the shared installation methods, fallback order, and how subsequent configuration is applied through Designs.

### Sample Applications

Meshery supports the deployment of a variety of sample applications on Meshery Adapter for Consul. Use Meshery to deploy any of these sample applications.

- httpbin
  - Httpbin is a simple HTTP request and response service.
- Bookinfo
  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- Image Hub
  - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.

### Performance management of Consul and its workloads

#### Prometheus and Grafana integration

The Meshery Adapter for Consul will connect to Prometheus and Grafana instances running in the control plane (typically found in a separate namespace) or other instances to which Meshery has network reachability.

### Custom configuration

You can paste (or type in) any Kubernetes manifest and have this adapter apply or delete it in the requested namespace. Use this for Consul-related resources, sample-app changes, or other cluster configuration Meshery can reach.

### Architecture

[![Consul Service Mesh Architecture](images/service-mesh-architecture-consul.png)](images/service-mesh-architecture-consul.png)

### Suggested Topics

- Examine [Meshery's architecture]({{< ref "concepts/architecture/_index.md" >}}) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{< ref "concepts/architecture/adapters.md" >}}), including [Day 1 installation]({{< ref "concepts/architecture/adapters.md" >}}#day-1-installing-infrastructure), [Day 2 configuration through Designs]({{< ref "concepts/architecture/adapters.md" >}}#day-2-ongoing-configuration-through-designs), and [adapter FAQs]({{< ref "concepts/architecture/adapters.md" >}}#adapter-faqs).
