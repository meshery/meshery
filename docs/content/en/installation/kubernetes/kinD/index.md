---
title: KinD
categories: [kubernetes]
aliases:
- /installation/platforms/kind
display_title: false
image: installation/kubernetes/kinD/images/kind.png
description: Install Meshery on KinD. Deploy Meshery in KinD in-cluster or outside of KinD out-of-cluster.
---

<h1>Quick Start with KinD <img src="images/kind.png" style="width:35px;height:35px;" /></h1>

Manage your KinD clusters with Meshery. Deploy Meshery in your [KinD cluster](#in-cluster-installation).

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://kind.sigs.k8s.io/docs/user/quick-start/#installation">KinD</a> on your local machine.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Alternative Installation: Using Helm](#alternative-installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

## In-cluster Installation

Follow the steps below to install Meshery in your KinD cluster.

### Preflight Checks

Before you deploy Meshery on KinD, check the following.

#### Preflight: Cluster Connectivity

If you don't have a KinD cluster yet, create one:

{{< code code="kind create cluster" >}}

List your KinD clusters:

{{< code code="kind get clusters" >}}

Verify your kubeconfig's current context.

{{< code code="kubectl config current-context" >}}

### Installation: Using `mesheryctl`

<details>
<summary>Verify your Meshery context</summary>
<p>
Verify that your current Meshery context is set for an in-cluster deployment (<code>platform: kubernetes</code>) by executing:
</p>

{{< code code="mesheryctl system context view" >}}

<p>
If the platform is not <code>kubernetes</code>, create a new context with Kubernetes as the platform:
</p>

{{< code code="mesheryctl system context create context-name --platform kubernetes --url http://localhost:9081 --set --yes" >}}

</details>

With your KinD cluster as your current context, start Meshery:

{{< code code="mesheryctl system start -p kubernetes" >}}

### Alternative Installation: Using Helm

See [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

### Post-Installation Steps

Meshery creates its service as type `LoadBalancer` by default, and KinD does not assign external addresses to `LoadBalancer` services on its own. The simplest option is `mesheryctl system dashboard --port-forward`. A KinD-specific option is [Cloud Provider KIND](https://kind.sigs.k8s.io/docs/user/loadbalancer/), a standalone binary that runs on your host, connects to your KinD cluster, and provisions load balancer containers for your services.

{{% mesheryctl/system-dashboard %}}

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}