---
layout: enhanced
title: Meshery Adapter for Network Service Mesh
name: Meshery Adapter for Network Service Mesh
mesh_name: Network Service Mesh
earliest_version: v0.2.1
port: 10004/gRPC
project_status: stable
lab: nsm-meshery-adapter
github_link: https://github.com/meshery/meshery-nsm
image: /assets/img/service-meshes/nsm.svg
white_image: /assets/img/service-meshes/nsm.svg
permalink: extensibility/adapters/nsm
redirect_from: service-meshes/adapters/nsm
language: en
---

{% assign sorted_tests_group = site.compatibility | group_by: "meshery-component" %}
{% for group in sorted_tests_group %}
{% if group.name == "meshery-nsm" %}
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

## Lifecycle management of {{ page.name }}

The {{page.name}} can install **{{page.earliest_version}}** of {{page.mesh_name}} service mesh. A number of sample applications can be installed using the {{page.name}}.

### Install {{ page.mesh_name }}

##### Choose the Meshery Adapter for {{ page.mesh_name }}

<a href="{{ site.baseurl }}/assets/img/adapters/nsm/nsm-adapter.png">
  <img style="width:500px; background: white" src="{{ site.baseurl }}/assets/img/adapters/nsm/nsm-adapter.png" />
</a>

Click on (+) and choose the {{page.earliest_version}} of the {{page.mesh_name}} service mesh.

<a href="{{ site.baseurl }}/assets/img/adapters/nsm/nsm-install.png">
  <img style="width:500px; background: white" src="{{ site.baseurl }}/assets/img/adapters/nsm/nsm-install.png" />
</a>

## Sample Applications

The ({{ page.name }}) includes a handful of sample applications. These applications represent different network services orchestrated by {{page.mesh_name}}. Use Meshery to deploy any of these sample applications:

### 1. Hello NSM Application

Watch this presentation to see the Hello NSM Application in-action:

<iframe width="560" height="315" src="https://www.youtube.com/embed/4xKixsDTtdM" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

See on YouTube: [Adopting Network Service Mesh with Meshery](https://www.youtube.com/watch?v=4xKixsDTtdM&list=PL3A-A6hPO2IOpTbdH89qR-4AE0ON13Zie)

#### 2. VPP-ICMP Application

_A simple example that connects a vpp based Pod to a Network Service using memif._

The simplest possible case for {{page.mesh_name}} is to have is connecting a Client via a vWire to another Pod that is providing a Network Service.
Network Service Mesh allows flexibility in the choice of mechanisms used to provide that vWire to a workload.

The icmp responder example does this with kernel interfaces. The vpp-icmp-responder provides and consumes the same 'icmp-responder' Network Service, but has Client's and Endpoint's that use a [memif](https://www.youtube.com/watch?v=6aVr32WgY0Q) high speed memory interfaces to achieve performance unavailable via kernel interfaces.

<a href="{{ site.baseurl }}/assets/img/adapters/nsm/vpp-icmp-responder-example.svg">
  <img style="width:500px; background: white" src="{{ site.baseurl }}/assets/img/adapters/nsm/vpp-icmp-responder-example.svg" />
</a>

Working process

This will install two Deployments:

| Name                   | Description                  |
| :--------------------- | :--------------------------- |
| vpp-icmp-responder-nsc | The Clients (four replicas)  |
| vpp-icmp-responder-nse | The Endpoints (two replicas) |

And cause each Client to get a vWire connecting it to one of the Endpoints. Network Service Mesh handles the Network Service Discovery and Routing, as well as the vWire 'Connection Handling' for setting all of this up.

<a href="{{ site.baseurl }}/assets/img/adapters/nsm/vpp-icmp-responder-example-2.svg">
  <img style="width:500px; background: white" src="{{ site.baseurl }}/assets/img/adapters/nsm/vpp-icmp-responder-example-2.svg" />
</a>

In order to make this case more interesting, Endpoint1 and Endpoint2 are deployed on two separate Nodes using PodAntiAffinity, so that the Network Service Mesh has to demonstrate the ability to string vWires between Clients and Endpoints on the same Node and Clients and Endpoints on different Nodes.

Verification

First verify that the vpp-icmp-responder example Pods are all up and running:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ kubectl get pods | grep vpp-icmp-responder
 </div></div>
 </pre>

To see the vpp-icmp-responder example in action, you can run:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ curl -s https://raw.githubusercontent.com/networkservicemesh/networkservicemesh/master/scripts/nsc_ping_all.sh | bash
 </div></div>
 </pre>

#### 3. ICMP Responder

The simplest possible case for Network Service Mesh is to have is connecting a Client via a vWire to another Pod that is providing a Network Service. We call this case the ‘icmp-responder’ example, because it allows the client to ping the IP address of the Endpoint over the vWire.

<a href="{{ site.baseurl }}/assets/img/adapters/nsm/vpp-icmp-responder-example.svg">
  <img style="width:500px; background: white" src="{{ site.baseurl }}/assets/img/adapters/nsm/vpp-icmp-responder-example.svg" />
</a>

Outcomes

This will install two Deployments:

| Name               | Description                 |
| :----------------- | :-------------------------- |
| icmp-responder-nsc | The Clients, four replicas  |
| icmp-responder-nse | The Endpoints, two replicas |

And cause each Client to get a vWire connecting it to one of the Endpoints. Network Service Mesh handles the
Network Service Discovery and Routing, as well as the vWire 'Connection Handling' for setting all of this up.

<a href="{{ site.baseurl }}/assets/img/adapters/nsm/icmp-responder-example-2.svg">
  <img style="width:500px; background: white" src="{{ site.baseurl }}/assets/img/adapters/nsm/icmp-responder-example-2.svg" />
</a>

In order to make this case more interesting, Endpoint1 and Endpoint2 are deployed on two separate Nodes using
_PodAntiAffinity_, so that the Network Service Mesh has to demonstrate the ability to string vWires between Clients and
Endpoints on the same Node and Clients and Endpoints on different Nodes.

Verification

1. Verify that the icmp-responder example Pods are all up and running:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ kubectl get pods | grep icmp-responder
 </div></div>
 </pre>

2. To see the icmp-responder example in action, you may run:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ curl -s https://raw.githubusercontent.com/networkservicemesh/networkservicemesh/master/scripts/nsc_ping_all.sh | bash
 </div></div>
 </pre>

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).

