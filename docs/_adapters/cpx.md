---
layout: page
title: Citrix Service Mesh (CPX) Adapter
name: Citrix Service Mesh
version: "1.0"
port: 10008/tcp
project_status: beta
github_link: https://github.com/layer5io/meshery-cpx
image: /docs/assets/img/service-meshes/cpx.png
---
# {{ page.name }}

| Service Mesh   | Adapter Status | Latest Supported Mesh Version |
| :------------: | :------------:   | :------------:              |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) | {{page.version}}  |

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the {{page.name}} service mesh. The SMI adapter for Kuma can also be installed using Meshery.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
