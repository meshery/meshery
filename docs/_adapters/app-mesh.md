---
layout: page
title: App Mesh Adapter
name: Meshery Adapter for App Mesh
version: "-"
port: 10005/tcp
project_status: alpha
github_link: https://github.com/layer5io/meshery-app-mesh
image: /docs/assets/img/service-meshes/aws-app-mesh.png
---

# {{ page.name }}

| Service Mesh   | Adapter Status | Latest Supported Mesh Version |
| :------------: | :------------:   | :------------:              |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) | {{page.version}}  |

### Lifecycle management

The {{page.name}} can install the App Mesh service mesh. The SMI adapter for App Mesh can also be installed using Meshery.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
