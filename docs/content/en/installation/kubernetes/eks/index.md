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

Manage your EKS clusters with Meshery. Deploy Meshery in EKS [in-cluster](#in-cluster-installation) or outside of EKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to install Meshery in your EKS clusters_**

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

1. Verify your connection to an Elastic Kubernetes Service cluster using the AWS CLI.
2. Log in to your AWS account using [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html). If you are using a different method of authentication in AWS, please refer to AWS documentation.
3. After successful login, set the cluster context.

{{< code code="aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]" >}}

1. _Optional:_ If you are using `eksctl`, follow the [AWS documentation steps](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html).
1. Verify your kubeconfig's current context.

{{< code code="kubectl config current-context" >}}

## Installation: Using `mesheryctl`

Execute <a href='{{< ref "reference/mesheryctl/system/start.md" >}}'>mesheryctl system start</a> command to start Meshery.

{{< code code="mesheryctl system start" >}}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation]({{< ref "installation/kubernetes/helm.md" >}}) guide. This guide is common for all cloud providers including AWS, Azure, GCP, and Minikube cluster services.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment using <a href='{{< ref "reference/mesheryctl/system/check.md" >}}'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{{< installation/accessing-meshery-ui >}}

For more information on accessing Meshery UI, including port-forwarding options, see the [Accessing Meshery UI]({{< ref "installation/kubernetes/_index.md#accessing-meshery-ui" >}}) section.

## Accessing Meshery UI

After installing Meshery in-cluster, you can access the Meshery UI using one of several methods depending on your environment and requirements.

### Option 1: LoadBalancer Service (Recommended for Production)

When Meshery is installed in-cluster, Meshery UI is served by a Kubernetes `Service` named `meshery` in the `meshery` namespace. This `Service` is created as type `LoadBalancer` by default, forwarding port `9081` (Meshery UI) to the Meshery Server container on port `8080`.

On EKS, a `LoadBalancer` `Service` instructs AWS to provision an external load balancer and assign it a routable `EXTERNAL-IP`. Retrieve the address with:

{{< code code="kubectl get service meshery --namespace meshery" >}}

Once the `EXTERNAL-IP` column shows an address instead of `<pending>`, open Meshery UI in your browser at `http://[EXTERNAL-IP]:9081`.

If the `Service` was previously set to another type (for example, `ClusterIP`), switch it back to `LoadBalancer`:

{{< code code=`kubectl patch service meshery --namespace meshery --type merge -p '{"spec":{"type":"LoadBalancer"}}'` >}}

{{% alert title="EXTERNAL-IP not assigned?" color="warning" %}}
A `LoadBalancer` `Service` is only assigned an external address when the cluster runs a load balancer controller. AWS provides one out of the box for EKS; however, if you encounter issues, ensure your EKS cluster has the AWS load balancer controller installed. Alternatively, you can use port-forwarding instead.
{{% /alert %}}

### Option 2: Port Forwarding (Universal)

Port forwarding is a universal method to access Meshery UI that works across all Kubernetes environments, including local clusters and managed clouds. This is particularly useful when:

- A LoadBalancer is not available or practical
- You're working in a local development environment
- You need temporary access without exposing Meshery externally
- You're behind a firewall or in a restricted network

For detailed instructions on using port forwarding, see the [Port Forwarding Guide]({{< ref "installation/kubernetes/port-forwarding.md" >}}).

**Quick example using mesheryctl:**

{{< code code="mesheryctl system dashboard --port-forward" >}}

**Quick example using kubectl:**

{{< code code="kubectl port-forward svc/meshery 9081:9081 --namespace meshery" >}}

Then open `http://localhost:9081` in your browser.

{{% alert title="Port Forwarding vs LoadBalancer" color="info" %}}

- **LoadBalancer**: Best for production deployments requiring external access on EKS
- **Port Forwarding**: Best for development, testing, and local access

Choose the method that best fits your use case. You can use both methods simultaneously if needed.
{{% /alert %}}

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
