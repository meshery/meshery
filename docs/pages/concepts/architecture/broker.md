---
layout: default
title: Broker
permalink: concepts/architecture/broker
type: concepts
redirect_from: architecture/broker
abstract: "Meshery broker component fascilitates data streaming between kubernetes cluster components and outside world."
language: en
list: include
---

Broker is a custom Kubernetes controller that provides data streaming across independent components of Meshery whether those components are running inside or outside of the Kubernetes cluster.

### Broker FAQs

#### How many Brokers can run?
It is recommended to run one broker instance for each kubernetes cluster, However the instance itself can be scaled up based on the incoming data volume in each of the cluster. The scaling is independent of the number of instances running.

#### What does an HA configuration look like?
We leverage on the kubernetes functionality in terms of the High-Availability behaviour. Meaning, the broker instance gets instantiated/restarted on its own when an issue occurs. In part, Meshery-Operator is also resposible for keeping the broker functional.

#### What stateful characteristics does the Broker have?
All the messages that are published to the broker is persisted in-memory within the broker instance until it get consumed. Persistent-volume/Disk-space is not currently being used by the Broker.

#### How do I know if the Broker is working? How do I troubleshoot the Broker?
The Broker instance is deployed inside the kubernetes cluster as a `Statefulset`. In the case where the broker doesnt seem to work, here are a few steps to troubleshoot the instance:

- Make sure the pods corresponding to the `Statefulset` is up and running.
- Make sure the kubernetes cluster has support for kubernetes `Service` type `LoadBalancer` or `NodePort`.
- Ensure connectivity between the Meshery-Server and the Broker service endpoint.
