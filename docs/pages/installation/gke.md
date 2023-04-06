---
layout: default
title: GKE
permalink: installation/platforms/gke
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/gke.png
---

{% include installation_prerequisites.html %}

To set up and run Meshery on GKE 

In order to provide Meshery with the necessary access to your managed Kubernetes instance, 
Meshery will need to be assigned a *ServiceAccount*. An existing ServiceAccount can be used or a new one created. Ensure that the ServiceAccount you use has the *cluster-admin* role assigned.

Meshery will use this *ServiceAccount* to interact with your managed Kubernetes instance. Use either of the following two methods to prepare a compatible kubeconfig file:

- [Automatic Configuration](#automatic-configuration-recommended)
- [Manual configuration](#manual-configuration-optional)

### Automatic Configuration (Recommended)

1. In your browser, navigate to Meshery (e.g., `http://localhost:9081`) and login.
1. Download your Meshery authentication token by clicking Get Token under your user profile.
1. Use this authentication token to execute the following command:
    
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">$ mesheryctl system config gke --token *PATH_TO_TOKEN*</div></div>
 </pre>

This command updates your kubeconfig to provide Meshery with access to your managed Kubernetes instance.
Once configured, proceed with using Meshery:
`mesheryctl system start`

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
`MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start`

### Manual Configuration (Optional)

If the [Automatic Configuration](#automatic-configuration-recommended) procedure fails or you would like to manually prepare your kubeconfig file to provide Meshery with the necessary access to your managed Kubernetes instance, perform the following actions:

1. Download the [generate_kubeconfig_gke.sh](./generate_kubeconfig_gke.sh) shell script.
1. Execute this shell script identifying ServiceAccount name and Namespace arguments, like so:
    
    <pre class="codeblock-pre"><div class="codeblock">
    <div class="clipboardjs">./generate_kubeconfig_gke.sh cluster-admin-sa-gke default</div></div>
    </pre>
    
1. Once the script is complete, you may proceed to start Meshery with the GKE-compatible configuration by executing:
    
    <pre class="codeblock-pre"><div class="codeblock">
    <div class="clipboardjs">$ mesheryctl system start</div></div>
    </pre>
    - Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
    <pre class="codeblock-pre"><div class="codeblock">
    <div class="clipboardjs">$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start</div></div>
    </pre>
1. In your browser, navigate to Meshery (e.g., `http://localhost:9081`) and login.
1. Under Settings --> Environment, provide the generated file (config-cluster-admin-sa-gke-default.yaml) as the kubeconfig file.

Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.