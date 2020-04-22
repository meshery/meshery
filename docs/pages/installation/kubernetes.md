---
layout: page
title: Kubernetes
permalink: installation/kubernetes
---

# Quick Start with Kubernetes
  
<a name="helm"></a>
## Using Helm

### Helm v3
Run the following:
```
$ git clone https://github.com/layer5io/meshery.git; cd meshery
$ kubectl create namespace meshery
$ helm install meshery --namespace meshery install/kubernetes/helm/meshery
```

### Helm v2
Run the following:
 ```
 $ git clone https://github.com/layer5io/meshery.git; cd meshery
 $ kubectl create namespace meshery
 $ helm template meshery --namespace meshery install/kubernetes/helm/meshery
 ```
## Using Kubernetes Manifests
Meshery can also be deployed on an existing Kubernetes cluster. See [compatibility table](#compatibility-matrix) for version compatibility. To install Meshery on your cluster, clone the Meshery repo:


```
$ git clone https://github.com/layer5io/meshery.git; cd meshery
```

Create a namespace as a new logical space to host Meshery and its components:

```
$ kubectl create ns meshery
```

All the needed deployment yamls for deploying Meshery are included in the `install/deployment_yamls/k8s` folder inside the cloned Meshery folder. To deploy the yamls on the cluster please run the following command:

```
$ kubectl -n meshery apply -f install/deployment_yamls/k8s
```

Once the yaml files are deployed, we need to expose the `meshery` service to be able to access the service from outside the cluster. There are several ways a service can be exposed on Kubernetes. Here we will describe 3 common ways we can expose a service:

* **Ingress** - If your Kubernetes cluster has a functional Ingress Controller, then you can configure an ingress to expose Meshery: 

```
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
```
* **LoadBalancer** - If your Kubernetes cluster has an external load balancer, this might be a logical route.
* **NodePort** - If your cluster does not have an Ingress Controller or a load balancer, then use NodePort to expose Meshery:
```
apiVersion: v1
kind: Service
spec:
    type: NodePort
```