---
layout: default
title: MeshSync
permalink: concepts/architecture/meshsync
type: concepts
redirect_from: architecture/meshsync
abstract: "Meshery offers support for kubernetes cluster/service mesh state synchronization with the help of MeshSync."
language: en
list: include
---

<p>
<img src="{{site.baseurl}}/assets/img/meshsync/meshsync.svg" align="left" style="margin-right:1rem;margin-bottom:.5rem;" width="20%" />

MeshSync is a custom Kubernetes controller that provides tiered discovery and continual synchronization with Meshery Server as to the state of the Kubernetes cluster, service meshes, and their workloads.

</p>

MeshSync's working snapshot of the state of each cluster and service mesh under management is stored in-memory and continuously refreshed.

MeshSync is managed by the <a href="{{site.baseurl}}/concepts/architecture/operator">Meshery Operator</a>.

## What are MeshSync's core responsibilities?

### Discover existing resources inside the kubernetes cluster

The resources that are present inside the cluster are discovered efficiently with the help of pipelines. The data is constructed in a particular format specific to meshery and publish across to different parts of the architecture.

### Listening to events/changes on every component

The informer in MeshSync actively listens to changes in resources and updates them in real time based on the provided informer 
configuration in the CRD.

## MeshSync FAQs

### How to reconfigure Meshsync informer?  
Meshsync is managed by [Operator]({{site.baseurl}}/concepts/architecture/operator) which watches for changes on Meshsync CRD for 
changes and updates the deployed resource accordingly. You can update the CRD using kubectl through the following steps
- Download the CRD with `kubectl get crd meshsyncs.meshery.layer5.io -o yaml > meshsync.yaml`
- Open the downloaded file and edit the field `informer_config` to blacklist all the types of resources that you don't want updates from.
- Apply the new definition with `kubectl apply -f meshsync.yaml` 

