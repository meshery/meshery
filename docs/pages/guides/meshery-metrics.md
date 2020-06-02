---
layout: guide
title: Using Metrics in Meshery
description: How to connected and use Prometheus and Grafana metrics in Meshery
permalink: guides/meshery-metrics
type: guide
---

In order to generate report/load testing of your k8s applications, Meshery uses `Grafana` and `Prometheus`.
Here you can find steps to generate these endpoints in different environments or for different service mesh adapters.

* TOC
{:toc}


### Istio + Minikube(dev env)



In order for Meshery to gain access to Grafana and Prometheus, their network ports need to be exposed.



**1) Get minikube ip**

Retreive the IP address of your Minikube cluster by executing:
    
```        
    $minikube ip
    o/p: 172.17.0.2
```

<i>Note: Istio is installed in `istio-system` namespace and the "BookInfo" sample app is installed in `book-info` namespace.
</i>

**2) Expose Prometheus service**

a) By default  `prometheus`  spec type is configured to `ClusterIP`. You can change it to `NodePort` by executing:

```
$kubectl patch svc prometheus -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

b) Get NodePort of `prometheus` service by executing:

```
$kubectl describe services prometheus -n istio-system|grep NodePort

o/p:NodePort:  http  30535/TCP
```

c) Prometheus endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30822
```

**3) Expose Grafana service**

a) By default  `Grafana`  spec type is configured to `ClusterIP`  you can change it to `NodePort` using below command.

```
$ kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

b) Get NodePort of `Grafana` service using below command

```
$kubectl describe services grafana -n istio-system|grep NodePort

o/p:NodePort:  http  32130/TCP
```

c) Grafana endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:32130
```

**4) Expose Istio BookInfo sample app `productpage` service**

a) By default  `productpage`  spec type is configured to `ClusterIP`  you can change it to `NodePort` using below command.

```
$ kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n book-info
```

b) Get NodePort of `productpage` service using below command

```
$ kubectl describe services productpage -n book-info|grep NodePort

o/p:NodePort:  http  30535/TCP
```

c) `productpage` endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30535/productpage
```

### Linkerd

**To be Updated**
