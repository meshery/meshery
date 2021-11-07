---
layout: default
title: Kubernetes
permalink: installation/platforms/kubernetes
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kubernetes.svg
---

{% include installation_prerequisites.html %}

## Available Deployment Methods

- [Using `mesheryctl`](#using-mesheryctl)
- [Using `helm`](#using-helm-charts)
- [Using Kubernetes manifests](#using-kubernetes-manifests)

### **Using mesheryctl**
Ensure that your `current-context` has `platform: kubernetes` configured in `~/.meshery/config.yaml`. Example context:

```
➜  ~ mesheryctl system context view
endpoint: http://localhost:9081
token: Default
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
```

Deploy Meshery to your Kubernetes cluster by executing:

<pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
    $ mesheryctl system start
 </div></div>
</pre>

Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
<pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
    $ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start
 </div></div>
</pre>

### **Using Helm Charts**

Run the following for default:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ helm repo add meshery https://meshery.io/charts/
 $ helm install meshery meshery/meshery
 </div></div>
 </pre>

Customize of deployment the Meshery adapters:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ helm repo add meshery https://meshery.io/charts/

 # Example: Pin your deployment to a specific Meshery version
 $ helm install meshery meshery/meshery --version v0.5.67

 # Example: Disabled the Meshery Adapter for Linkerd and verify the deployment manifest
 $ helm install --set meshery-linkerd.enabled=false meshery/meshery --dry-run

 # Example: Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way
 $ helm install meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host meshery/meshery
 </div></div>
 </pre>

The key of Meshery adapters you can find [here](https://artifacthub.io/packages/helm/meshery/meshery#values)

### **Using Kubernetes Manifests [deprecated]**
Meshery can also be deployed on an existing Kubernetes cluster. See [compatibility table](#compatibility-matrix) for version compatibility. To install Meshery on your cluster, clone the Meshery repo:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ git clone https://github.com/layer5io/meshery.git;
 $ cd meshery
 </div></div>
 </pre>

Create a namespace as a new logical space to host Meshery and its components:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ kubectl create ns meshery
 </div></div>
 </pre>

All the needed deployment yamls for deploying Meshery are included in the *install/deployment_yamls/k8s* folder inside the cloned Meshery folder. To deploy the yamls on the cluster please run the following command:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ kubectl -n meshery apply -f install/deployment_yamls/k8s
 </div></div>
 </pre>

Once the yaml files are deployed, we need to expose the *meshery* service to be able to access the service from outside the cluster. There are several ways a service can be exposed on Kubernetes. Here we will describe 3 common ways we can expose a service:

#### **Ingress**

If your Kubernetes cluster has a functional Ingress Controller, then you can configure an ingress to expose Meshery:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 apiVersion: extensions/v1beta1
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
        # Please kindly check your service name and service port to confirm the Ingress can work well
         serviceName: meshery-service
         servicePort: 9081
 </div></div>
 </pre>

* **LoadBalancer** - If your Kubernetes cluster has an external load balancer, this might be a logical route.

* **NodePort** - If your cluster does not have an Ingress Controller or a load balancer, then use NodePort to expose Meshery:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 apiVersion: v1
 kind: Service
 spec:
     type: NodePort
 </div></div>
 </pre>

Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.
