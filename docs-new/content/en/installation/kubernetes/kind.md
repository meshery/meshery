---
title: "KinD"
description: "Install Meshery on KinD. Deploy Meshery in-cluster or out-of-cluster."
weight: 50
aliases:
  - /installation/platforms/kind
image: /images/platforms/kind.png
display_title: "false"
---

<h1>Quick Start with KinD <img src="/images/platforms/kind.png" style="width:35px;height:35px;" /></h1>

Manage your KinD clusters with Meshery. Deploy Meshery in your [KinD cluster](#in-cluster-installation).

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://kind.sigs.k8s.io/docs/user/quick-start/#installation">KinD</a> on your local machine.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes](/installation/kubernetes)

## In-cluster Installation

Follow the steps below to install Meshery in your KinD cluster.

### Preflight: Cluster Connectivity

Start KinD if it is not already running:

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

With your KIND cluster configured as your `current-context`, start Meshery:

{{< code >}}
$ mesheryctl system start -p kubernetes
{{< /code >}}

### Alternative Installation: Using Helm

See [Helm Installation](/installation/kubernetes/helm) guide.

### Post-Installation Steps

Meshery deploys with LoadBalancer service type by default. If you are using KinD, you may need to expose the Meshery service using `mesheryctl system dashboard --port-forward`.

Optionally, you can verify the health of your Meshery deployment using [mesheryctl system check](/reference/mesheryctl/system/check).

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

{{< related-discussions tag="meshery" >}}
