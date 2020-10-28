---
layout: page
title: Overview
permalink: overview
---

# Introducing Meshery
The service mesh management plane adopting, operating and developing on different service meshes. 
Meshery facilitates learning about functionality and performance of service meshes and incorporates the collection and display of metrics from applications running on or across service meshes.
Meshery provides this high-level functionality:

1. Service Mesh Performance Management
1. Service Mesh Configuration Management
    - Configuration best practices
1. Service Mesh Lifecycle Management
1. Service Mesh Interoperability and Federation

<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/CFj1O_uyhhs" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

<div style="text-align:center;width:100%"><emphasis>Delivered at Service Mesh Day 2019</emphasis></div>


<h2>What challenges does Meshery solve?</h2>
<b>Service mesh management - one or multiple service meshes.</b>

Anytime performance questions are to be answered, they are subjective to the specific workload and infrastructure used for measurement. Given this challenge, the Envoy project, for example, refuses to publish performance data because such tests can be:
- Involved
- Misinterpreted

Beyond the need for performance and overhead data under a permutation of different workloads (applications) and types and sizes of infrastructure resources, the need for cross-project, apple-to-apple comparisons are also desired in order to facilitate a comparison of behavioral differences between service meshes and selection of their use. Individual projects shy from publishing test results of other, competing service meshes. An independent, unbiased, credible analysis is needed.

Meshery is intended to be a vendor and project-neutral utility for uniformly benchmarking the performance of service meshes. Between service mesh and proxy projects (and surprisingly, within a single project), a number of different tools and results exist. Meshery allows you to pick an efficient set of tools for your ecosystem by providing performance evaluation and metrics.

1. By leveraging Meshery you could achieve apples-to-apples performance comparison of service meshes
1. Track your service mesh performance from release to release.
1. Understand behavioral differences between service meshes.
1. Track your application performance from version to version.

## Meshery is for Adopters and Operators
Whether making a Day 0 adoption choice or maintaining a Day 2 deployment, Meshery has useful capabilities in either circumstance. Targeted audience for Meshery project would be any technology operators that leverage service mesh in their ecosystem; this includes developers, devops engineers, decision makers, architects, and organizations that rely on microservices platform. 

## Meshery is for performance management: testing and benchmarking
Meshery helps users weigh the value of their service mesh deployment against the overhead incurred in running a service mesh. Meshery provides statistical analysis of the request latency and throughput seen across various permutations of your workload, infrastructure and service mesh configuration.
In addition to request latency and throughput, Meshery also tracks memory and CPU overhead in of the nodes in your cluster. Measure your data plane and control plane against different sets of workloads and infrastructures.

<a href="https://raw.githubusercontent.com/layer5io/meshery/master/docs/assets/img/readme/meshery_lifecycle_management.png"><img alt="Layer5 Service Mesh Community" src="{{ site.baseurl }}{% link assets/img/readme/meshery_lifecycle_management.png %}"  width="100%" align="center"/></a>
Establish a performance benchmark and track performance against this baseline as your environment changes over time.

## Meshery is for any service mesh
Infrastructure diversity is a reality for any enterprise. Whether you're running a single service mesh or multiple types of service meshes, you'll find that Meshery supports your infrastructure diversity (or lack thereof).

The Meshery service mesh adapters are present in three stages:

##### **Available** - Service mesh adapters that Meshery currently supports.

| Platform      &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Status        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "stable" -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }})                     |       {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

##### **In-progress** - Service mesh adapters which are currently under development.

| Platform      &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Status        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "beta" -%}
| <img src="{{ adapter.image }}" style="width:20px" />  [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

##### **Help-wanted** - Service mesh adapters adapters for which we are seeking community-contributed support.

| Platform          &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;     | Status        |
| :------------               | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "alpha" -%}
| <img src="{{ adapter.image }}" style="width:20px" />  [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}