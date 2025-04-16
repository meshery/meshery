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
To check if your Broker instance is running smoothly (it's deployed as a Kubernetes StatefulSet), follow these quick checks:

- Confirm that the Broker pods are running.
- Verify your cluster supports `LoadBalancer` or `NodePort` service types.
- Make sure the Meshery Server can reach the Broker service.

Still seeing issues? The **[Meshery Troubleshooting Guide](https://docs.meshery.io/guides/troubleshooting/meshery-operator-meshsync)** covers common problems with the Broker, MeshSync, and Operator — and offers clear steps to resolve them.


This guide provides:  
- **Step-by-step troubleshooting** for common errors.  
- **Insights into Meshery's custom controllers** to help you diagnose issues.  
- **Best practices** for maintaining a smooth and stable Meshery deployment.  

Whether you're dealing with **installation problems, connectivity issues, or unexpected behavior**, this guide will walk you through resolving them efficiently.
