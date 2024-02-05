---
layout: default
title: Meshery Overview
permalink: project/overview
# redirect_from: project/overview/
language: en
display-title: true
type: project
category: none
list: exclude
published: true
abstract: Meshery is the self-service engineering platform, enabling collaborative design and operation of cloud and cloud native infrastructure.
---
Meshery is an extensible engineering platform for the collaborative design and operation of cloud and cloud native infrastructure and applications.

Kubernetes-centric. Kubernetes not required.

## Meshery is for all cloud and cloud native infrastructure

Infrastructure diversity is a reality for any enterprise. Whether you're running a single Kubernetes cluster or multiple Kubernetes clusters, on one cloud or multiple clouds, you'll find that Meshery supports your infrastructure diversity (or lack thereof).

## Meshery's Functionality

Meshery supports all Kubernetes-based infrastructure including most cloud services of AWS and GCP platforms. Meshery features can be categorized by:

1. Performance Management
   - Workload and performance characterization with both built-in and external load generators
   - Prometheus and Grafana integration
1. Lifecycle Management (Day 0, Day 1)
   - Cloud and cloud native provisioning
   - Discovery and onboarding of existing environments and workloads
1. Configuration Management (Day 2)
   - Cloud native patterns catalog
   - Configuration best practices
   - Policy engine for relationship inference and context-aware design
1. Collaboration
   - Multi-player infrastructure design and operation
1. Data Plane Intelligence
   - Registry and configuration of WebAssembly filters for Envoy
1. Interoperability and Federation
   - Integration with thousands of cloud services and cloud native projects
   - Manage multiple service meshes concurrently
   - Connect to multiple clusters independently

### Meshery is for Developers, Operators, and Product Owners

Whether making a Day 0 adoption choice or maintaining a Day 2 deployment, Meshery has useful capabilities in either circumstance. Targeted audience for Meshery project would be any technology operators that leverage service mesh in their ecosystem; this includes developers, devops engineers, decision makers, architects, and organizations that rely on microservices platform.

### Meshery is for cloud native patterns

Through [Models]({{site.baseurl}}/concepts/logical/models), Meshery describes infrastructure under management, enabling you to define cloud native designs and patterns and then to export those designs and share within the <a href="https://meshery.io/catalog" target="_self_">Meshery Catalog</a>.

### Meshery is for performance management

Meshery helps users weigh the value of their cloud native deployments against the overhead incurred in running different deployment scenarios and different configruations. Meshery provides statistical analysis of the request latency and throughput seen across various permutations of your workload, infrastructure and infrastructure configuration. In addition to request latency and throughput, Meshery also tracks memory and CPU overhead in of the nodes in your cluster. Establish a performance benchmark and track performance against this baseline as your environment changes over time.

<!-- ### Supported Integrations

#### **Stable**

| Service Mesh | Status |
| :----------- | -----: |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "stable" -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

##### **Beta**

| Service Mesh | Status |
| :----------- | -----: |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "beta" -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

##### **Alpha** - Service mesh adapters for which we are seeking community-contributed support.

| Service Mesh | Status |
| :----------- | -----: |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "alpha" -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
 -->

<!-- ## Meshery as a project and its community

{% assign sorted_pages = site.pages | sort: "type" | reverse %}

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="project" and item.language=="en" and item.list != "exclude" %}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.description != " " %}
        -  {{ item.description }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul> -->
