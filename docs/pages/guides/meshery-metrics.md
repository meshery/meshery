---
layout: default
title: Using Metrics in Meshery
description: How to connect and use Prometheus and Grafana metrics in Meshery
permalink: guides/meshery-metrics
type: Guides
---


### **Connect and use metrics in Meshery**

Meshery provides performance reports, including performance test results, node resource metrics etc. so that operators may easily understand the overhead of their service mesh's control plane and data plane in context of the overhead incurred on nodes running within the cluster. In order to generate performance test reports of service meshes and their workloads, Meshery uses `Grafana` and/or `Prometheus` as visualization and metrics systems, respectively. This guide outlines the requirements necessary for Meshery to connect to these systems. The steps may vary depending upon the service mesh and its configuration.

In order to pull in these environment metrics, you can also manually configure Meshery to connect with your existing Grafana and/or Prometheus instances through the Meshery dashboard. Once they have been loaded and are displayed on the screen, you may also perform an *ad-hoc* test to check Meshery's connection status.

#### Tutorial Guide

Follow along with this tutorial to set up and integrate Grafana and Prometheus instances on top of your service mesh and run performance tests:

* Deploy a service mesh and any available sample application

  For the purpose of this tutorial guide, we will be installing [Istio](/docs/service-meshesadapters/istio) and deploying the [BookInfo sample application](/docs/guides/sample-apps#bookinfo):

    - [Install Istio](/docs/service-meshesadapters/istio) on Meshery in the `istio-system` namespace. 
      <a href="/docs/assets/img/adapters/istio/istio-install.png">
      <img style="width:300px;padding-top:5px;" src="/docs/assets/img/adapters/istio/istio-install.png" />
      </a>
    - [Deploy the BookInfo sample application](/docs/guides/sample-apps#to-deploy-a-sample-app-on-meshery) 

      The [BookInfo](/docs/guides/sample-apps#bookinfo) sample app should be installed in the `default` namespace unless otherwise specified upon deployment of BookInfo.

* Set up Grafana and/or Prometheus:
1. [Connect Meshery to metric systems](#connect-meshery-to-metric-systems)
2. [Connect Meshery to metric systems using Minikube](#Connect-Meshery-to-metric-systems-using-Minikube)

* Expose the service metric ports - The service ports of Grafana and Prometheus need to be exposed in order for Meshery to connect to and interact with these visualizations and metrics systems.
* Access the port assigned to the metric service
* [Run Performance tests](#run-performance-tests)


#### **Connect Meshery to metric systems**

###### 1. Using kubectl, edit the Grafana and Prometheus services in the *Istio-system* namespace:

```sh
kubectl edit svc grafana -n istio-system
```

###### 2. Change specification type

By default, the service specification types, like `grafana` is configured to `ClusterIP`. You can change it to `NodePort` by executing:

```
$ kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n istio-system
```

###### 3. Get the newly assigned port for your chosen service

Run:
```
$ kubectl get svc grafana -n istio-system
```

Example output:
```
NAME      TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
grafana   NodePort   10.100.67.144   <none>        3000:30188/TCP   3d11h
```

###### 4. Access the service port

 Using your host's IP address, enter http://<host-ip>:30188 into the Metrics panel in Meshery's Settings page (replace 30188 with your environment's port number).

#### **Connect Meshery to metric systems using Minikube**

In this scenario, Meshery will be configured to deploy Istio to a Kubernetes cluster running in Minikube. The service ports of Grafana and Prometheus need to be exposed in order for Meshery to connect to and interact with these visualizations and metrics systems.

- [Grafana spec](#expose-grafana-service)
- [Prometheus spec](#expose-prometheus-service)
- [Bookinfo - `productpage` service](#expose-istio-bookinfo-sample-app-productpage-service)

###### 1. Get minikube ip

Retreive the IP address of your Minikube cluster by executing:
    
```        
    $ minikube ip
    172.17.0.2
```

###### 2. Expose the service


By default, the service specification types, like `prometheus`, `grafana`, and the `productpage`, are configured to **ClusterIP**. You can change it to **NodePort** by replacing *service spec type* with the spec you wish to run and executing:

```
$ kubectl patch svc `service spec type` -p '{"spec": {"type": "NodePort"}}' -n istio-system
```


###### **Expose `Grafana` service**

* Get the NodePort of `grafana` service using below command

```
$kubectl describe services grafana -n istio-system|grep NodePort

o/p:NodePort:  http  32130/TCP
```

* Find the `grafana` endpoint

The Grafana endpoint will be *http://$MINIKUBE_IP:NODE_PORT*

```
http://172.17.0.2:32130
```

<a href="{{ site.baseurl }}/assets/img/meshery-metrics/grafana-server-settings.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/meshery-metrics/grafana-server-settings.png" />
</a>

###### **Expose `Prometheus` service**

Meshery allows you to expose Prometheus as a service with a single click. You can do this:
- [Through the Meshery UI](#meshery-ui)
- Alternatively, you can also attempt [Manual Integration](#manual-steps)

###### **Meshery UI**

Meshery auto-discovers all Prometheus instances available on your local system and will offer you a list of options to choose from. You can select the Prometheus Server that you wish to employ.


<a href="{{ site.baseurl }}/assets/img/meshery-metrics/prometheus-settings.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/meshery-metrics/prometheus-settings.png" />
</a>

Meshery also provides you the option of simply pasting in your Prometheus queries:

- Navigate to the management page for Istio on the Meshery UI
- Click on <i class="fas fa-caret-right fa-lg"></i>, located under **Apply Custom Configuration** 
- Paste in your Prometherus query. Click on <i class="fas fa-caret-right fa-lg"></i>

<a href="{{ site.baseurl }}/assets/img/meshery-metrics/prometheus-query.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/meshery-metrics/prometheus-query.png" />
</a>

###### **Manual Steps**

* Get the NodePort of `prometheus` service by executing:

```
$ kubectl describe services prometheus -n istio-system|grep NodePort

NodePort:  http  30535/TCP
```

* Find the `prometheus` endpoint 

The Prometheus endpoint will be *http://$MINIKUBE_IP:NODE_PORT*

```
http://172.17.0.2:30822
```


###### **Expose Istio BookInfo sample app `productpage` service**

* Get the NodePort of `productpage` service by executing:

```
$ kubectl describe services productpage -n book-info|grep NodePort

NodePort:  http  30535/TCP
```

* Find the `productpage` endpoint

The `productpage` endpoint will be http://$MINIKUBE_IP:NODE_PORT

```
http://172.17.0.2:30535/productpage
```

###### **Run Performance Tests**

After successfully setting up a connection between your metric service and Meshery, you may proceed to run performance tests by navigating to the Performance Test tab on Meshery:

<a href="/docs/assets/img/performance-management/performance-meshery.png"><img style="width:450px;padding-top:5px;" src="/docs/assets/img/performance-management/performance-meshery.png" /></a>

**Run Test Results**

<a href="/docs/assets/img/performance-management/meshery-and-grafana.png"><img style="width:450px;padding-top:5px;" src="/docs/assets/img/performance-management/meshery-and-grafana.png" /></a>

##### Suggested Reading

- Guide: [Interpreting Performance Test Results](/docs/guides/interpreting-performance-test-results)
