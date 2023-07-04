---
layout: default
title: Using Metrics in Meshery
description: How to connect and use Prometheus and Grafana metrics in Meshery
permalink: guides/meshery-metrics
language: en
type: Guides
category: performance
---

## Connect and use metrics in Meshery

Meshery provides performance reports, including performance test results, node resource metrics etc. so that operators may easily understand the overhead of their service mesh's control plane and data plane in context of the overhead incurred on nodes running within the cluster. In order to generate performance test reports of service meshes and their workloads, Meshery uses Grafana and/or Prometheus as visualization and metrics systems, respectively. This guide outlines the requirements necessary for Meshery to connect to these systems. The steps may vary depending upon the service mesh and its configuration.

In order to pull in these environment metrics, you can also manually configure Meshery to connect with your existing Grafana and/or Prometheus instances through the Meshery dashboard. Once they have been loaded and are displayed on the screen, you may also perform an *ad-hoc* test to check Meshery's connection status.

### Prometheus Charts

User needs to set the Prometheus URL and API key to create and query boards.

[![Prometheus Charts]({{ site.baseurl }}/assets/img/architecture/PrometheusCharts.svg)]({{ site.baseurl }}/assets/img/architecture/PrometheusCharts.svg)


### Grafana Charts

User needs to set the Grafana URL and API key to create and query boards.

[![Grafana Charts]({{ site.baseurl }}/assets/img/architecture/GrafanaBoards.svg)]({{ site.baseurl }}/assets/img/architecture/GrafanaBoards.svg)

### Dynamic Boards

Dynamic Boards can be generated from Prometheus or Grafana. These boards are defined by the user. Grafana SDK is used for these boards.

### Static Boards

Static Boards capture the Service Mesh Performance. Certain protos are defined and tracked as a part of performance tests. Static Boards Queries Prometheus SDK directly.



<!-- ## Tutorial Guide

Connect Meshery to your Grafana and Prometheus instances to enable enhanced service mesh performance management. Deploy a service mesh and any available sample application

Retreive the IP address of your Minikube cluster by executing:
    
 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">       
 minikube ip
 172.17.0.2
 </div></div>
 </pre>

* Set up Grafana and/or Prometheus:
1. [Connect Meshery to metric systems](#connect-meshery-to-metric-systems)
2. [Connect Meshery to metric systems using Minikube](#Connect-Meshery-to-metric-systems-using-Minikube)

* Expose the service metric ports - The service ports of Grafana and Prometheus need to be exposed in order for Meshery to connect to and interact with these visualizations and metrics systems.
* Access the port assigned to the metric service
* [Run Performance tests](#run-performance-tests)


### **Connect Meshery to metric systems**

#### 1. Using kubectl, edit the Grafana and Prometheus services in the *Istio-system* namespace:

```sh
kubectl edit svc grafana -n istio-system
```

#### 2. Change specification type

By default, the service specification types, like Grafana is configured to `ClusterIP`. You can change it to `NodePort` by executing:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 kubectl patch svc prometheus -p '{"spec": {"type": "NodePort"}}' -n istio-system
 </div></div>
 </pre>

#### 3. Get the newly assigned port for your chosen service

Run:
```
kubectl get svc grafana -n istio-system
```

Example output:
```
NAME      TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
grafana   NodePort   10.100.67.144   <none>        3000:30188/TCP   3d11h
```

#### 4. Access the service port

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 http://172.17.0.2:30822
 </div></div>
 </pre>

###### 2. Expose the service

By default, the service specification types, like Prometheus, Grafana, and the `productpage`, are configured to **ClusterIP**. You can change it to **NodePort** by replacing *service spec type* with the spec you wish to run and executing:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n istio-system
 </div></div>
 </pre>

### **Expose Grafana service**

* Get the NodePort of Grafana service using below command

```
kubectl describe services grafana -n istio-system|grep NodePort

o/p:NodePort:  http  32130/TCP
```

* Find the Grafana endpoint

The Grafana endpoint will be *http://$MINIKUBE_IP:NODE_PORT*

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 http://172.17.0.2:32130
 </div></div>
 </pre>

<a href="{{ site.baseurl }}/assets/img/meshery-metrics/grafana-server-settings.png">
  <img style="width:500px;" src="{{ site.baseurl }}/assets/img/meshery-metrics/grafana-server-settings.png" />
</a>

### **Expose Prometheus service**

Meshery allows you to expose Prometheus as a service with a single click. You can do this:
- [Through the Meshery UI](#meshery-ui)
- Alternatively, you can also attempt [Manual Integration](#manual-steps)

#### **Meshery UI**

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

#### **Manual Steps**

* Get the NodePort of Prometheus service by executing:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 kubectl patch svc grafana -p '{"spec": {"type": "NodePort"}}' -n book-info
 </div></div>
 </pre>

* Find the Prometheus endpoint 

The Prometheus endpoint will be *http://$MINIKUBE_IP:NODE_PORT*

```
http://172.17.0.2:30822
```

#### **Expose Istio BookInfo sample app `productpage` service**

* Get the NodePort of `productpage` service by executing:

```
kubectl describe services productpage -n book-info|grep NodePort

NodePort:  http  30535/TCP
```

* Find the `productpage` endpoint

The `productpage` endpoint will be http://$MINIKUBE_IP:NODE_PORT

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 http://172.17.0.2:30535/productpage
 </div></div>
 </pre>

#### **Run Performance Tests**

After successfully setting up a connection between your metric service and Meshery, you may proceed to run performance tests by navigating to the Performance Test tab on Meshery:

<a href="{{ site.baseurl }}/assets/img/performance-management/performance-meshery.png"><img style="width:450px;padding-top:5px;" src="{{ site.baseurl }}/assets/img/performance-management/performance-meshery.png" /></a>

**Run Test Results**

<a href="{{ site.baseurl }}/assets/img/performance-management/meshery-and-grafana.png"><img style="width:450px;padding-top:5px;" src="{{ site.baseurl }}/assets/img/performance-management/meshery-and-grafana.png" /></a> -->

##### Suggested Reading

- Guide: [Interpreting Performance Test Results]({{ site.baseurl }}/guides/interpreting-performance-test-results)
