---
layout: page
title: Architecture
permalink: /architecture
parent: Meshery
nav_order: 3
---

## Table of contents
{: .no_toc }

1. TOC
{:toc}

---
### Architecture
<div class="iframe-container">
    <iframe src="https://docs.google.com/presentation/d/e/2PACX-1vRqJ2cFC9LSGVi1ReDjphzOfpWwINPg4__XyMZu85R8KNPnuHmnswEnHnTkY-8FQl-GrjuXpQ7WlCLE/embed?start=false&loop=false&delayms=3000" 
            class="meshery-arch" frameborder="0" height="445px" width="739px" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"
            allowtransparency="true"></iframe>
</div> 

### Value provided by Meshery cloud
<div class="iframe-container">
    <iframe src="https://docs.google.com/presentation/d/e/2PACX-1vRqJ2cFC9LSGVi1ReDjphzOfpWwINPg4__XyMZu85R8KNPnuHmnswEnHnTkY-8FQl-GrjuXpQ7WlCLE/embed?start=false&loop=false&delayms=3000&slide=id.g4f68f671f0_0_0" 
            class="meshery-arch" frameborder="0" height="445px" width="739px" allowfullscreen="true" mozallowfullscreen="true" webkitallowfullscreen="true"
            allowtransparency="true"></iframe>
</div> 

### Network Ports (that Meshery uses and needs)
- Meshery (thru docker-compose) needs port 9081
- Meshery Istio adapter port 10000
- Linkerd adapter port 10001


### Adapters (what they are, what purpose they serve
* Istio
* Linkerd
* Octarine
* Consul Connect
