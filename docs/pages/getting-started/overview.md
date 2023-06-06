---
layout: default
title: Overview
permalink: /getting-started/overview
language: en
---

As the cloud native management plane, Meshery enables the adoption, operation, and management of Kubernetes clusters and their lifecycle. Meshery's powerful performance management functionality is useful whether your running workloads on Kubernetes or outside of Kubernetes. 

Meshery implmenents both [Service Mesh Performance](https://smp-spec.io) (SMP) and [Service Mesh Interface](https://smi-spec.io) (SMI) and Meshery is the conformance tool for SMI. Meshery integrates with Open Application Model (OAM) to enable users to deploy service mesh patterns. Meshery enables operators to deploy WebAssembly filters to Envoy-based data planes. Meshery facilitates learning about functionality and performance of service meshes and incorporates the collection and display of metrics from applications running on or across service meshes.

## Meshery's Functionality

Meshery features can be categorized by:

1. Cloud Native Performance Management
   - Workload and cloud native performance characterization
   - Prometheus and Grafana integration
1. Cloud Native Configuration Management
   - Configuration best practices
1. Cloud Native Lifecycle Management
   - Cloud native provisioning and workload onboarding
   - Meshery Operator and MeshSync
   - Cloud native patterns and Open Application Model integration
1. Data Plane Intelligence
   - Registry and configuration of WebAssembly filters for Envoy
1. Cloud Native Interoperability and Federation
   - Manage multiple service meshes concurrently
   - Connect to multiple clusters independently

### Meshery is for Developers, Operators, and Product Owners

Whether making a Day 0 adoption choice or maintaining a Day 2 deployment, Meshery has useful capabilities in either circumstance. Targeted audience for Meshery project would be any technology operators that leverage service mesh in their ecosystem; this includes developers, devops engineers, decision makers, architects, and organizations that rely on microservices platform.

### Meshery is for cloud native patterns

Through MeshModel, Meshery describes infrastructure under management, enabling you to define cloud native designs and patterns and then to export those designs and share within the <a href="https://meshery.io">Meshery Catalog</a>. 

### Meshery is for performance management

Meshery helps users weigh the value of their cloud native deployments against the overhead incurred in running different deployment scenarios and different configruations. Meshery provides statistical analysis of the request latency and throughput seen across various permutations of your workload, infrastructure and infrastructure configuration. In addition to request latency and throughput, Meshery also tracks memory and CPU overhead in of the nodes in your cluster. Measure your data plane and control plane against different sets of workloads and infrastructures.

Anytime performance questions are to be answered, they are subjective to the specific workload and infrastructure used for measurement. Given this challenge, many projects refuse to publish their own performance data, because such tests can be quite invovled and misinterpreted.

Beyond the need for performance and overhead data under a permutation of different workloads (applications) and types and sizes of infrastructure resources, the need for cross-project, apple-to-apple comparisons are also desired in order to facilitate a comparison of behavioral differences between cloud antive  and selection of their use. Individual projects shy from publishing test results of other, competing service meshes. An independent, unbiased, credible analysis is needed.

Meshery is intended to be a vendor and project-neutral utility for uniformly benchmarking the performance of cloud native infrastructure. Between service mesh and proxy projects (and surprisingly, within a single project), a number of different tools and results exist. Meshery allows you to pick an efficient set of tools for your ecosystem by providing performance evaluation and metrics.

1. By leveraging Meshery you can achieve apples-to-apples performance comparison
1. Track your service mesh performance from release to release.
1. Understand behavioral differences between cloud native infrastructure.
1. Track your application performance from version to version.

<a href="https://raw.githubusercontent.com/meshery/meshery/master/docs/assets/img/readme/meshery_lifecycle_management.png"><img alt="Meshery Lifecycle Management" src="{{ site.baseurl }}{% link assets/img/readme/meshery_lifecycle_management.png %}"  width="100%" align="center"/></a>
Establish a performance benchmark and track performance against this baseline as your environment changes over time.

## Meshery is for all cloud native infrastructure

Infrastructure diversity is a reality for any enterprise. Whether you're running a single Kubernetes cluster or multiple Kubernetes clusters, you'll find that Meshery supports your infrastructure diversity (or lack thereof).

<!-- ### Supported Service Meshes

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
{% endfor %} -->
