---
layout: default
title: Meshery Adapter for App Mesh
name: Meshery Adapter for App Mesh
component: App Mesh
earliest_version: v1.4.1
port: 10005/gRPC
project_status: beta
github_link: https://github.com/meshery/meshery-app-mesh
image: /assets/img/service-meshes/app-mesh.svg
white_image: /assets/img/service-meshes/app-mesh-white.svg
permalink: extensibility/adapters/app-mesh
redirect_from: service-meshes/adapters/app-mesh
language: en
---

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
      {% if group.name == "meshery-app-mesh" %}
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

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of {{page.component}}. A number of sample applications can be installed using the {{page.name}}.

### Features

1. Lifecycle management of {{page.component}}
1. Lifecycle management of sample applications
1. Performance testing

### Sample Applications

The {{ page.name }} includes a handful of sample applications. Use Meshery to deploy any of these sample applications.

- [Emojivoto]({{site.baseurl}}/guides/sample-apps#emojivoto)

  - A microservice application that allows users to vote for their favorite emoji, and tracks votes received on a leaderboard.

- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)

  - Follow this [tutorial workshop](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md) to set up and deploy the BookInfo sample app on Istio using Meshery.

- [Httpbin]({{site.baseurl}}/guides/sample-apps#httpbin)

  - Httpbin is a simple HTTP request and response service.

Identify overhead involved in running {{page.component}}, various {{page.component}} configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

The [{{page.name}}]({{ page.github_link }}) will connect to {{page.component}}'s Prometheus and Grafana instances running in the control plane.
