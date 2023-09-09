---
layout: default
title: Project
permalink: project
redirect_from: project/
language: en
display-title: "false"
list: exclude
---

# Meshery Overview

As the collaborative cloud native manager plane, Meshery enables the design, operation, and management of Kubernetes clusters and their workloads. Meshery's powerful performance management functionality is accomplished through implementation of [Cloud Native Performance](https://smp-spec.io). <!-- Meshery's cloud native manager functionality leverages [Service Mesh Interface](https://smi-spec.io) (SMI) and Meshery is the conformance tool for SMI. --> Meshery Catalog enable users to capture and share cloud native design patterns. Meshery enables operators to deploy WebAssembly filters to Envoy-based data planes. Meshery facilitates learning about functionality and performance of infrastructure and incorporates the collection and display of metrics from applications.

##### **Meshery as a project and its community**

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
</ul>
---

## Meshery's Functionality

Meshery features can be categorized by:

1. Cloud Native Performance Management
   - Workload and infrastructure performance characterization
   - Prometheus and Grafana integration
1. Cloud Native Configuration Management
   - Configuration best practices
   - Configuration support for hundreds of [integrations](/extensibility/integrations)
1. Cloud Native Lifecycle Management
   - Cloud native provisioning and workload onboarding
   - Cloud native design patterns and sharing via catalog
1. Data Plane Intelligence
   - Registry and configuration of WebAssembly filters for Envoy
   - Filter chaining and publishing to content catalog
1. Cloud Native Interoperability and Federation
   - Manage multiple Kubernetes clusters concurrently

### Meshery is for Developers, Operators, and Product Owners

Whether making a Day 0 adoption choice or maintaining a Day 2 deployment, Meshery has useful capabilities in either circumstance. Targeted audience for Meshery project would be any technology operators that leverage cloud native infrastructure in their ecosystem; this includes developers, devops engineers, decision makers, architects, and organizations that rely on microservices platform.

### Meshery is for cloud native design patterns

Meshery integrates with hundreds of [infrastructure and application components](/extensibility/integrations) to enable users to deploy cloud native design patterns.

### Meshery is for performance management

Meshery helps users weigh the value of their workload deployments against the overhead incurred while running cloud native infrastructure. Meshery provides statistical analysis of the request latency and throughput seen across various permutations of your workload, infrastructure and configuration.
In addition to request latency and throughput, Meshery also tracks memory and CPU overhead in of the nodes in your cluster. Measure your data plane and control plane against different sets of workloads and infrastructures.

Anytime performance questions are to be answered, they are subjective to the specific workload and infrastructure used for measurement. Given this challenge, the Envoy project, for example, refuses to publish performance data because such tests can be:

- Involved
- Misinterpreted

Beyond the need for performance and overhead data under a permutation of different workloads (applications) and types and sizes of infrastructure resources, the need for cross-project, apple-to-apple comparisons are also desired in order to facilitate a comparison of behavioral differences between cloud native infrastructure and selection of their use. Individual projects shy from publishing test results of other, competing projects. An independent, unbiased, credible analysis is needed.

Meshery is vendor and project-neutral. Meshery's performance management features allow you to uniformly benchmarking the performance of cloud native infrastructure. Meshery allows you to pick an efficient set of tools for your ecosystem by providing performance evaluation and metrics.

1. By leveraging Meshery you could achieve apples-to-apples performance comparison of cloud native infrastructure
2. Track your instructure and application performance from release to release.
3. Understand behavioral differences between releases of your infrastructure.
4. Track your application performance from version to version.

<a href="https://raw.githubusercontent.com/meshery/meshery/master/docs/assets/img/readme/meshery_lifecycle_management.png"><img alt="Meshery Lifecycle Management" src="{{ site.baseurl }}{% link assets/img/readme/meshery_lifecycle_management.png %}"  width="100%" align="center"/></a>
Establish a performance benchmark and track performance against this baseline as your environment changes over time.

## Meshery is for all cloud native infrastructure

Infrastructure diversity is a reality for any enterprise. Whether you're running a single Kubernetes cluster or multiple clusters, you'll find that Meshery supports your infrastructure diversity (or lack thereof).

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
