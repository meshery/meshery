---
title: "EKS"
description: "Install Meshery on Elastic Kubernetes Service. Deploy Meshery in-cluster or out-of-cluster."
weight: 30
aliases:
  - /installation/platforms/eks
image: /images/platforms/eks.png
display_title: "false"
---

<h1>Quick Start with EKS <img src="/images/platforms/eks.png" style="width:35px;height:35px;" /></h1>

Manage your EKS clusters with Meshery. Deploy Meshery in EKS [in-cluster](#in-cluster-installation) or outside of EKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your EKS clusters_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html">AWS CLI</a>, configured for your environment.</li>
    <li>Access to an active EKS cluster in an AWS account.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes](/installation/kubernetes)

### Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
- [Out-of-cluster Installation](#out-of-cluster-installation)

# In-cluster Installation

Follow the steps below to install Meshery in your EKS cluster.

**Prerequisites: Cluster Connectivity**

1. Verify your connection to an Elastic Kubernetes Service cluster using the AWS CLI.
2. Log in to your AWS account using [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html).
3. After successful login, set the cluster context:

{{< code >}}
aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]
{{< /code >}}

4. Verify your kubeconfig's current context:

{{< code >}}
kubectl config current-context
{{< /code >}}

## Installation: Using `mesheryctl`

Execute [mesheryctl system start](/reference/mesheryctl/system/start) command to start Meshery:

{{< code >}}
$ mesheryctl system start
{{< /code >}}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using [mesheryctl system check](/reference/mesheryctl/system/check).

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< accessing-meshery-ui >}}

# Out-of-cluster Installation

Install Meshery on Docker (out-of-cluster) and connect it to your EKS cluster.

## Install Meshery on Docker

{{< code >}}
$ mesheryctl system start -p docker
{{< /code >}}

Configure Meshery to connect to your cluster:

{{< code >}}
$ mesheryctl system config eks
{{< /code >}}

{{< related-discussions tag="meshery" >}}
