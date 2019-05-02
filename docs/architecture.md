---
layout: page
title: Architecture
permalink: /architecture
parent: Meshery
nav_order: 3
---

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---
### Diagram
![Meshery Architecture](assets/images/arch.jpg)

### Value provided by meshery cloud
![Meshery Hosted](assets/images/hosted.jpg)
### Network Ports (that Meshery uses and needs)
- Meshery (thru docker-compose) needs port 9081
- Meshery Istio adapter port 10000
- Linkerd adapter port 10001
### Adapters (what they are, what purpose they serve
* Istio
* Linkerd
* Octarine
* Consul Connect
