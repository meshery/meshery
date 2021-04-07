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

MeshSync is a custom Kubernetes controller that provides tiered discovery and continual synchronization with Meshery Server as to the state of the Kubernetes cluster, service meshes, and their workloads.

## What are its core responsibilities?

### Discover existing resources inside the kubernetes cluster
The resources that are present inside the cluster are discovered efficiently with the help of pipelines. The data is constructed in a particular format specific to meshery and publish across to different parts of the architecture.

### Listening to events/changes on every component
MeshSync implements several informers/listeners across each resource to listen to changes that occur to them. These are reflected at real time and are updated in their respective areas.

### MeshSync FAQs
