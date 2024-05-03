---
layout: page
title: Meshery Adapter for Kuma
name: Meshery Adapter for Kuma
component: Kuma
earliest_version: v1.2.2
port: 10007/gRPC
project_status: stable
lab: kuma-meshery-adapter
github_link: https://github.com/meshery/meshery-kuma
image: /assets/img/service-meshes/kuma.svg
white_image: /assets/img/service-meshes/kuma-white.svg
permalink: extensibility/adapters/kuma
redirect_from: service-meshes/adapters/kuma
language: en
---
{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
      {% if group.name == "meshery-kuma" %}
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

The {{page.name}} can install **{{page.earliest_version}}** of {{page.component}} infrastructure. A number of sample applications can be installed using the {{page.name}}.

### Install {{ page.component }}

Choose the Meshery Adapter for {{ page.component }}.

<a href="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-adapter.png" />
</a>

Click on (+) and choose the {{page.earliest_version}} of the {{page.component}} infrastructure.

<a href="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-install.png" />
</a>

## Workload Management

The following sample applications are available in this adapter.

- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)
  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.
