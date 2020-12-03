---
layout: default
title: Kubernetes
permalink: installation/platforms/kubernetes
type: installation
language: en
list: include
image: /docs/assets/img/platforms/kubernetes.svg
---

{% include installation_prerequisites.html %}

**To set up and run Meshery on Kubernetes** 

- [ Use Helm and set up a Kubernetes cluster](#using-helm)
- [Run Meshery on existing Kubernetes cluster](#using-kubernetes-manifests)

### **Using Helm**

#### 1. **Helm v3**
Run the following:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ git clone https://github.com/layer5io/meshery.git; cd meshery
 $ kubectl create namespace meshery
 $ helm install meshery --namespace meshery install/kubernetes/helm/meshery
 </div></div>
 </pre>

#### 2. **Helm v2**
Run the following:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ git clone https://github.com/layer5io/meshery.git; cd meshery
 $ kubectl create namespace meshery
 $ helm template meshery --namespace meshery install/kubernetes/helm/meshery | kubectl apply -f -
 </div></div>
 </pre>

### **Using Kubernetes Manifests**
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
  $ kubectl create ns meshery```
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
         serviceName: meshery-service
         servicePort: 8080
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

Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides](/docs/guides) for advanced usage tips.
