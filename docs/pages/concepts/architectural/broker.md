---
layout: default
title: Broker
permalink: concepts/architecture/broker
type: components
redirect_from: architecture/broker
abstract: "Meshery broker component facilitates data streaming between kubernetes cluster components and outside world."
language: en
list: include
---

Broker is a custom Kubernetes controller that provides data streaming across independent components of Meshery whether those components are running inside or outside of the Kubernetes cluster.

[![Meshery Log Viewer]({{ site.baseurl }}/assets/img/architecture/meshery-log-viewer.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-log-viewer.svg)

### Broker FAQs

#### How many Brokers can run?
It is recommended to run one broker instance for each kubernetes cluster, However the instance itself can be scaled up based on the incoming data volume in each of the cluster. The scaling is independent of the number of instances running.

#### What does an HA configuration look like?
We leverage on the kubernetes functionality in terms of the High-Availability behaviour. Meaning, the broker instance gets instantiated/restarted on its own when an issue occurs. In part, Meshery-Operator is also responsible for keeping the broker functional.

#### What stateful characteristics does the Broker have?
All the messages that are published to the broker is persisted in-memory within the broker instance until it get consumed. Persistent-volume/Disk-space is not currently being used by the Broker.

#### How do I know if the Broker is working? How do I troubleshoot the Broker?
The Broker instance is deployed inside the kubernetes cluster as a `Statefulset`. In the case where the broker doesn't seem to work, here are a few steps to troubleshoot the instance:

- Make sure the pods corresponding to the `Statefulset` is up and running.
- Make sure the kubernetes cluster has support for kubernetes `Service` type `LoadBalancer` or `NodePort`.
- Ensure connectivity between the Meshery-Server and the Broker service endpoint.

#### Where can I find help troubleshooting Meshery Operator, MeshSync, and Broker?
If you're experiencing issues with **Meshery Operator, MeshSync, or Broker**, the best place to find solutions is the **[Meshery Troubleshooting Guide](https://docs.meshery.io/guides/troubleshooting/meshery-operator-meshsync)**.  

This guide provides:  
✅ **Step-by-step troubleshooting** for common errors.  
✅ **Insights into Meshery's custom controllers** to help you diagnose issues.  
✅ **Best practices** for maintaining a smooth and stable Meshery deployment.  

Whether you're dealing with **installation problems, connectivity issues, or unexpected behavior**, this guide will walk you through resolving them efficiently.
