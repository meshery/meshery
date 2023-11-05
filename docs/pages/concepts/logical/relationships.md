---
layout: default
title: Relationships
permalink: concepts/relationshps
type: concepts
abstract: "Meshery Relationships identify and facilitate genealogy between Components."
language: en
list: include
---

[Relationships](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships) define the nature of interaction between interconnected components in MeshModel. They represent various types of connections and dependencies between components, such as hierarchical, network, or default relationships. Relationships have selectors, metadata, and optional parameters.


## Semantic and Non-Semantic Relationships
Shapes (all styles) convertable to Components

The `isAnnotation` attribute of a Component determines whether the give Component reflects an infrastructure concern - is sematically meaningful, and whose lifecycle is managed by Meshery.

## Core Constructs of Relationships

- Kinds
- Subtypes


## Kind - Subtypes of Relationships

1. Edge - Network

#### Example
Service --> Deployment
Service --> Pod
Ingress Controller --> Ingress --> Service

1. Edge - Mount

#### Example
Assignment of Pods to Persistent Volumes via PVC.

Pod --> Persisten Volume Claim --> Persistent Volume
 
1. Edge - Permission

#### Example
The set of Service Accounts that are entitled with the one or more Cluster Roles bound via Cluster Role Binding.

Cluster Role --> Cluser Role Binding --> Service Account

1. Edge - Firewall

#### Example
Kubernetes Network Policy for controlling ingress and egress traffic from Pod-to-Pod

1. Heirarchical - Inventory

#### Example
WASM Filter (binary and configuration) --> Istio WASM Plugin
WASM Filter (Binary and Configuration) --> Istio Envoy Filter

1. Heirarchical - Parent

#### Example
Any namespaced Kubernetes component --> Kubernetes Namespace

## Meshery Registry
You can reference and search the full set of registered relationships in Meshery UI --> Setttings --> Registry



{% include/alert.html title="Future Feature" %}

```
mesheryctl model import -f [ oci:// | file:// ]`
```
