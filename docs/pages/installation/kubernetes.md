---
layout: default
title: Kubernetes
permalink: installation/kubernetes
type: installation
category: kubernetes
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kubernetes.svg
---

{% include installation_prerequisites.html %}

## Available Deployment Methods

- [Using `mesheryctl`](#using-mesheryctl)
- [Using `helm`](#using-helm-charts)
<!-- - [Using Kubernetes manifests](#using-kubernetes-manifests) -->

### Using mesheryctl

<pre class="codeblock-pre">
<div class="codeblock">
 <div class="clipboardjs">
 $ mesheryctl system context create k8s -p kubernetes -s
 </div></div>
</pre>

<i>Don't have `mesheryctl`? <a href="{{site.baseurl}}/installation/mesheryctl">Install with Bash, Brew, or Scoop</a></i>.

<details>
<summary>Optional: Verify Step</summary>

<p>
Ensure that your <code>current-context</code> has <code>platform: kubernetes</code> configured in <code>~/.meshery/config.yaml</code>. Example context:</p>

<pre>
<code>
➜  ~ mesheryctl system context view
endpoint: http://localhost:9081
token: default
platform: kubernetes
adapters:
- meshery-istio
- meshery-linkerd
- meshery-consul
- meshery-nsm
- meshery-kuma
- meshery-cpx
- meshery-osm
- meshery-traefik-mesh
channel: stable
version: latest
</code>
</pre>

</details>

Deploy Meshery to your Kubernetes cluster by executing:

{% capture code_content %}$ mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}


## Using Helm

<pre class="codeblock-pre">
<div class="codeblock">
 <div class="clipboardjs">
 $ kubectl create ns meshery
 $ helm repo add meshery https://meshery.io/charts/
 $ helm install meshery meshery/meshery -n meshery
 </div></div>
</pre>

### Using Helm

{% capture code_content %}$ helm repo add meshery https://meshery.io/charts/
$ helm install meshery-operator meshery/meshery-operator
$ helm install meshery meshery/meshery{% endcapture %}
{% include code.html code=code_content %}

Customize of deployment the Meshery adapters:

###### Example: Pin your deployment to a specific Meshery version
<pre>
<code>$ helm install meshery meshery/meshery --version v0.7.0</code>
</pre>


###### Example: Disabled the Meshery Adapter for Linkerd and verify the deployment manifest
<pre>
<code>$ helm install --set meshery-linkerd.enabled=false meshery/meshery --dry-run</code>
</pre>

The key for Meshery Adapters can be found [here](https://artifacthub.io/packages/helm/meshery/meshery#values)


<!-- ### **[deprecated] Using Kubernetes Manifests **
Meshery can also be deployed on an existing Kubernetes cluster. See [compatibility table](#compatibility-matrix) for version compatibility. To install Meshery on your cluster, clone the Meshery repo:

{% capture code_content %} $ git clone https://github.com/layer5io/meshery.git;
 $ cd meshery{% endcapture %}
{% include code.html code=code_content %}

Create a namespace as a new logical space to host Meshery and its components:

{% capture code_content %}$ kubectl create ns meshery{% endcapture %}
{% include code.html code=code_content %}

All the needed deployment yamls for deploying Meshery are included in the *install/deployment_yamls/k8s* folder inside the cloned Meshery folder. To deploy the yamls on the cluster please run the following command:

{% capture code_content %}$ kubectl -n meshery apply -f install/deployment_yamls/k8s{% endcapture %}
{% include code.html code=code_content %}

Once the yaml files are deployed, we need to expose the *meshery* service to be able to access the service from outside the cluster. There are several ways a service can be exposed on Kubernetes. Here we will describe 3 common ways we can expose a service: -->

#### Exposing Meshery Service

If your Kubernetes cluster has a functional Ingress Controller, then you can configure an ingress to expose Meshery:

* **Ingress (example)** 

{% capture code_content %} apiVersion: extensions/v1beta1
 kind: Ingress
 metadata:
 name: meshery-ingress
 annotations:
     kubernetes.io/ingress.class: "nginx"
 spec:
 rules:
 - host: *
     http:
     paths:
     - path: /
         backend:
         serviceName: meshery-service
         servicePort: 9081{% endcapture %}
{% include code.html code=code_content %}

* **LoadBalancer** - If your Kubernetes cluster has support of a LoadBalancer resource, this is an ideal choice.

* **NodePort** - If your cluster does not have an Ingress Controller or a LoadBalancer support, then use NodePort to expose Meshery:

{% capture code_content %} apiVersion: v1
 kind: Service
 spec:
     type: NodePort{% endcapture %}


### Sample Commands to verify Meshery UI access:

{% capture code_content %}export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")

kubectl --namespace meshery port-forward $POD_NAME 9081:8080{% endcapture %}
{% include code.html code=code_content %}

# Advanced Installations

**Configurable OAuth Callback URL** ([learn more]({{site.baseurl}}/extensibility/providers#configurable-oauth-callback-url))

Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

<pre><code>
$ helm install meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host meshery/meshery
</code></pre>


Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}
