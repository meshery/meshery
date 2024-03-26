---
layout: enhanced
title: Meshery Adapter for Istio
name: Meshery Adapter for Istio
component: Istio
earliest_version: v1.6.0
port: 10000/gRPC
project_status: stable
lab: istio-meshery-adapter
github_link: https://github.com/meshery/meshery-istio
image: /assets/img/service-meshes/istio.svg
white_image: /assets/img/service-meshes/istio-white.svg
permalink: extensibility/adapters/istio
redirect_from: service-meshes/adapters/istio
language: en
---

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
{% if group.name == "meshery-istio" %}
{% assign items = group.items | sort: "meshery-component-version" | reverse %}
{% for item in items %}
{% if item.meshery-component-version != "edge" %}
{% if item.overall-status == "passing" %}
{% assign adapter_version_dynamic = item.meshery-component-version %}
{% break %}
{% elsif item.overall-status == "failing" %}
{% continue %}
{% endif %}
{% endif %}
{% endfor %}
{% endif %}
{% endfor %}

{% include compatibility/adapter-status.html %}

<!-- {% include adapter-labs.html %} -->

## Features

1. {{page.component}} Lifecycle Management
1. Workload Lifecycle Management
1. Cloud Native Performance (SMP)
   1. Prometheus and Grafana connections
1. Configuration Analysis, Patterns, and Best Practices
   1. Custom Configuration

### Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of the {{page.component}}.

### Install {{ page.component }}

In Meshery's UI, choose the Meshery Adapter for {{ page.component }}.

<a href="{{ site.baseurl }}/assets/img/adapters/istio/istio-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/istio/istio-adapter.png" />
</a>

Click on (+) and choose the {{page.earliest_version}} of the {{page.component}}.

<a href="{{ site.baseurl }}/assets/img/adapters/istio/istio-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/istio/istio-install.png" />
</a>

### Workload Management

The ({{page.name}}) includes a handful of sample applications. Use Meshery to deploy any of these sample applications:

- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)
  - Follow this [tutorial workshop](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md) to set up and deploy the BookInfo sample app on Istio using Meshery.
- [Httpbin]({{site.baseurl}}/guides/sample-apps#httpbin)
  - Httpbin is a simple HTTP request and response service.
- [Online Boutique]({{site.baseurl}}/guides/sample-apps#online-boutique)
  - Online Boutique Application is a web-based, e-commerce demo application from the Google Cloud Platform.
- [Image Hub]({{site.baseurl}}/guides/sample-apps#imagehub)
  - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.

## Using Cloud Native Standards

Meshery's powerful performance management functionality is accomplished through implementation of [Cloud Native Performance](https://smp-spec.io). Meshery enables operators to deploy WebAssembly filters to Envoy-based data planes. Meshery facilitates learning about functionality and performance of infrastructure and workloads and incorporates the collection and display of metrics from applications using Prometheus and Grafana integrations.

### Design Patterns and Meshery Models

### Prometheus and Grafana connections

The {{page.name}} allows you to quickly deploy (or remove) an Istio add-ons. Meshery will deploy the Prometheus and Grafana add-ons (including Jaeger and Kiali) into Istio's control plane (typically the `istio-system` namespace). You can also connect Meshery to Prometheus, Grafana instances not running in the control plane.

If you already have existing Prometheus or Grafana deployments in your cluster, MeshSync will discover them and attempt to automatically register them for use.

## Configuration Management

{{page.name}} provides

### Configuration best practices

On demand, the {{page.name}} will parse all of Istio's configuration and compare the running configuration of the infrastructure against known best practices for an {{page.title}} deployment.

### Custom infrastructure configuration

Meshery allows you to apply configuration to your infrastructure deployment. You can paste (or type in) any Kubernetes manifest that you would like to have applied to your infrastructure, in fact, you can apply any configuration that you would like to your Kubernetes cluster. This configuration may be VirtualServices, DestinationRules or any other custom Istio resource.

<a href="{{ site.baseurl }}istio-adapter-custom-configuration.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/istio/istio-adapter-custom-configuration.png" />
</a>

Add-on resources can be applied **or** deleted using this custom configuration operation.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).

