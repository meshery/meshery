---
layout: page
title: Kuma Adapter
name: Meshery Adapter for Kuma
version: v1.0
port: 10007/tcp
project_status: beta
github_link: https://github.com/layer5io/meshery-kuma
image: /docs/assets/img/service-meshes/kuma.svg
---

# {{ page.name }}

|  Service Mesh  |                   Adapter Status                    | Latest Supported Mesh Version |
| :------------: | :-------------------------------------------------: | :---------------------------: |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) |       {{page.version}}        |

### Features

1. Lifecycle management of Kuma
1. SMI Conformance Capability of Kuma

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of {{page.mesh_name}}. A number of sample applications can be installed using the {{page.name}}.

### SMI Conformance Capability

Meshery allows you to analyze the compliance status and functional capabilities of your service mesh. This allows you to compare high-level functional differences between service meshes and verify whether your service mesh is conformant with the SMI specification.

Learn more about the SMI specification and [Meshery's conformance test suite](https://meshery.layer5.io/docs/functionality/smi-conformance).
