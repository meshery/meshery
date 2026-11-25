---
title: EKS
categories: [kubernetes]
aliases:
- /installation/platforms/eks
display_title: false
image: installation/kubernetes/eks/images/eks.png
description: Install Meshery on Elastic Kubernetes Service. Deploy Meshery in EKS in-cluster or outside of EKS out-of-cluster.
---

<h1>Quick Start with EKS <img src="images/eks.png" style="width:35px;height:35px;" /></h1>

Manage your EKS clusters with Meshery. Deploy Meshery inside your EKS cluster ([in-cluster](#in-cluster-installation)) or outside of it ([out-of-cluster](#out-of-cluster-installation)). **_Note: We recommend installing Meshery in your EKS clusters._**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="{{< ref "installation/mesheryctl/_index.md" >}}" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html">AWS CLI</a>, configured for your environment.</li>
    <li>Access to an active EKS cluster in an AWS account.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})

### Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Install Meshery on Docker](#install-meshery-on-docker)

# In-cluster Installation

Follow the steps below to install Meshery in your EKS cluster.

**Prerequisites: Cluster Connectivity**

Use the AWS CLI to connect `kubectl` to your Elastic Kubernetes Service (EKS) cluster.

1. Sign in to your AWS account with [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html). If you authenticate to AWS another way, see the AWS documentation.
2. Add the cluster to your kubeconfig and make it your current context:

{{< code code="aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]" >}}
3. _Optional:_ If you use `eksctl`, follow the [AWS documentation](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html).
4. Verify your kubeconfig's current context:

{{< code code="kubectl config current-context" >}}

## Installation: Using `mesheryctl`

Run <a href='{{< ref "reference/mesheryctl/system/start.md" >}}'>mesheryctl system start</a> to deploy Meshery to your EKS cluster:

{{< code code="mesheryctl system start -p kubernetes" >}}

## Installation: Using Helm

To install Meshery with Helm 3, see the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide.

## Post-Installation Steps

Optionally, verify the health of your Meshery deployment with <a href='{{< ref "reference/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery. Open Meshery UI in your browser.

{{< installation/accessing-meshery-ui >}}

# Out-of-cluster Installation

{{% alert title="Out-of-cluster EKS deployments not currently supported" color="warning" %}}
Out-of-cluster support for EKS is still beta and on the [roadmap](https://github.com/meshery/meshery/blob/master/ROADMAP.md).
{{% /alert %}}

Install Meshery on Docker (out-of-cluster) and connect it to your EKS cluster.

## Install Meshery on Docker

{{< code code="mesheryctl system start -p docker" >}}

Once Meshery is running, sign in and connect it to your cluster. `mesheryctl connection create` prompts for the AWS region and cluster name, writes the cluster's kubeconfig to `~/.meshery/kubeconfig.yaml`, and uploads it to Meshery Server. (`mesheryctl system config eks` is deprecated.)

{{< code code=`mesheryctl system login
mesheryctl connection create --type eks` >}}

Meshery UI is available on your local machine at [`http://localhost:9081`](http://localhost:9081).

{{< related-discussions tag="meshery" >}}
