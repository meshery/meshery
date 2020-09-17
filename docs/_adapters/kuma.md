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

### Conformance

Defining "Conformance" - It is important to acknowledge that conformance consists of both capabilities and compliance status. We define conformance as a combination of these two concepts.

1. SMI Conformance acknowledges that
   ...some participating service meshes may conscientiously never fully implement functions (SMI specs)...

2. SMI Conformance identifies
   ...a difference between full implementation of a specification and compliance with the portions that it implements...

[Learn more](https://meshery.layer5.io/docs/functionality/smi-conformance) about SMI conformance capabilities here...
