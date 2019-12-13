---
layout: page
title: Architecture
permalink: architecture
---

# Architecture

<div class="iframe-container">
 <iframe src="https://docs.google.com/presentation/d/e/2PACX-1vSj6eYr6AgZ4mBgOL_Gv9T4WyLBFkPv49asNtdw1_Gn_xCsk37QRhOjdBRB-3Jp1ehneFmm2dpgFie-/embed?start=false&loop=false&delayms=3000#slide=id.g55c4016581_0_0" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>
</div> 

# Value provided by Meshery Cloud
Performance test reports area anonymously sent to Meshery Cloud for statistical analysis and sharing back to the community. Data points like the average overhead of service meshes as seen in various environments are shared through this mechanism.
<div class="iframe-container">
    <iframe src="https://docs.google.com/presentation/d/e/2PACX-1vSj6eYr6AgZ4mBgOL_Gv9T4WyLBFkPv49asNtdw1_Gn_xCsk37QRhOjdBRB-3Jp1ehneFmm2dpgFie-/embed?start=false&loop=false&delayms=3000#slide=id.g4f68f671f0_0_0" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>

</div> 

## Network Ports 
Meshery uses the following list of network ports to interface with its various components:

| Adapter | Port |
|:-------|:----|
| Meshery web-based UI | 9081/tcp |
| [Meshery Adapter for Istio](/docs/installation/adapters/istio) | 10000/tcp |
| [Meshery Adapter for Linkerd](/docs/installation/adapters/linkerd) | 10001/tcp |
| [Meshery Adapter for Consul](/docs/installation/adapters/consul) | 10002/tcp |
| [Meshery Adapter for Octarine](/docs/installation/adapters/octarine) | 10003/tcp |
| [Meshery Adapter for NSM](/docs/installation/adapters/nsm) | 10004/tcp |
| [Meshery Adapter for App Mesh](/docs/installation/adapters/app-mesh) |10005/tcp |
| [Meshery Adapter for Maesh](/docs/installation/adapters/maesh) |10006/tcp |
| [Meshery Adapter for Kuma](/docs/installation/adapters/kuma) |10007/tcp |
| [Meshery Adapter for Citrix CPX](/docs/installation/adapters/cpx) |10008/tcp |

See the [Adapters](installation/adapters) section for more information on the function of an adapter.
