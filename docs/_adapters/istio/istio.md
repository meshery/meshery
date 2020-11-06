---
layout: page
title: Istio
name: Meshery Adapter for Istio
version: v1.7.3
mesh_name: Istio
port: 10000/tcp
project_status: stable
github_link: https://github.com/layer5io/meshery-istio
image: /docs/assets/img/service-meshes/istio.svg
---

# {{ page.name }}

|  Service Mesh  |                   Adapter Status                    | Latest Supported Mesh Version |
| :------------: | :-------------------------------------------------: | :---------------------------: |
| <img src="{{ page.image }}" style="width:20px" /> {{ page.title }} | [{{ page.project_status }}]({{ page.github_link }}) |       {{page.version}}        |

### Lifecycle management

The {{page.name}} can install **{{page.version}}** of the {{page.title}} service mesh. The SMI adapter for Istio can also be installed using Meshery.

### Install {{ page.mesh_name }}

##### **Choose the Meshery adapter for {{ page.mesh_name }}**

<a href="#istio-adapter">
  <img style="width:500px;" src="/docs/assets/img/adapters/istio/istio-adapter.png" />
</a>
<a href="#" class="lightbox" id="istio-adapter">
  <span style="background-image: url('/docs/assets/img/adapters/istio/istio-adapter.png')"></span>
</a>

##### **Click on (+) and choose the `{{page.version}}` of the {{page.mesh_name}} service mesh.**

<a href="#istio-install">
  <img style="width:500px;" src="/docs/assets/img/adapters/istio/istio-install.png" />
</a>
<a href="#" class="lightbox" id="istio-install">
  <span style="background-image: url('/docs/assets/img/adapters/istio/istio-install.png')"></span>
</a>

### Features

1. Lifecycle management of Istio
1. Lifecycle management of sample applications
1. SMI Conformance Capability of Istio
1. Configuration best practices
1. Custom service mesh configuration
1. Prometheus and Grafana connections

### Sample applications

The ({{page.name}}) includes a handful of sample applications. Use Meshery to deploy any of these sample applications:

- [Bookinfo](/docs/guides/sample-apps#bookinfo)
    - Follow this [tutorial workshop](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md) to set up and deploy the BookInfo sample app on Istio using Meshery. 

- [Httpbin](/docs/guides/sample-apps#httpbin)
    - Httpbin is a simple HTTP request and response service.

- [Hipster](/docs/guides/sample-apps#hipster)
    - Hipster Shop Application is a web-based, e-commerce demo application from the Google Cloud Platform.

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
