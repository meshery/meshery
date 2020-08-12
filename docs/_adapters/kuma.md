---
layout: page
title: Kuma Adapter
name: Kuma
version: v1.0
port: 10007/tcp
project_status: alpha
github_link: https://github.com/layer5io/meshery-kuma
image: /docs/assets/img/service-meshes/kuma.svg
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
