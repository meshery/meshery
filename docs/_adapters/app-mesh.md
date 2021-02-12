---
layout: default
title: Meshery Adapter for App Mesh
name: Meshery Adapter for App Mesh
mesh_name: App Mesh
version: v0.1.0
port: 10005/tcp
project_status: alpha
github_link: https://github.com/layer5io/meshery-app-mesh
image: /assets/img/service-meshes/app-mesh.svg
permalink: service-meshes/adapters/app-mesh
---

{% include adapter-status.html %}

The {{ page.name }} is currently under construction ({{ page.project_status }} state), which means that the adapter is not functional and cannot be interacted with through the <a href="{{ site.baseurl }}/installation#6-you-will-now-be-directed-to-the-meshery-ui"> Meshery UI </a>at the moment. Check back here to see updates.

Want to contribute? Check our [progress]({{page.github_link}}).

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the {{page.mesh_name}} service mesh. 

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
