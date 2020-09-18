---
layout: page
title: Open Service Mesh Adapter
name: Meshery Adapter for Open Service Mesh
mesh_name: Open Service Mesh
version: OSM v0.3.0
port: 10009/tcp
project_status: alpha
github_link: https://github.com/layer5io/meshery-osm
image: /docs/assets/img/service-meshes/osm.svg
---

# {{ page.name }}

|  Service Mesh  |                   Adapter Status                    | Latest Supported Mesh Version |
| :------------: | :-------------------------------------------------: | :---------------------------: |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) |       {{page.version}}        |

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the {{page.mesh_name}} service mesh. A number of sample applications for {{page.mesh_name}} can also be installed using Meshery. Using the Meshery adapter for Open Service Mesh, you can run the full suite of Service Mesh Interface (SMI) conformance tests to validate OSM’s compliance with the SMI specification.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).

### SMI Conformance Capabiliy

### Conformance

Defining "Conformance" - Acknowledging that conformance consists of both capabilities and compliance status is important. We define conformance as a combination of these two concepts.

    1. SMI Conformance acknowledge that ...some participating service meshes may conscientiously never fully implement functions (SMI specs)...
    2. SMI Conformance identifies ...a difference between full implementation of a specification and compliance with the portions that it implements...

You can [Learn more](https://meshery.layer5.io/docs/functionality/smi-conformance) about SMI conformance capabilities here...
