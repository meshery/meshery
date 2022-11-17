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

{% include installation_prerequisites.html %}

## Set up and run Meshery on EKS:

The following guide will help you in installing Meshery and making it work with EKS clusters.<br/>
In order to install Meshery, `mesheryctl` needs to access your EKS cluster. For this, you must have:

- Any one of the [aws CLIs](https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html), for managing EKS, installed and configured to use your resources.
- A [valid kubeconfig](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html) for your cluster.

Once you fulfil the above two requirements, you should be able to install Meshery:

- [Install Meshery into your EKS cluster](#install-meshery-into-your-eks-cluster)
- [Install Meshery in Docker and connect it to your EKS cluster](#install-meshery-in-docker-and-connect-it-to-your-eks-cluster)

**_Note: It is advisable to [Install Meshery into your EKS clusters](#install-meshery-into-your-eks-cluster)_**

### Install Meshery into your EKS cluster

To set the context to Kubernetes

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ mesheryctl system context create [context-name] -p kubernetes -s
 </div></div>
 </pre>
 <br>

Execute the following to start Meshery

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ mesheryctl system start 
 </div></div>
 </pre>
 <br>

<button class="toggle-button" onclick="HideToggleFunction()">Optional</button> Customize your Meshery Provider Callback URL

<div id="hiddendiv">
Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start
 </div></div>
 </pre>
 <br>
Meshery should now be running in your EKS cluster and Meshery UI should be accessible at the `EXTERNAL IP` of `meshery` service.
<br/>
</div>
<br/>
Also see: [Install Meshery into Kubernetes](https://docs.meshery.io/installation/platforms/kubernetes)

### Install Meshery in Docker and connect it to your EKS cluster

**_Note: Out-of-cluster support for EKS is still beta and on [roadmap](https://github.com/meshery/meshery/blob/master/ROADMAP.md)._**

Install Meshery in Docker

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ mesheryctl system start -p docker 
 </div></div>
 </pre>
 <br>

Configure Meshery to connect to your cluster by executing:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 $ mesheryctl system config eks
 </div></div>
 </pre>
 <br>
Once you have verified that all the services are up and running, Meshery UI will be accessible on your local machine on port 9081. Open your browser and access Meshery at [`http://localhost:9081`](http://localhost:9081).
