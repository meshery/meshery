---
layout: page
title: App Mesh Adapter
name: Meshery Adapter for App Mesh
mesh_name: App Mesh
version: v0.1.0
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

The {{page.name}} can install **{{page.version}}** of the {{page.mesh_name}} service mesh. 

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
