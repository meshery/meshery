---
layout: page
title: Overview
permalink: overview
---
# Introducing Meshery
The service mesh management plans adopting, operating and developing on different service meshes. 
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



- **Available service mesh adapters** - Service mesh adapters that Meshery currently supports.

| Platform      &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Status        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "stable" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }})                     |       {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
<br>
- **In-progress service mesh adapters** - Service mesh adapters for which community-contributed support has been committed and are currently under development.

| Platform      &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Status        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "beta" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
<br>
- **Help-wanted service mesh adapters** - Service mesh adapters adapters for which we are seeking community-contributed support.

| Platform          &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;     | Status        |
| :------------               | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "alpha" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

## Community
This project is community-built and welcomes collaboration! [Fork here on Github](https://github.com/layer5io/meshery)

* Join [weekly community meeting](https://docs.google.com/document/d/1c07UO9dS7_tFD-ClCWHIrEzRnzUJoFQ10EzfJTpS7FY/edit?usp=sharing) [Fridays from 10am to 11am Central](/assets/projects/meshery/Meshery-Community-Meeting.ics). 
  * Watch community [meeting recordings](https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0) and subscribe to the [community calendar](https://bit.ly/2SbrRhe).

* Access the [community drive](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) (request access).

# FAQ 

## Why use Meshery?
* Because its an open source, vendor neutral projects that facilitates testing across meshes.
* Because fortio is not packaged into a mesh testing utility, but is only a load-generator unto its own.
* Because regpatrol is closed source, binary is not released, scripted for one mesh, and is produced by a vendor of that mesh.

## Why create Meshery and not use another benchmark tool?
Meshery is purpose built for facilitating benchmarking of service meshes and their workloads. Other benchmark tools are not. There are some other tools used for service mesh benchmarking, like regpatrol. Regpatrol is used by IBM is not open source or available in binary form to use and has the following differences from Meshery:
- Telemetry - regpatrol sources telemetry from the Mixer Prometheus adapter and uses IBM's proprietary node agent.
- Meshery sources from the Mixer Prometheus adapter and uses Prometheus node-exporter.
- Traffic type - regpatrol uses jmeter, which can parse responses and perform functional tests.
- Meshery is using fortio, which is for load-gen and perf-testing only.

# Resources

## Meshery Presentations

- [O'Reilly OSCON 2020](https://conferences.oreilly.com/oscon/oscon-or)
- [O'Reilly Infrastructure & Ops 2020](https://conferences.oreilly.com/infrastructure-ops/io-ca/public/schedule/speaker/226795)
- [InnoTech Dallas 2020](https://innotechdallas2020.sched.com/event/aN7E/a-management-plane-for-service-meshes)
- [KubeCon EU 2020](https://kccnceu20.sched.com/event/Zetg/discreetly-studying-the-effects-of-individual-traffic-control-functions-lee-calcote-layer5?iframe=no&w=100%&sidebar=yes&bg=no)
- DockerCon 2020 ([deck](https://calcotestudios.com/talks/decks/slides-dockercon-2020-service-meshing-with-docker-desktop-and-webassembly.html), [video](https://www.youtube.com/watch?v=5BrbbKZOctw&list=PL3A-A6hPO2IN_HSU0pSfijBboiHggs5mC&index=4&t=0s))
- [Open Source 101 at Home](https://calcotestudios.com/talks/decks/slides-open-source-101-at-home-solving-the-service-mesh-adopters-dilemma.html)
- [Docker Captains Roundtable 2020](https://calcotestudios.com/talks/decks/slides-docker-captains-2020-meshery-the-multi-service-mesh-manager.html)
- [Cloud Native Austin 2020](https://www.meetup.com/Cloud-Native-Austin/events/267784090/)
- NSMCon 2019 talk ([video](https://www.youtube.com/watch?v=4xKixsDTtdM), [deck](https://calcotestudios.com/talks/decks/slides-nsmcon-kubecon-na-2019-adopting-network-service-mesh-with-meshery.html))
- [Service Mesh Day 2019](https://youtu.be/CFj1O_uyhhs)
- [DockerCon 2019 Open Source Summit](https://www.docker.com/dockercon/2019-videos?watch=open-source-summit-service-mesh)
- KubeCon EU 2019 ([video](https://www.youtube.com/watch?v=LxP-yHrKL4M&list=PLYjO73_1efChX9NuRaU7WocTbgrfvCoPE), [deck](https://calcotestudios.com/talks/decks/slides-kubecon-eu-2019-service-meshes-at-what-cost.html))
- [KubeCon EU 2019 Istio Founders Meetup](https://calcotestudios.com/talks/decks/slides-istio-meetup-kubecon-eu-2019-istio-at-scale-large-and-small.html)
- [Cloud Native Rejekts EU 2019](https://calcotestudios.com/talks/decks/slides-cloud-native-rejekts-2019-evaluating-service-meshes.html)
- [Container World 2019](https://calcotestudios.com/talks/decks/slides-container-world-2019-service-meshes-but-at-what-cost.html)
- Solving the Service Mesh Adopter’s Dilemma ([deck](https://calcotestudios.com/talks/decks/slides-open-source-101-at-home-solving-the-service-mesh-adopters-dilemma.html), [event](https://opensource101.com/sessions/solving-the-service-mesh-adopters-dilemma/),[video](https://www.youtube.com/watch?v=Q1zSWbO0RmI&list=PL3A-A6hPO2IN_HSU0pSfijBboiHggs5mC&index=2&t=0s))

## Other Resources
- [Service Mesh Comparison](https://layer5.io/landscape)
- [Service Mesh Tools](https://layer5.io/landscape#tools)
- [Service Mesh Books](https://layer5.io/books)
- [Service Mesh Workshops](https://layer5.io/workshops)
