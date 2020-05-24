---
layout: guide
title: meshery metrics
description: How to use meshery metrics
permalink: guides/mesherymetrics
type: guide
---

In order to generate report/load testing of your k8s applications, Meshery uses `Grafana` and `Prometheus`.
Here you can find steps to generate these endpoints in different environemnts/service mesh adaptors.

* TOC
{:toc}


### Istio + Minikube(dev env)



In order for Meshery to gain access to grafana and prometheus, their network ports need to be exposed.



**1) Get minikube ip**

Get Minikube ip by executed below command.
    
```        
    $minikube ip
    o/p: 172.17.0.2
```
<i>Note: Here Istio is installed in `istio-system` name space and `productpage` sample app is installed in `book-info` name space.
</i>

**2)Expose Prometheus service**

a)By default  `prometheus`  spec type is configured to `ClusterIP`  you can chage it to `NodePort` using below command.

```
$kubectl patch svc prometheus -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

b)Get NodePort of `prometheus` service using below command

```
$kubectl describe services prometheus -n istio-system|grep NodePort

o/p:NodePort:  http  30535/TCP
```

c) Prometheus endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30822
```


**3)Expose Grafana service**

a)By default  `Grafana`  spec type is configured to `ClusterIP`  you can chage it to `NodePort` using below command.

```
$kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

b)Get NodePort of `Grafana` service using below command

```
$kubectl describe services grafana -n istio-system|grep NodePort

o/p:NodePort:  http  32130/TCP
```

c) Grafana endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:32130
```


**4)Expose  istio sample app service productpage**

a)By default  `productpage`  spec type is configured to `ClusterIP`  you can change it to `NodePort` using below command.

```
$kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n book-info
```

b)Get NodePort of `productpage` service using below command

```
$kubectl describe services productpage -n book-info|grep NodePort

o/p:NodePort:  http  30535/TCP
```

c) productpage endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30535
```






### Linkerd

**To be Updated**
