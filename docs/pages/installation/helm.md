---
layout: default
title: Helm
permalink: installation/kubernetes/helm
type: installation
category: kubernetes
redirect_from:
- installation/platforms/helm
display-title: "false"
language: en
list: include
---
# Install Meshery on Kubernetes Using Helm

<div class="prereqs"><h4>Prerequisites</h4>
<ol>
<li><a href="https://helm.sh/docs/intro/install/" class="meshery-light">Helm</a> should be installed on your local machine.</li>
<li>You should have access to the cluster/platform where you want to deploy Meshery.</li>
<li>Ensure that the kubeconfig file has the correct current context/cluster configuration.</li>
</ol>
</div>

## Install Meshery on Your Kubernetes Cluster Using Helm

{% capture code_content %}helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace
helm install meshery-operator meshery/meshery-operator{% endcapture %}
{% include code.html code=code_content %}

Optionally, Meshery Server supports customizing the callback URL for your remote provider, like so:

{% capture code_content %}helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace{% endcapture %}
{% include code.html code=code_content %}

### Accessing Meshery UI for Clusters

To access Meshery's UI , please refer to the [accessing-ui](/tasks/accessing-meshery-ui) guide for detailed instructions.

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}
