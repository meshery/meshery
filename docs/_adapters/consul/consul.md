---
layout: default
title: Meshery Adapter for Consul
name: Meshery Adapter for Consul
mesh_name: Consul
earliest_version: v1.8.4
port: 10002/gRPC
project_status: stable
lab: consul-meshery-adapter
github_link: https://github.com/meshery/meshery-consul
image: /assets/img/service-meshes/consul.svg
white_image: /assets/img/service-meshes/consul.svg
language: en
language: en
#redirect_from: service-meshes/adapters/consul
permalink: service-meshes/adapters/consul
---

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
      {% if group.name == "meshery-consul" %}
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

{% include adapter-labs.html %}

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of the {{page.mesh_name}} service mesh.

### Install {{ page.mesh_name }}

##### Choose the Meshery Adapter for {{page.mesh_name}}

<a href="{{ site.baseurl }}/assets/img/adapters/consul/consul-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/consul/consul-adapter.png" />
</a>

##### Click on (+) and choose the {{page.earliest_version}} of the {{page.mesh_name}} service mesh.

<a href="{{ site.baseurl }}/assets/img/adapters/consul/consul-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/consul/consul-install.png" />
</a>

A number of [sample applications](#sample-applications) for {{page.mesh_name}} can also be installed using Meshery.

### Features

1. Lifecycle management of {{page.mesh_name}}
1. Lifecycle management of sample applications
1. Performance management of {{page.mesh_name}} and it workloads
   - Prometheus and Grafana integration
1. Configuration management and best practices of {{page.mesh_name}}
1. Custom service mesh configuration

### Sample Applications

Meshery supports the deployment of a variety of sample applications on {{ page.name }}. Use Meshery to deploy any of these sample applications.

- [httpbin]({{site.baseurl}}/guides/sample-apps#httpbin)
  - Httpbin is a simple HTTP request and response service.
- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)
  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
- [Image Hub]({{site.baseurl}}/guides/sample-apps#imagehub)
  - Image Hub is a sample application written to run on Consul for exploring WebAssembly modules used as Envoy filters.

[![Layer5 Image Hub on HashiCorp Consul]({{ site.baseurl }}/service-meshes/adapters/consul/layer5-image-hub-on-hashicorp-consul.png)]({{ site.baseurl }}/service-meshes/adapters/consul/layer5-image-hub-on-hashicorp-consul.png)

### Performance management of Consul and it workloads

#### Prometheus and Grafana integration

The {{ page.name }} will connect to {{ page.name }}'s Prometheus and Grafana instances running in the control plane (typically found in a separate namespace) or other instances to which Meshery has network reachability.

### Architecture

[![Consul Service Mesh Architecture]({{ site.baseurl }}/service-meshes/adapters/consul/service-mesh-architecture-consul.png)]({{ site.baseurl }}/service-meshes/adapters/consul/service-mesh-architecture-consul.png)

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
