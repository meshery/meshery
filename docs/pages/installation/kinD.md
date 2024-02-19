---
layout: default
title: kinD
permalink: installation/kubernetes/kind
type: installation
category: kubernetes
redirect_from:
- installation/platforms/kind
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kind.png
abstract: Install Meshery on kinD. Deploy Meshery in kinD in-cluster or outside of kinD out-of-cluster.
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:35px;height:35px;" /></h1>

Manage your kinD clusters with Meshery. Deploy Meshery in your [kinD cluster](#in-cluster-installation).

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li>Install the Meshery command line client, <a href="{{ site.baseurl }}/installation/mesheryctl" class="meshery-light">mesheryctl</a>.</li>
<li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> on your local machine.</li>
<li>Install <a href="https://kind.sigs.k8s.io/docs/user/quick-start/#installation">kinD</a>, on your local machine.</li>
</ol>
</div>

Also see: [Install Meshery on Kubernetes]({{ site.baseurl }}/installation/kubernetes)

## Available Deployment Methods

- [In-cluster Installation](#in-cluster-installation)
  - [Preflight Checks](#preflight-checks)
    - [Preflight: Cluster Connectivity](#preflight-cluster-connectivity)
    - [Preflight: Plan your access to Meshery UI](#preflight-plan-your-access-to-meshery-ui)
  - [Installation: Using `mesheryctl`](#installation-using-mesheryctl)
  - [Installation: Using Helm](#installation-using-helm)
- [Post-Installation Steps](#post-installation-steps)
  - [Access Meshery UI](#access-meshery-ui)

# In-cluster Installation

Follow the steps below to install Meshery in your kinD cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on kinD.

### Preflight: Cluster Connectivity

Start the kinD, if not started using the following command:
{% capture code_content %}kind create cluster{% endcapture %}
{% include code.html code=code_content %}
Check up on your kinD cluster :
{% capture code_content %}kind get clusters{% endcapture %}
{% include code.html code=code_content %}
Verify your kubeconfig's current context.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [port-forwarding](/tasks/accessing-meshery-ui) guide for detailed instructions.
2. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

Meshery should now be running in your kinD cluster and Meshery UI should be accessible at the `INTERNAL IP` of `meshery` service.

## Installation: Using `mesheryctl`

Once kinD cluster is configured as current cluster-context, execute the below command.

Before executing the below command, go to `~/.meshery/config.yaml` and ensure that current platform is set to kubernetes.
{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

# Post-Installation Steps

## Access Meshery UI

To access Meshery's UI, please refer to the [instruction](/tasks/accessing-meshery-ui) for detailed guidance.

{% if page.suggested-reading != false and page.title and page.type and page.category and page.url %}
{% include_cached suggested-reading.html  title=page.title type=page.type category=page.category url=page.url language="en" %}
{% endif %}

{% include related-discussions.html tag="meshery" %}