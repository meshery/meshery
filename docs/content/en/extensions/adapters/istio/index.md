---
title: Meshery Adapter for Istio
name: Meshery Adapter for Istio
component: Istio
earliest_version: v1.6.0
adapter_version: v0.8.3
port: 10000/gRPC
project_status: stable
lab: istio-meshery-adapter
github_link: https://github.com/meshery/meshery-istio
image: /extensions/adapters/istio/images/istio.svg
white_image: /extensions/adapters/istio/images/istio-white.svg
aliases: 
- /service-meshes/adapters/istio
- /extensibility/adapters/istio
---

## Features

1. Istio Lifecycle Management
1. Workload Lifecycle Management
1. Cloud Native Performance (SMP)
   1. Prometheus and Grafana connections
1. Configuration Analysis, Patterns, and Best Practices
   1. Custom Configuration

### Workload Management

The Meshery Adapter for Istio includes a handful of sample applications. Use Meshery to deploy any of these sample applications:

- [Bookinfo](/guides/infrastructure-management/sample-apps#bookinfo)
- [Httpbin](/guides/infrastructure-management/sample-apps#httpbin)
  - Httpbin is a simple HTTP request and response service.
- [Online Boutique](/guides/infrastructure-management/sample-apps#online-boutique)
  - Online Boutique Application is a web-based, e-commerce demo application from the Google Cloud Platform.
- [Image Hub](/guides/infrastructure-management/sample-apps#imagehub)
  - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.

## Using Cloud Native Standards

Meshery's powerful performance management functionality is accomplished through implementation of [Cloud Native Performance](https://smp-spec.io). Meshery enables operators to deploy WebAssembly filters to Envoy-based data planes. Meshery facilitates learning about functionality and performance of infrastructure and workloads and incorporates the collection and display of metrics from applications using Prometheus and Grafana integrations.

### Design Patterns and Meshery Models

### Prometheus and Grafana connections

The Meshery Adapter for Istio allows you to quickly deploy (or remove) an Istio add-ons. Meshery will deploy the Prometheus and Grafana add-ons (including Jaeger and Kiali) into Istio's control plane (typically the `istio-system` namespace). You can also connect Meshery to Prometheus, Grafana instances not running in the control plane.

If you already have existing Prometheus or Grafana deployments in your cluster, MeshSync will discover them and attempt to automatically register them for use.

## Configuration Management

Meshery Adapter for Istio provides

### Configuration best practices

On demand, the Meshery Adapter for Istio will parse all of Istio's configuration and compare the running configuration of the infrastructure against known best practices for an Meshery Adapter for Istio deployment.

### Custom infrastructure configuration

Meshery allows you to apply configuration to your infrastructure deployment. You can paste (or type in) any Kubernetes manifest that you would like to have applied to your infrastructure, in fact, you can apply any configuration that you would like to your Kubernetes cluster. This configuration may be VirtualServices, DestinationRules or any other custom Istio resource.

<a href="/extensions/adapters/istio/images/istio-adapter-custom-configuration.png">
  <img style="width:500px;" src="/extensions/adapters/istio/images/istio-adapter-custom-configuration.png" />
</a>

Add-on resources can be applied **or** deleted using this custom configuration operation.

### Suggested Topics

- Examine [Meshery's architecture](/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters](/architecture/adapters).

