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

<p style="display:block">
<img src="{{site.baseurl}}/assets/img/meshsync/meshsync.svg" align="left" 
    style="margin-right:1rem;margin-bottom:.5rem;" width="20%" />

MeshSync is a custom Kubernetes controller that provides tiered discovery and continual synchronization with Meshery Server as to the state of the Kubernetes cluster, service meshes, and their workloads.

</p>

MeshSync's working snapshot of the state of each cluster and service mesh under management is stored in-memory and continuously refreshed.

MeshSync is managed by the <a href="{{site.baseurl}}/concepts/architecture/operator">Meshery Operator</a>.

## What are MeshSync's core responsibilities?

### Discover existing resources inside the kubernetes cluster

The resources that are present inside the cluster are discovered efficiently with the help of pipelines. The data is constructed in a particular format specific to meshery and publish across to different parts of the architecture.

### Listening to events/changes on every component

MeshSync implements several informers/listeners across each resource to listen to changes that occur to them. These are reflected at real time and are updated in their respective areas.

## MeshSync FAQs

