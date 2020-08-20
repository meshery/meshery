---
layout: page
title: Open Service Mesh Adapter
name: Open Service Mesh
version: "-"
port: 10010/tcp
project_status: alpha
github_link: https://github.com/layer5io/meshery-osm
image: /docs/assets/img/service-meshes/osm.svg
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

