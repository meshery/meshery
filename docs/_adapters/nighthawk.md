---
layout: enhanced
title: Meshery Adapter for Nighthawk
name: Meshery Adapter for Nighthawk
component: Nighthawk
earliest_version: v0.7
port: 10013/gRPC
project_status: alpha
#lab: nighthawk-meshery-adapter
github_link: https://github.com/meshery/meshery-nighthawk
image: /assets/img/adapters/nighthawk/nighthawk.svg
white_image: /assets/img/adapters/nighthawk/nighthawk-white.svg
permalink: extensibility/adapters/nighthawk
language: en
---

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
      {% if group.name == "meshery-linkerd" %}
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

### Features

1. Lifecycle management of {{page.component}}
1. Performance characterization

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of {{page.component}}. Performance tests of various configurations can run using the {{page.name}}.

### Install {{ page.component }}

##### Choose the Meshery Adapter for {{ page.component }}

<!-- <a href="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-adapter.png" />
</a> -->

##### Click on (+) and choose the {{page.earliest_version}} of the {{page.component}}.

<!-- <a href="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-install.png" />
</a> -->

## Performance characterization

Identify overhead involved in running {{page.component}}, various {{page.component}} configurations while running different workloads and on different infrastructure. The adapter facilitates data plane and control plane performance testing.

1. Prometheus integration
1. Grafana integration

### Performance tests
