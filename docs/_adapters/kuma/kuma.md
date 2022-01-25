---
layout: page
title: Meshery Adapter for Kuma
name: Meshery Adapter for Kuma
mesh_name: Kuma
earliest_version: v1.2.2
port: 10007/gRPC
project_status: stable
adapter_version: v0.5.2
lab: kuma-meshery-adapter
github_link: https://github.com/meshery/meshery-kuma
image: /assets/img/service-meshes/kuma.svg
permalink: service-meshes/adapters/kuma
---

{% include adapter-status.html %}

{% include adapter-labs.html %}

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of {{page.mesh_name}} service mesh. A number of sample applications can be installed using the {{page.name}}.

### Install {{ page.mesh_name }}

Choose the Meshery Adapter for {{ page.mesh_name }}.

<a href="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-adapter.png" />
</a>

Click on (+) and choose the {{page.earliest_version}} of the {{page.mesh_name}} service mesh.

<a href="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/kuma/kuma-install.png" />
</a>

## Workload Management

The following sample applications are available in this adapter.

- [Bookinfo]({{site.baseurl}}/guides/sample-apps#bookinfo)
  - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.

### SMI Conformance Capability

Meshery allows you to analyze the compliance status and functional capabilities of your service mesh. This allows you to compare high-level functional differences between service meshes and verify whether your service mesh is conformant with the SMI specification.

Learn more about the SMI specification and [Meshery's conformance test suite]({{ site.baseurl }}/functionality/service-mesh-interface).
