---
title: KinD
description: "Install Meshery on KinD. Deploy Meshery in-cluster or out-of-cluster."
weight: 50
aliases:
  - /installation/platforms/kind
image: /images/platforms/kind.png
display_title: false
---

<h1>Quick Start with KinD <img src="/images/platforms/kind.png" style="width:35px;height:35px;" /></h1>

Manage your KinD clusters with Meshery. Deploy Meshery in your [KinD cluster](#in-cluster-installation).

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command-line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://kind.sigs.k8s.io/docs/user/quick-start/#installation">KinD</a> on your local machine.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes](/installation/kubernetes)

## Available Deployment Methods

- [Available Deployment Methods](#available-deployment-methods)
- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Alternative Installation: Using Helm](#alternative-installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)

## In-cluster Installation

Follow the steps below to install Meshery in your KinD cluster.

### Preflight Checks

Read through the following considerations prior to deploying Meshery on KinD.

#### Preflight: Cluster Connectivity

Start KinD if it is not already started using the following command:

{{< code >}}
kind create cluster
{{< /code >}}

Check the status of your KinD cluster:

{{< code >}}
kind get clusters
{{< /code >}}

Verify your kubeconfig's current context:

{{< code >}}
kubectl config current-context
{{< /code >}}

### Installation: Using `mesheryctl`

<details>
<summary>Verify your Meshery context</summary>
<p>
Verify that your current Meshery context is set for an in-cluster deployment (<code>platform: kubernetes</code>) by executing:
</p>

{{< code >}}
mesheryctl system context view
{{< /code >}}

<p>
If the context is not set to <code>platform: kubernetes</code>, you can create a new context with Kubernetes as the platform using the following command.
</p>

{{< code >}}
mesheryctl system context create context-name --platform kubernetes --url http://localhost:9081 --set --yes
{{< /code >}}

<br/>
</details>

With your KIND cluster configured as your `current-context`, start Meshery:

{{< code >}}
mesheryctl system start -p kubernetes
{{< /code >}}

Meshery server supports customizing authentication flow callback URL. To learn more, refer to the [customizing authentication](/extensibility/providers) guide.

### Alternative Installation: Using Helm

For detailed instructions on installing Meshery using Helm, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

### Post-Installation Steps

Meshery deploys with LoadBalancer service type by default. If you are using KinD, you may need to expose the Meshery service. A universal option is to use `mesheryctl system dashboard --port-forward`. A KIND-specific option is to use the [Cloud Provider KIND](https://kind.sigs.k8s.io/docs/user/loadbalancer/). Cloud Provider KIND runs as a standalone binary in your host and connects to your KIND cluster and provisions new LoadBalancer containers for your Services.

**mesheryctl system dashboard**

To open the Meshery UI in your browser, use the following command:

{{< code >}}
mesheryctl system dashboard
{{< /code >}}

To port-forward to Meshery UI, use the following command:

{{< code >}}
mesheryctl system dashboard --port-forward
{{< /code >}}

Optionally, you can verify the health of your Meshery deployment using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}