---
layout: default
title: Rancher(rke2)
permalink: installation/kubernetes/rancher
type: installation
category: kubernetes
redirect_from:
- installation/platforms/rancher
display-title: "false"
language: en
list: include
image: /assets/img/platforms/rke21.svg
abstract: Install Meshery on RKE2
---

<h1>Quick Start with {{ page.title }} <img src="{{ page.image }}" style="width:100px;height:100px;vertical-align:middle;" /></h1>

Manage your Rke2 clusters with Meshery. Deploy Meshery in Rke2 [in-cluster](#in-cluster-installation) .

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Install the Rke2 cluster in machine, <a href="https://docs.rke2.io/" class="meshery-light"></a>.</li>
    <li>Install <a href="https://kubernetes.io/docs/tasks/tools/">kubectl</a> installed on your local machine.</li>
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
  - [Installation: Manual Steps](#installation-manual-steps)
  - [Installation: Upload Config File in Meshery Web UI](#installation-upload-config-file-in-meshery-web-ui)
  - [Post-Installation Steps](#post-installation-steps)

# In-cluster Installation

Follow the steps below to install Meshery in your Rke2 cluster.

## Preflight Checks

Read through the following considerations prior to deploying Meshery on Rke2.

### Preflight: Cluster Connectivity

Start the rke2 cluster, if not started using the following command:
{% capture code_content %}systemctl start rke2-server.service{% endcapture %}
{% include code.html code=code_content %}
Check up on your rke2 cluster :
{% capture code_content %}systemctl status rke2-server {% endcapture %}
{% include code.html code=code_content %}
Verify your kubeconfig's current context.
{% capture code_content %}kubectl config current-context{% endcapture %}
{% include code.html code=code_content %}
To switch your kubeconfig's current context.
{% capture code_content %}kubectl config use-context <rke2 context>{% endcapture %}
{% include code.html code=code_content %}

### Preflight: Plan your access to Meshery UI

1. If you are using port-forwarding, please refer to the [port-forwarding]({{ site.baseurl }}/reference/mesheryctl/system/dashboard) guide for detailed instructions.
2. Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

Meshery should now be running in your RKE2 cluster and Meshery UI should be accessible at the `INTERNAL IP` of `meshery` service.

## Installation: Using `mesheryctl`

Use Meshery's CLI to streamline your connection to your RKE2 cluster. Configure Meshery to connect to your RKE2 cluster by executing:

{% capture code_content %}$ mesheryctl system config <rke2 context>{% endcapture %}
{% include code.html code=code_content %}

Once configured, execute the following command to start Meshery.

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

If you encounter any authentication issues, you can use `mesheryctl system login`. For more information, click [here](/guides/mesheryctl/authenticate-with-meshery-via-cli) to learn more.

## Installation: Using Helm

For detailed instructions on installing Meshery using Helm V3, please refer to the [Helm Installation](/installation/kubernetes/helm) guide.

## Installation: Manual Steps

You may also manually generate and load the kubeconfig file for Meshery to use:

**The following configuration yaml will be used by Meshery. Copy and paste the following in your config file** :

{% capture code_content %}apiVersion: v1
clusters:

- cluster:
  certificate-authority-data: <rke2-cert-shortcutted>
  server: https://<rke2-server-ip>:<rke2-server-port>
  name: rke2
  contexts:
- context:
  cluster: rke2
  user: rke2-user
  name: rke2
  current-context: rke2-context
  kind: Config
  preferences: {}
  users:
- name: rke2-user
  user:
  client-certificate-data: <rke2-cert-shortcutted>
  client-key-data: <rke2-key-shortcutted>{% endcapture %}
  {% include code.html code=code_content %}

_Note_: Make sure _current-context_ is set to _rke2-context_.

<br />
**To allow Meshery to auto detect your config file, Run** :
{% capture code_content %}kubectl config view --minify --flatten > config_minikube.yaml{% endcapture %}
{% include code.html code=code_content %}

<br />
Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.




## Installation: Upload Config File in Meshery Web UI

- Run the below command to generate the _"config_rke2.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">kubectl config view --minify --flatten > config_rke2.yaml</div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option.

## Post-Installation Steps

Optionally, you can verify the health of your Meshery deployment, using <a href='/reference/mesheryctl/system/check'>mesheryctl system check</a>.

You're ready to use Meshery! Open your browser and navigate to the Meshery UI.

{% include_cached installation/accessing-meshery-ui.md %}

{% include related-discussions.html tag="meshery" %}