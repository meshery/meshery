---
layout: default
title: Minikube
permalink: installation/platforms/minikube
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/minikube.png
---

{% include installation_prerequisites.html %}

**To Setup and run Meshery on Minikube** :

- [Steps](#steps)
  - [1. Start minikube](#1-start-minikube)
  - [2. Install Meshery](#2-install-meshery)
  - [2. Configure Meshery to use minikube](#2-configure-meshery-to-use-minikube)
  - [Manual Steps](#manual-steps)

##### Compatibility

The following minimum component versions are required:

<table id="compatibility-table">
  <tr>
    <th id="model">Name</th>
    <th id="model">Version</th>
  </tr>
  <tr>
    <td><a href="https://kubernetes.io/docs/tasks/tools/install-minikube/">Minikube</a></td>
    <td>1.0.0 </td>
  </tr>
  <tr>
    <td><a href="https://istio.io/docs/setup/kubernetes/prepare/platform-setup/minikube/">Kubernetes</a></td>
    <td>1.14.1</td>
  </tr>
  <tr>
    <td><a href="https://kubernetes.io/docs/tasks/tools/install-kubectl/">kubectl</a></td>
    <td>1.14.1</td>
  </tr>
</table>

## Steps

Perform the following steps in order:

### 1. Start minikube

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1
 </div></div>
 </pre>

_Note: minimum memory required is --memory=4096 (for Istio deployments only)_

**Check up on your minikube cluster** :

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
  minikube status 
 </div></div>
 </pre>

### 2. Install Meshery

Follow the [installation steps](/installation/quick-start) to setup the mesheryctl CLI and install Meshery.

**Users using docker driver**:
After completing the Meshery installation, execute the following commands to establish connectivity between Meshery Server and Kubernetes cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 docker network connect bridge meshery_meshery_1
 </div></div>
 </pre>

<br/>

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 docker network connect minikube meshery_meshery_1
 </div></div>
 </pre>

To establish connectivity between a particular Meshery Adapter and Kubernetes server, use _"docker ps"_ to identify the name of the desired container, and execute the following commands:

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 docker network connect bridge &#60; container name of the desired adapter &#62;
 </div></div>
 </pre>

<br/>

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 docker network connect minikube &#60; container name of the desired adapter &#62;
 </div></div>
 </pre>

### 3. Configure Meshery to use minikube

1. Login to Meshery. Under your user profile, click _Get Token_.

2. Use [mesheryctl]({{ site.baseurl }}/installation#using-mesheryctl) to configure Meshery to use minikube. To allow Meshery to detect your config file, execute the following commands:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system config minikube -t ~/Downloads/auth.json
 </div></div>
 </pre>
<br/>

**Optionally run the command below to expose the LoadBalancer services to the host machine** :

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 minikube tunnel
 </div></div>
 </pre>

**Optionally configure Meshery to use minikube through the Web UI** :

- Run the below command to generate the _"config_minikube.yaml"_ file for your cluster:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 kubectl config view --minify --flatten > config_minikube.yaml
 </div></div>
 </pre>

- Upload the generated config file by navigating to _Settings > Environment > Out of Cluster Deployment_ in the Web UI and using the _"Upload kubeconfig"_ option.

### Manual Steps

You may also manually generate and load the kubeconfig file for Meshery to use:

**The following configuration yaml will be used by Meshery. Copy and paste the following in your config file** :

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 apiVersion: v1
 clusters:
 - cluster:
     certificate-authority-data: < cert shortcutted >
     server: https://192.168.99.100:8443
   name: minikube
 contexts:
 - context:
     cluster: minikube
     user: minikube
   name: minikube
 current-context: minikube
 kind: Config
 preferences: {}
 users:
 - name: minikube
   user:
     client-certificate-data: < cert shortcutted >
     client-key-data: < key shortcutted >
 </div></div>
 </pre>

_Note_: Make sure _current-context_ is set to _minikube_.

<br />
**To allow Meshery to auto detect your config file, Run** :
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 kubectl config view --minify --flatten > config_minikube.yaml
 </div></div>
</pre>

<br />
Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.
