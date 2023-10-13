---
layout: default
title: Helm Installation
permalink: installation/helm
type: installation
display-title: "true"
language: en
list: exclude
---

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li><a href="https://helm.sh/docs/intro/install/" class="meshery-light">Helm</a>  should be installed on your local machine.</li>
    <li>You should have access to the cluster/platform where you want to deploy Meshery.</li>
    <li>Ensure that the kubeconfig file has the correct current context/cluster configuration.</li>
  </ol>
</div> 

## Install Meshery on Your Kubernetes Cluster Using Helm V3

We strongly recommend using Helm v3 because this version no longer includes the Tiller component, making it lighter and safer.

{% capture code_content %}helm repo add meshery https://meshery.io/charts/
helm install meshery meshery/meshery --namespace meshery --create-namespace{% endcapture %}
{% include code.html code=code_content %}
<br />

Optionally, Meshery Server supports customizing the callback URL for your remote provider, like so:

{% capture code_content %}helm install meshery meshery/meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host --create-namespace{% endcapture %}
{% include code.html code=code_content %}

### Accessing Meshery UI for Clusters

Access the [Meshery UI](/services/port-forward) by port-forwarding it as a Kubernetes service.