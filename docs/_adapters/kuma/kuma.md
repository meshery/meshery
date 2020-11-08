---
layout: page
title: Kuma
name: Meshery Adapter for Kuma
mesh_name: Kuma
version: v1.0
port: 10007/tcp
project_status: beta
github_link: https://github.com/layer5io/meshery-kuma
image: /docs/assets/img/service-meshes/kuma.svg
permalink: service-meshes/adapters/kuma
---

{% include adapter-status.html %}

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of {{page.mesh_name}}. A number of sample applications can be installed using the {{page.name}}.

### Install {{ page.mesh_name }}

##### **Choose the Meshery Adapter for {{ page.mesh_name }}**

<a href="#kuma-adapter">
  <img style="width:500px;" src="/docs/assets/img/adapters/kuma/kuma-adapter.png" />
</a>
<a href="#" class="lightbox" id="kuma-adapter">
  <span style="background-image: url('/docs/assets/img/adapters/citrix/citrix-adapter.png')"></span>
</a>

##### **Click on (+) and choose the `{{page.version}}` of the {{page.mesh_name}} service mesh.**

<a href="#kuma-install">
  <img style="width:500px;" src="/docs/assets/img/adapters/kuma/kuma-install.png" />
</a>
<a href="#" class="lightbox" id="kuma-install">
  <span style="background-image: url('/docs/assets/img/adapters/kuma/kuma-install.png')"></span>
</a>

### Features

1. Lifecycle management of Kuma
1. SMI Conformance Capability of Kuma

### Sample Applications

- [Bookinfo](/docs/guides/sample-apps#bookinfo) 
    - The sample BookInfo application displays information about a book, similar to a single catalog entry of an online book store.

### SMI Conformance Capability

Meshery allows you to analyze the compliance status and functional capabilities of your service mesh. This allows you to compare high-level functional differences between service meshes and verify whether your service mesh is conformant with the SMI specification.

Learn more about the SMI specification and [Meshery's conformance test suite](https://meshery.layer5.io/docs/functionality/smi-conformance).
