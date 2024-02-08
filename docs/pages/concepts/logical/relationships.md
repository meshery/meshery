---
layout: enhanced
title: Relationships
permalink: concepts/logical/relationships
type: concepts
abstract: "Meshery Relationships identify and facilitate genealogy between Components."
language: en
list: include
redirect_from:
- concepts/relationships
---

[Relationships](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships) define the nature of interaction between interconnected components in Meshery. They represent various types of connections and dependencies between components no matter the genealogy of the relationship such as parent, siblings, binding. Relationships have selectors, metadata, and optional parameters.

{% include/alert.html type="info" title="Contributing a new Relationship" content="<a href='https://docs.meshery.io/project/contributing/contributing-models#contribute-to-meshmodel-relationships'>contributing to relationships</a>" %}

## Semantic and Non-Semantic Relationships

Shapes (all styles) convertable to Components

The `isAnnotation` attribute of a Component determines whether the give Component reflects an infrastructure concern - is sematically meaningful, and whose lifecycle is managed by Meshery.

## Core Constructs of Relationships

- Kinds
- Subtypes
- Selectors

## Kind - Subtypes of Relationships
The combination of `kind` and `subType` uniquely determines the visual paradigm for a given relationship i.e., relationships with the same `kind` and `subType` will share an identical visual representation regardless of the specific components involved in the relationship.
### 1. Edge - Network
Configures the networking.

**Example**
- Service --> Deployment,
- Service --> Pod,
- IngressController --> Ingress --> Service

### 2. Edge - Mount
**Example**

Assignment of PersistentVolumes to Pods via PersistenVolumeClaim .

- Pod --> PersistenVolumeClaim --> PersistentVolume

### 3. Edge - Permission
**Example**

The set of Service Accounts that are entitled with the one or more Roles/ClusterRoles bound via Role/ClusterRoleBinding.

- ClusterRole --> CluserRoleBinding --> ServiceAccount
- Role --> RoleBinding --> ServiceAccount

### 4. Edge - Firewall

#### Example
Kubernetes Network Policy for controlling ingress and egress traffic from Pod-to-Pod

- Pod --> NetworkPolicy --> Pod

### 5. Heirarchical - Inventory

**Example**

- WASMFilter (binary and configuration) --> IstioWASMPlugin
- WASMFilter (binary and configuration) --> IstioEnvoyFilter

### 6. Heirarchical - Parent

**Example**

Any namespaced Kubernetes component --> Kubernetes Namespace

## Selectors in Relationships


### Example Selector

Selectors can be applied to various components, enabling a wide range of relationship definitions. Here are some examples:
1. **ConfigMap - Pod**
2. **ConfigMap - Deployment**
3. **WASMFilter - EnvoyFilter**

The above pairs have hierarchical inventory relationships, and visual paradigm remain consistent across different components.
A snippet of the selector backing this relationship is listed below.


 <!-- add images -->
### Types of Relationships
1. **Inventory** 
2. **Parent**
3. **Mount** 
4. **Firewall** 
5. **Permission**
6. **Network**

### How Relationships are formed?
1. You can create relationships manually by using the edge handles, bringing related components to close proximity or dragging a node inisde other node.

_Note: It may happen that, you created a relationship from the UI, but the [Policy Engine]({{site.baseurl}}/concepts/logical/policies) disapproved/overrided the decision if all the constraints for a particular relationship are not satisfied._

2. Relationships gets auto-created if the user update's the node config such that the relationship criteria is satisfied. _Open the [catalog item](https://meshery.io/catalog/deployment/7dd39d30-7b14-4f9f-a66c-06ba3e5000fa.html) in playground and follow the steps in the description._

When the relationships are created by the user, almost in all cases the config of the involved components are patched. To see the specific of patching refer [Patch Strategies](#patch-strategies)

The Designs are evaluated by the [Policy Engine]({{site.baseurl}}/concepts/logical/policies) for potential relationships 
<!-- Explain how and what configs get patched when relationships are created -->
<!-- Explain real time evaluationof relationships on -->
<!-- 1. Import -->
<!-- 2. When compoennt config is update and it statisfied the condition for the relationship -->

### Patch Strategies

Patches in Meshery relationships utilize strategies and references (mutatorRef/mutatedRef) for the from and to fields. These convey the property path that will be updated as the relationship is created.

## Meshery Registry

You can reference and search the full set of registered relationships in Meshery UI --> Setttings --> Registry

{% include/alert.html type="info" title="Future Feature" %}

```
mesheryctl model import -f [ oci:// | file:// ]`
```

