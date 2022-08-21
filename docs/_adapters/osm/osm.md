---
layout: default
title: Meshery Adapter for Open Service Mesh
name: Meshery Adapter for Open Service Mesh
mesh_name: Open Service Mesh
port: 10009/gRPC
project_status: stable
lab: osm-meshery-adapter
github_link: https://github.com/meshery/meshery-osm
image: /assets/img/service-meshes/osm.svg
permalink: service-meshes/adapters/osm
earliest_version: v0.9.2
---

{% include adapter-status.html %}

{% include adapter-labs.html %}

## Lifecycle management

The {{page.name}} can install **{{page.earliest_version}}** of the {{page.mesh_name}} service mesh. Sample applications for {{page.mesh_name}} can also be installed using Meshery. Using the {{page.name}}, you may also run the complete suite of Service Mesh Interface (SMI) conformance tests to validate OSM’s compliance with the SMI specification.

### Install {{ page.mesh_name }}

Choose the Meshery Adapter for {{ page.mesh_name }}.

<a href="{{ site.baseurl }}/assets/img/adapters/osm/osm-adapter.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/osm/osm-adapter.png" />
</a>

Click on (+) and choose the `{{page.earliest_version}}` of the {{page.mesh_name}} service mesh.

<a href="{{ site.baseurl }}/assets/img/adapters/osm/osm-install.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/adapters/osm/osm-install.png" />
</a>

## Using Service Mesh Standards

### Complying with Service Mesh Interface (SMI)

Meshery allows you to analyze the compliance status and functional capabilities of your service mesh. This allows you to compare high-level functional differences between service meshes and verify whether your service mesh is conformant with the SMI specification.

Learn more about the SMI specification and [Meshery's conformance test suite]({{ site.baseurl }}/functionality/service-mesh-interface).

### Managing Service Mesh Performance (SMP)

### Prometheus and Grafana connections

The {{page.name}} will connect to Open Service Mesh's Prometheus and Grafana instances running in the control plane (typically found in the `osm-system` namespace). You can also connect Meshery to Prometheus and Grafana instances not running in the service mesh's control plane.

### Sample Applications

The {{ page.name }} does not support Sample Applications yet.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
