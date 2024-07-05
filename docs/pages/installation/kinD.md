---
layout: default
title: KinD
permalink: installation/kubernetes/kind
type: installation
category: kubernetes
redirect_from:
- installation/platforms/kind
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kind.png
abstract: Install Meshery on KinD. Deploy Meshery in KinD in-cluster or outside of KinD out-of-cluster.
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your KinD clusters with Meshery. Deploy Meshery in your [KinD cluster](#in-cluster-installation).

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://kind.sigs.k8s.io/docs/user/quick-start/#installation">KinD</a>, on your local machine.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{ site.baseurl }}/installation/kubernetes)

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

Start the KinD, if not started using the following command:
{% capture code_content %}kind create cluster{% endcapture %}
{% include code.html code=code_content %}
Check up on your KinD cluster :
{% capture code_content %}kind get clusters{% endcapture %}
{% include code.html code=code_content %}
Verify your kubeconfig's current context.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

### Installation: Using `mesheryctl`

<details>
<summary>Verify your Meshery context</summary>
<p>
Verify that your current Meshery context is set for an in-cluster deployment (`platform: kubernetes`) by executing:
</p>

{% capture code_content %}$ mesheryctl system context view{% endcapture %}
{% include code.html code=code_content %}
<p>
If the context is not set to <code>platform: kubernetes</code>, you can create a new context with Kubernetes as the platform using the following command.
</p>

{% capture code_content %}$ mesheryctl system context create context-name --platform kubernetes --url http://localhost:9081 --set --yes{% endcapture %}
{% include code.html code=code_content %}
<br/>
</details>

With your KIND cluster configured your `current-context`, start Meshery.

{% capture code_content %}$ mesheryctl system start -p kubernetes{% endcapture %}
{% include code.html code=code_content %}

### Alternative Installation: Using Helm

See [Helm Installation](/installation/kubernetes/helm) guide.

### Post-Installation Steps

Meshery deploys with LoadBalancer service type by default. If you are using KinD, you may need to expose the Meshery service. A universal option is to use `mesheryctl system dashboard --port-forward`. A KIND-specific option to use use the [Cloud Provider KIND](https://kind.sigs.k8s.io/docs/user/loadbalancer/). Cloud Provider KIND runs as a standalone binary in your host and connects to your KIND cluster and provisions new Load Balancer containers for your Services.

{% include mesheryctl/system-dashboard.md %}

Optionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md display-title="true" %}

{% include related-discussions.html tag="meshery" %}