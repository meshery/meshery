---
layout: page
title: Meshery Documentation
permalink: /
---
# What is Meshery 
A multi-service mesh management plane adopting, operating and developing on different service meshes. 
Meshery faciliates learning about functionality and performance of service meshes and incorporates collection and display of metrics from applications running on or across service meshes. 
Meshery provides this high-level functionality: 

1. Performance Benchmarking
1. Service Mesh Lifecycle Management
1. Service Mesh Interoperability and Federation

<div style="text-align:center;">
<iframe width="560" height="315" src="https://www.youtube.com/embed/CFj1O_uyhhs" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe><br /><i>Delivered at Service Mesh Day 2019</i></div>

## What challenges does Meshery solve? 
<p style="text-align:center;"><b>Service mesh management - one or multiple service meshes.</b><br /></p>

<p style="margin-bottom:1em; margin-top:1em;">Anytime performance questions are to be answered, they are subjective to the specific workload and infrastructure used for measurement. Given this challenge, the Envoy project, for example, refuses to publish performance data because such tests can be 1) involved and 2) misinterpreted.</p>


<p style="margin-bottom:1em; margin-top:1em;">Beyond the need for performance and overhead data under a permutation of different workloads (applications) and types and sizes of infrastructure resources, the need for cross-project, apple-to-apple comparisons are also desired in order to facilitate a comparison of behavioral differences between service meshes and selection of their use. Individual projects shy from publishing test results of other, competing service meshes. An independent, unbiased, credible analysis is needed.<br /></p>

<p style="margin-bottom:1em; margin-top:1em;">Meshery is intended to be a vendor and project-neutral utility for uniformly benchmarking the performance of service meshes. Between service mesh and proxy projects (and surprisingly, within a single project), a number of different tools and results exist. Meshery allows you to pick an efficient set of tools for your ecosystem by providing performance evaluation and metrics.<br /></p>

1. By leveraging Meshery you could achieve apples-to-apples performance comparision of service meshes
1. Track your service mesh performance from release to release.
1. Understand behavioral differences between service meshes.
1. Track your application performance from version to version.

## Meshery is for Adopters and Operators
Whether making a Day 0 adoption choice or maintaining a Day 2 deployment, Meshery has useful capabilities in either circumstance. Targeted audience for Meshery project would be any technology operators that leverage service mesh in their ecosystem; this includes developers, devops engineers, decision makers, architects, and organizations that rely on microservices platform. 

## Meshery is for performance benchmarking
- Identify permutations of workloads, infrastructure types, and measurements to use for: 
    1. Data plane testing
    1. Control plane testing.
    - Against a fixed set of:
        1. Workload(s)
        1. Infrastructure(s)

## Supported Service Meshes
See the full list of [adapters](installation/adapters).

| Platform      | Status        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

**Available service mesh adapters** - Service mesh adapters that Meshery currently supports.

**In-progress service mesh adapters** - Service mesh adapters for which community-contributed support has been committed and are currently under development.

**Help-wanted service mesh adapters** - Service mesh adapters adapters for which we are seeking community-contributed support.

# Contributing

## Community
This project is community-built and welcomes collaboration! [Fork here on Github](https://github.com/layer5io/meshery)

* Join [weekly community meeting](https://docs.google.com/document/d/1c07UO9dS7_tFD-ClCWHIrEzRnzUJoFQ10EzfJTpS7FY/edit?usp=sharing) [Fridays from 10am to 11am Central](/assets/projects/meshery/Meshery-Community-Meeting.ics). 
  * Watch community [meeting recordings](https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0) and subscribe to the [community calender](https://bit.ly/2SbrRhe).

* Access the [community drive](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) (request access).

# FAQ 

## Why use Meshery?
* because its an open source, vendor neutral projects that facilitates testing across meshes.
* because fortio is not packaged into a mesh testing utility, but is only a load-generator unto its own.
* because regpatrol is closed sourcej, binary is not released, scripted for one mesh, and is produced by a vendor of that mesh.

## Why create Meshery and not use another benchmark tool?
<p style="margin-bottom:1em; margin-top:1em;">Meshery is purpose built for factilitating benchmarking of service meshes and their workloads. Other benchmark tools are not. There are some other tools used for service mesh benchmarking, like regpatrol. Regpatrol is used by IBM is not open source or available in binary form to use and has the following differences from Meshery:</p>
- Telemetry - regpatrol sources telemetry from Mixer Prometheus adapter and uses IBM's proprietary node agent.
- Meshery sources from Mixer Prometheus adapter and uses Prometheus node-exporter.
- Traffic type - regpatrol uses jmeter, which can parse responses and perform functional tests.
- Meshery is using fortio, which is for load-gen and perf-testing only.

# Resources

## Meshery Presentations
- [Service Mesh Day 2019](https://youtu.be/CFj1O_uyhhs)
- [DockerCon 2019 Open Source Summit](https://www.docker.com/dockercon/2019-videos?watch=open-source-summit-service-mesh)
- KubeCon EU 2019 ([video](https://www.youtube.com/watch?v=LxP-yHrKL4M&list=PLYjO73_1efChX9NuRaU7WocTbgrfvCoPE), [deck](https://calcotestudios.com/talks/decks/slides-kubecon-eu-2019-service-meshes-at-what-cost.html))
- [KubeCon EU 2019 Istio Founders Meetup](https://calcotestudios.com/talks/decks/slides-istio-meetup-kubecon-eu-2019-istio-at-scale-large-and-small.html)
- [Cloud Native Rejekts EU 2019](https://calcotestudios.com/talks/decks/slides-cloud-native-rejekts-2019-evaluating-service-meshes.html)
- [Container World 2019](https://calcotestudios.com/talks/decks/slides-container-world-2019-service-meshes-but-at-what-cost.html)

## Other Resources
- [Service Mesh Comparison](https://layer5.io/landscape)
- [Service Mesh Tools](https://layer5.io/landscape#tools)
- [Service Mesh Books](https://layer5.io/books)
- [Service Mesh Workshops](https://layer5.io/workshops)
