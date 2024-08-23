---
layout: default
title: EKS
permalink: installation/kubernetes/eks
type: installation
category: kubernetes
redirect_from:
- installation/platforms/eks
display-title: "false"
language: en
list: include
image: /assets/img/platforms/eks.png
abstract: Install Meshery on Elastic Kubernetes Service. Deploy Meshery in EKS in-cluster or outside of EKS out-of-cluster.
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your EKS clusters with Meshery. Deploy Meshery in EKS [in-cluster](#in-cluster-installation) or outside of EKS [out-of-cluster](#out-of-cluster-installation). **_Note: It is advisable to [Install Meshery in your EKS clusters](#install-meshery-into-your-eks-cluster)_**

<div class="prereqs"><h4>Prerequisites</h4>
  <ol>
    <li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
    <li>Install <a href="https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html">AWS CLI</a>, configured for your environment.</li>
    <li>Access to an active EKS cluster in AWS Account.</li>
  </ol>
</div>

Also see: [Install Meshery on Kubernetes]({{ site.baseurl }}/installation/kubernetes)

### Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
  - [Post-Installation Steps](#post-installation-steps)
- [Out-of-cluster Installation](#out-of-cluster-installation)
  - [Install Meshery on Docker](#install-meshery-on-docker)

# In-cluster Installation

Follow the steps below to install Meshery in your EKS cluster.

**Prerequisites: Cluster Connectivity**

1. Verify your connection to an Elastic Kubernetes Services Cluster using AWS CLI.
1. Login to AWS account using [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html), if you are using a different method of authentication in AWS, please refer to AWS documentation.
1. After successful login, set the cluster context.
{% capture code_content %}aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]{% endcapture %}
{% include code.html code=code_content %}
1. _Optional:_ If you are using `eksctl`, follow the [AWS documentation steps](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html).
1. Verify your kubeconfig's current context.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your EKS cluster. Configure Meshery to connect to your EKS cluster by executing:

{% capture code_content %}$ mesheryctl system config eks{% endcapture %}
{% include code.html code=code_content %}

Once configured, execute the following command to start Meshery.

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md display-title="true" %}

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

{% include related-discussions.html tag="meshery" %}
