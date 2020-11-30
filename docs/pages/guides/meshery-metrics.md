---
layout: default
title: Using Metrics in Meshery
description: How to connected and use Prometheus and*Grafana*metrics in Meshery
permalink: guides/meshery-metrics
type: Guides
---

In order to generate performance test reports of service meshes and their workloads, Meshery uses **Grafana** and/or **Prometheus** as visualization and metrics systems, respectively. This document outlines the requirements necessary for Meshery to connect to these systems. Steps vary depending upon the service mesh and its configuration.

* TOC
{:toc}


## Istio + Minikube

In this scenario, Meshery has deployed Istio to a Kubernetes cluster running in Minikube. The service ports of*Grafana*and Prometheus need to be exposed in order for Meshery to connect to and interact with these visualization and metrics systems.

### Get minikube ip

Retreive the IP address of your Minikube cluster by executing:
    
```        
    $ minikube ip
    172.17.0.2
```

<i>Note: Istio is installed in `istio-system` namespace and the "BookInfo" sample app is installed in `default` namespace unless otherwise specified upon deployment of "BookInfo".
</i>

### Expose Prometheus service

a) By default  *prometheus* service spec type is configured to *ClusterIP*. You can change it to **NodePort** by executing:

```
$ kubectl patch svc prometheus -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

b) Get NodePort of *prometheus* service by executing:

```
$ kubectl describe services prometheus -n istio-system|grep NodePort

NodePort:  http  30535/TCP
```

c) Prometheus endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30822
```

### Expose Grafana service

a) By default  *Grafana*  spec type is configured to *ClusterIP*  you can change it to **NodePort** by executing:

```
$ kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

b) Get NodePort of *Grafana* service using below command

```
$kubectl describe services grafana -n istio-system|grep NodePort

o/p:NodePort:  http  32130/TCP
```

c) Grafana endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:32130
```

### Expose Istio BookInfo sample app `productpage` service

a) By default *productpage*  spec type is configured to *ClusterIP*  you can change it to **NodePort** by executing:

```
$ kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n book-info
```

b) Get NodePort of *productpage* service by executing:

```
$ kubectl describe services productpage -n book-info|grep NodePort

NodePort:  http  30535/TCP
```

c) *productpage* endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30535/productpage
```

## Linkerd

_coming soon..._
