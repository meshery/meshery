---
layout: default
title: Meshery Adapter for Istio
name: Meshery Adapter for Istio
mesh_name: Istio
earliest_version: v1.6.0
port: 10000/gRPC
project_status: stable
lab: istio-meshery-adapter
github_link: https://github.com/meshery/meshery-istio
image: /assets/img/service-meshes/istio.svg
permalink: service-meshes/adapters/istio
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

{% include adapter-status.html %}

{% include adapter-labs.html %}

## Features

1. {{page.mesh_name}} Lifecycle Management
1. Workload Lifecycle Management
   1. Using Service Mesh Standards
      1. Service Mesh Performance (SMP)
         1. Prometheus and Grafana connections
      1. Service Mesh Interface (SMI)
1. Configuration Analysis, Patterns, and Best Practices
   1. Custom Service Mesh Configuration

### Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of the {{page.mesh_name}} service mesh. The SMI adapter for Istio can also be installed using Meshery.

### Install {{ page.mesh_name }}

In Meshery's UI, choose the Meshery Adapter for {{ page.mesh_name }}.

<a href="{{ site.baseurl }}/assets/img/adapters/istio/istio-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/istio/istio-adapter.png" />
</a>

Click on (+) and choose the {{page.earliest_version}} of the {{page.mesh_name}} service mesh.

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

## Using Service Mesh Standards

As the open source, cloud native management plane, Meshery enables the adoption, operation, and management of Kubernetes, any service mesh, and their workloads. Meshery's powerful performance management functionality is accomplished through implementation of [Service Mesh Performance](https://smp-spec.io) (SMP). Meshery's cloud native manager functionality leverages [Service Mesh Interface](https://smi-spec.io) (SMI) and Meshery is the conformance tool for SMI. Meshery integrates with Open Application Model (OAM) to enable users to deploy service mesh patterns. Meshery enables operators to deploy WebAssembly filters to Envoy-based data planes. Meshery facilitates learning about functionality and performance of service meshes and incorporates the collection and display of metrics from applications using Prometheus and Grafana integrations.

### Service Mesh Patterns and Open Application Model (OAM)

### Complying with Service Mesh Interface (SMI)

Meshery allows you to analyze the compliance status and functional capabilities of your service mesh. This allows you to compare high-level functional differences between service meshes and verify whether your service mesh is conformant with the SMI specification.

Learn more about the SMI specification and [Meshery's conformance test suite]({{ site.baseurl }}/functionality/service-mesh-interface).

### Managing Service Mesh Performance (SMP)

### Prometheus and Grafana connections

The {{page.name}} allows you to quickly deploy (or remove) an Istio add-ons. Meshery will deploy the Prometheus and Grafana add-ons (including Jaeger and Kiali) into Istio's control plane (typically the `istio-system` namespace). You can also connect Meshery to Prometheus, Grafana instances not running in the service mesh's control plane.

If you already have existing Prometheus or Grafana deployments in your cluster, MeshSync will discover them and attempt to automatically register them for use.

## Configuration Management

{{page.name}} provides

### Configuration best practices

On demand, the {{page.name}} will parse all of Istio's configuration and compare the running configuration of the service mesh against known best practices for an {{page.title}} deployment.

### Custom service mesh configuration

Meshery allows you to apply configuration to your service mesh deployment. You can paste (or type in) any Kubernetes manifest that you would like to have applied to your service mesh, in fact, you can apply any configuration that you would like to your Kubernetes cluster. This configuration may be VirtualServices, DestinationRules or any other custom Istio resource.

<a href="{{ site.baseurl }}istio-adapter-custom-configuration.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/istio/istio-adapter-custom-configuration.png" />
</a>

Service mesh resources can be applied **or** deleted using this custom configuration operation.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
