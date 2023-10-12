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

Manage your EKS clusters with Meshery. Deploy Meshery on EKS or outside of EKS.
{% include installation_prerequisites.html %}

### General Prerequisites:
1. Access to an active EKS cluster in AWS Account.
2. Any one of the [aws CLIs](https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html), for managing EKS, installed and configured to use your resources.
3. Ensure [kubectl](https://kubernetes.io/docs/tasks/tools/) is installed on your local machine.

## Connect to an Elastic Kubernetes Services Cluster using AWS CLI: 

Connect to the EKS Cluster by following these steps:

1. Install [AWS CLI(awscli)](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) if not installed, and login to aws account using [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-authentication-user.html), if you are using a different method of authentication in AWS please refer to the AWS Docs.

2. After successful login, set the cluster context
{% capture code_content %}aws eks update-kubeconfig --name [YOUR_CLUSTER_NAME] --region [YOUR_REGION]{% endcapture %}
{% include code.html code=code_content %}

[Optional]If you are using eksctl, follow the steps [here](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html)

3. Verify the current context of the cluster:
{% capture code_content %}kubectl cluster-info{% endcapture %}
{% include code.html code=code_content %}

## Install Meshery on an EKS Cluster using mesheryctl

### Connect to Elastic Kubernetes Services Cluster using mesheryctl

Use Meshery's CLI to streamline your connection to your EKS cluster. Configure Meshery to connect to your EKS cluster by executing:

{% capture code_content %}$ mesheryctl system config eks{% endcapture %}
{% include code.html code=code_content %}
 <br>

Once configured, install Meshery with this command:

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}
 <br>

<button class="toggle-button" onclick="HideToggleFunction()">Optional</button> Customize your Meshery Provider Callback URL

<div id="hiddendiv">
Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way
<br>
{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}
 <br>
Meshery should now be running in your EKS cluster and Meshery UI should be accessible at the `EXTERNAL IP` of `meshery` service.
<br/>
</div>
<br/>

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## [Alternative] Install Meshery on your Elastic Kubernetes Services Cluster using Helm V3

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/helm) guide.
<br />

### [Optional] Access Meshery UI

To access Meshery's UI via port-forwarding, please refer to the [port-forwarding](/services/port-forward) guide for detailed instructions.

### Install Meshery in Docker and connect it to your EKS cluster

**_Note: Out-of-cluster support for EKS is still beta and on [roadmap](https://github.com/meshery/meshery/blob/master/ROADMAP.md)._**

Install Meshery in Docker

{% capture code_content %}$ mesheryctl system start -p docker{% endcapture %}
{% include code.html code=code_content %}
 <br>

Configure Meshery to connect to your cluster by executing:

{% capture code_content %}$ mesheryctl system config eks{% endcapture %}
{% include code.html code=code_content %}
 <br>
Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).