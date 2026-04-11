---
title: EKS
category: [kubernetes]
aliases:
- /installation/platforms/eks
display_title: false
image: /installation/kubernetes/images/eks.png
description: Install Meshery on Elastic Kubernetes Service. Deploy Meshery in EKS in-cluster or outside of EKS out-of-cluster.
---

<h1>Quick Start with EKS <img src="/installation/kubernetes/images/eks.png" style="width:35px;height:35px;" /></h1>

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
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Install Meshery on Docker](#install-meshery-on-docker)

# In-cluster Installation

Follow the steps below to install Meshery in your EKS cluster.

**Prerequisites: Cluster Connectivity**

1. Verify your connection to an Elastic Kubernetes Service cluster using the AWS CLI.
2. Log in to your AWS account using [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html). If you are using a different method of authentication in AWS, please refer to AWS documentation.
3. After successful login, set the cluster context.

{{< code code="aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]" >}}
1. _Optional:_ If you are using `eksctl`, follow the [AWS documentation steps](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html).
1. Verify your kubeconfig's current context.

{{< code code="kubectl config current-context" >}}

## Installation: Using `mesheryctl`

Execute <a href='/reference/mesheryctl/system/start'>mesheryctl system start</a> command to start Meshery.

{{< code code="mesheryctl system start" >}}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< installation/accessing-meshery-ui >}}

# Out-of-cluster Installation

{{% alert title="Out-of-cluster EKS deployments not currently supported" color="warning" %}}
Out-of-cluster support for EKS is still beta and on the [roadmap](https://github.com/meshery/meshery/blob/master/ROADMAP.md).
{{% /alert %}}

Install Meshery on Docker (out-of-cluster) and connect it to your EKS cluster.

## Install Meshery on Docker

{{< code code="mesheryctl system start -p docker" >}}

Configure Meshery to connect to your cluster by executing:

{{< code code="mesheryctl system config eks" >}}

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).

{{< related-discussions tag="meshery" >}}
