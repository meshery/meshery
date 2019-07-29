---
layout: page
title: Architecture
permalink: /architecture
nav_order: 3
---

## Table of contents
{: .no_toc }

1. TOC
{:toc}

---
### Architecture

<div class="iframe-container">
 <iframe src="https://docs.google.com/presentation/d/e/2PACX-1vSj6eYr6AgZ4mBgOL_Gv9T4WyLBFkPv49asNtdw1_Gn_xCsk37QRhOjdBRB-3Jp1ehneFmm2dpgFie-/embed?start=false&loop=false&delayms=3000#slide=id.g55c4016581_0_0" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>
</div> 

### Value provided by Meshery Cloud
Performance test reports area anonymously sent to Meshery Cloud for statistical analysis and sharing back to the community. Data points like the average overhead of service meshes as seen in various environments are shared through this mechanism.
<div class="iframe-container">
    <iframe src="https://docs.google.com/presentation/d/e/2PACX-1vSj6eYr6AgZ4mBgOL_Gv9T4WyLBFkPv49asNtdw1_Gn_xCsk37QRhOjdBRB-3Jp1ehneFmm2dpgFie-/embed?start=false&loop=false&delayms=3000#slide=id.g4f68f671f0_0_0" frameborder="0" width="960" height="569" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"></iframe>

</div> 

### Network Ports (that Meshery uses and needs)
- Meshery (thru docker-compose) needs port 9081
- Meshery Istio adapter port 10000
- Linkerd adapter port 10001
- Consul 10002

### Adapters (what they are, what purpose they serve)
* Istio
* Linkerd
* Octarine
* Consul Connect
