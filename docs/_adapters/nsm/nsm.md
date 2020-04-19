---
layout: page
title: Network Service Mesh Adapter
name: Network Service Mesh
version: v0.2.1
port: 10004/tcp
project_status: stable
---
| Adapter Status |
| :------------: |
| [{{ page.project_status }}]({{ page.github_link }})|

# Sample Applications

The Meshery adapter for Network Service Mesh includes a handful of sample applications. These applications represent different network services orchestrated by Network Service Mesh. Use Meshery to deploy any of these sample applications.

## Hello NSM Application

Watch this presentation to see the Hello NSM Application in-action:
<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/4xKixsDTtdM" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

See on YouTube: [Adopting Network Service Mesh with Meshery](https://www.youtube.com/watch?v=4xKixsDTtdM&list=PL3A-A6hPO2IOpTbdH89qR-4AE0ON13Zie)

## VPP-ICMP Application

_A simple example that connects a vpp based Pod to a Network Service using memif._

The simplest possible case for Network Service Mesh is to have is connecting a Client via a vWire to another Pod that is providing a Network Service.
Network Service Mesh allows flexibility in the choice of mechanisms used to provide that vWire to a workload.

The icmp responder example does this with kernel interfaces.  The vpp-icmp-responder provides and consumes the same 'icmp-responder' Network Service, but has Client's and Endpoint's that use a [memif](https://www.youtube.com/watch?v=6aVr32WgY0Q) high speed memory interfaces to achieve performance unavailable via kernel interfaces.

![vpp-icmp-responder-example](/docs/adapters/nsm/vpp-icmp-responder-example.svg)

**What it does**

This will install two Deployments:

Name | Description
:--------|:--------
**vpp-icmp-responder-nsc** | The Clients, four replicas
**vpp-icmp-responder-nse** | The Endpoints, two replicas

And cause each Client to get a vWire connecting it to one of the Endpoints.  Network Service Mesh handles the Network Service Discovery and Routing, as well as the vWire 'Connection Handling' for setting all of this up.

![vpp-icmp-responder-example-2](/docs/adapters/nsm/vpp-icmp-responder-example-2.svg)

In order to make this case more interesting, Endpoint1 and Endpoint2 are deployed on two separate Nodes using PodAntiAffinity, so that the Network Service Mesh has to demonstrate the ability to string vWires between Clients and Endpoints on the same Node and Clients and Endpoints on different Nodes.

**Verify**

First verify that the vpp-icmp-responder example Pods are all up and running:

```bash
kubectl get pods | grep vpp-icmp-responder
```

To see the vpp-icmp-responder example in action, you can run:

```bash
curl -s https://raw.githubusercontent.com/networkservicemesh/networkservicemesh/master/scripts/nsc_ping_all.sh | bash
```

## ICMP Responder

The simplest possible case for Network Service Mesh is to have is connecting a Client via a vWire to another Pod that is providing a Network Service. We call this case the ‘icmp-responder’ example, because it allows the client to ping the IP address of the Endpoint over the vWire.

![icmp-responder-example](/docs/adapters/nsm/icmp-responder-example.svg)

**What it Does**

This will install two Deployments:

Name | Description |
:--------|:--------
**icmp-responder-nsc** | The Clients, four replicas |
**icmp-responder-nse** | The Endpoints, two replicas |

And cause each Client to get a vWire connecting it to one of the Endpoints.  Network Service Mesh handles the
Network Service Discovery and Routing, as well as the vWire 'Connection Handling' for setting all of this up.

![icmp-responder-example-2](/docs/adapters/nsm/icmp-responder-example-2.svg)

In order to make this case more interesting, Endpoint1 and Endpoint2 are deployed on two separate Nodes using
PodAntiAffinity, so that the Network Service Mesh has to demonstrate the ability to string vWires between Clients and
Endpoints on the same Node and Clients and Endpoints on different Nodes.

**Verifying**

First verify that the icmp-responder example Pods are all up and running:

```bash
kubectl get pods | grep icmp-responder
```

To see the icmp-responder example in action, you can run:

```bash
curl -s https://raw.githubusercontent.com/networkservicemesh/networkservicemesh/master/scripts/nsc_ping_all.sh | bash
```
