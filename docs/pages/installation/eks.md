---
layout: default
title: EKS
permalink: installation/platforms/eks
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/eks.png
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your EKS clusters with Meshery. Deploy Meshery in EKS [in-cluster](#in-cluster-installation) or outside of EKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to [Install Meshery in your EKS clusters](#install-meshery-into-your-eks-cluster)_**

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html">AWS CLI</a>, configured for your environment.</li>
    <li>Access to an active EKS cluster in AWS Account.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes]({{ site.baseurl }}/installation/platforms/kubernetes)
## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
- [Post-Installation Steps](#post-installation-steps)
  - [Access Meshery UI](#access-meshery-ui)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Install Meshery on Docker](#install-meshery-on-docker)

# In-cluster Installation

Follow the steps below to install Meshery in your EKS cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on EKS.

### Preflight: Cluster Connectivity

1. Verfiy you connection to an Elastic Kubernetes Services Cluster using AWS CLI.
1. Login to AWS account using [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html), if you are using a different method of authentication in AWS, please refer to AWS documentation.
1. After successful login, set the cluster context.
{% capture code_content %}aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]{% endcapture %}
{% include code.html code=code_content %}
1. _Optional:_ If you are using `eksctl`, follow the [AWS documentation steps](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html).
1. Verify your kubeconfig's current context.
{% capture code_content %}kubectl cluster-info{% endcapture %}
{% include code.html code=code_content %}

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [port-forwarding]({{ site.baseurl }}/reference/mesheryctl/system/dashboard) guide for detailed instructions.
2. If you are using a LoadBalancer, please refer to the [LoadBalancer]({{ site.baseurl }}/installation/platforms/kubernetes#exposing-meshery-serviceloadbalancer) guide for detailed instructions.
3. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

Meshery should now be running in your EKS cluster and Meshery UI should be accessible at the `EXTERNAL IP` of `meshery` service.

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your EKS cluster. Configure Meshery to connect to your EKS cluster by executing:

{% capture code_content %}$ mesheryctl system config eks{% endcapture %}
{% include code.html code=code_content %}

Once configured, execute the following command to start Meshery.

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/helm) guide.

# Post-Installation Steps

## Access Meshery UI

To access Meshery's UI via port-forwarding, please refer to the [port-forwarding](/services/port-forward) guide for detailed instructions.

# Out-of-cluster Installation

{% include alert.html title='Out-of-cluster EKS deployments not currently supported' type="warning" alert='Out-of-cluster support for EKS is still beta and on <a href="https://github.com/meshery/meshery/blob/master/ROADMAP.md">roadmap</a>.' %}

Install Meshery on Docker (out-of-cluster) and connect it to your EKS cluster.

## Install Meshery on Docker

{% capture code_content %}$ mesheryctl system start -p docker{% endcapture %}
{% include code.html code=code_content %}

Configure Meshery to connect to your cluster by executing:

{% capture code_content %}$ mesheryctl system config eks{% endcapture %}
{% include code.html code=code_content %}

Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).