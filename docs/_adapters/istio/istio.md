---
layout: page
title: Istio
name: Meshery Adapter for Istio
version: v1.7.3
port: 10000/tcp
project_status: stable
github_link: https://github.com/layer5io/meshery-istio
image: /docs/assets/img/service-meshes/istio.svg
---

# {{ page.name }}

|  Service Mesh  |                   Adapter Status                    | Latest Supported Mesh Version |
| :------------: | :-------------------------------------------------: | :---------------------------: |
| {{page.title}} | [{{ page.project_status }}]({{ page.github_link }}) |       {{page.version}}        |

### Features

1. Lifecycle management of Istio
1. Lifecycle management of sample applications
1. SMI Conformance Capability of Istio
1. Configuration best practices
1. Custom service mesh configuration
1. Prometheus and Grafana connections

### Lifecycle management

The Meshery Adapter for Istio can install **v1.7.3** of the Istio service mesh. The SMI adapter for Istio can also be installed using Meshery.

### Sample applications

The ({{page.name}}) includes a handful of sample applications. Use Meshery to deploy any of these sample applications:

- [Bookinfo](https://github.com/istio/istio/tree/master/samples/bookinfo)
- [Httpbin](https://httpbin.org/)
- [Hipster](https://github.com/GoogleCloudPlatform/microservices-demo)

Once BookInfo is deployed, you can use Meshery to apply configuration to control traffic, inject latency, perform context-based routing, and so on.

### SMI Conformance Capability

Meshery allows you to analyze the compliance status and functional capabilities of your service mesh. This allows you to compare high-level functional differences between service meshes and verify whether your service mesh is conformant with the SMI specification.

Learn more about the SMI specification and [Meshery's conformance test suite](https://meshery.layer5.io/docs/functionality/smi-conformance).

### Configuration best practices

The {{page.name}} will parse all of Istio's configuration and compare the running configuration of the service mesh against known best practices for an {{page.title}} deployment.

### Custom service mesh configuration

Meshery allows you to paste (or type in) any Kubernetes manifest that you would like to have applied to the cluster. This configuraiton may be new VirtualServices or new DestinationRules or other.

![Custom Istio Configuration in Meshery]({{ relative_url }}istio-adapter-custom-configuration.png)

### Prometheus and Grafana connections

The {{page.name}} will connect to Istio's Prometheus and Grafana instances running in the control plane (typically found in the `istio-system` namespace). You can also connect Meshery to Prometheus and Grafana instances not running in the service mesh's control plane.

### Suggested Topics

- Examine [Meshery's architecture]({{ site.baseurl }}/architecture) and how adapters fit in as a component.
- Learn more about [Meshery Adapters]({{ site.baseurl }}/architecture/adapters).
