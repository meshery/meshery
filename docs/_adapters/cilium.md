---
layout: enhanced
title: Meshery Adapter for Cilium Service Mesh
name: Meshery Adapter for Cilium Service Mesh
component: Cilium
earliest_version: v1.10.6
port: 10012/gRPC
project_status: stable
lab: cilium-meshery-adapter
github_link: https://github.com/meshery/meshery-cilium
image: /assets/img/service-meshes/cilium.svg
white_image: /assets/img/service-meshes/cilium-white.svg
permalink: extensibility/adapters/cilium
redirect_from: service-meshes/adapters/cilium
language: en
---

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
      {% if group.name == "meshery-cilium" %}
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

The {{ page.name }} is currently under construction ({{ page.project_status }} state), which means that the adapter is not functional and cannot be interacted with through the <a href="{{ site.baseurl }}/installation#6-you-will-now-be-directed-to-the-meshery-ui"> Meshery UI </a>at the moment. Check back here to see updates.

Want to contribute? Check our [progress]({{page.github_link}}).

### Install {{ page.component }}

##### Choose the Meshery Adapter for {{ page.component }}

<a href="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-adapter.png" />
</a>

##### Click on (+) and choose the {{page.earliest_version}} of the {{page.component}}.

<a href="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/linkerd/linkerd-install.png" />
</a>

