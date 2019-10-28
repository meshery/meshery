---
layout: page
title: Kubernetes
permalink: installation/kubernetes
---

# Quick Start with Kubernetes
Meshery can also be deployed on an existing Kubernetes cluster. See [compatibility table](#compatibility-matrix) for version compatibility. To install Meshery on your cluster, let us first create a namespace to host the Meshery components:
```
kubectl create ns meshery
```

All the needed deployment yamls for deploying Meshery are included in the `deployment_yamls/k8s` folder inside the cloned Meshery folder. To deploy the yamls on the cluster please run the following command:
```
kubectl -n meshery apply -f deployment_yamls/k8s
```
Once the yaml files are deployed, we need to expose the `meshery` service to be able to access the service from outside the cluster. 

There are several ways a service can be exposed on Kubernetes. 

Here we will describe 3 common ways we can expose a service:

* **Ingress** If your Kubernetes cluster has a functional Ingress Controller, then you can configure an ingress to 

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

* **LoadBalancer**
    If your Kubernetes cluster has an external load balancer, this might be a logical route.

* **NodePort**
    If your cluster does not have an Ingress Controller or a load balancer, then NodePort is probably the last resort.