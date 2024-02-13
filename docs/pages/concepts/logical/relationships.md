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

Relationships define the nature of interaction between interconnected components in Meshery. They represent various types of connections and dependencies between components no matter the genealogy of the relationship such as parent, siblings, binding. Relationships have selectors, metadata, and optional parameters.

{% include/alert.html type="dark" title="Contributor Guide to Meshery Relationships" content="If you want to create a new relationship definition or modify existing relationship definitions, see the <a href='https://docs.meshery.io/project/contributing/contributing-models#contribute-to-meshmodel-relationships'>Contributing to Meshery Relationships</a> guide." %}

## Semantic and Non-Semantic Relationships

### Semantic Relationships

Semantic relationships are those that are meaningful in the context of the application or infrastructure. For example, a `Service` in Kubernetes is semantically related to a `Deployment` or a `Pod`. These relationships are meaningful and are managed by Meshery.

_[TODO: a visual example is needed here]_

### Non-Semantic Relationships

Non-semantic relationships are those that are meaningful to you as a user and your mental representation of your infrastructure and applications, but are not meaningful in terms of how Meshery evaluates the design or manages these relationships and their associated components. Non-sematic relationships are ignored by Meshery's lifecycle management engine. For example, a `Rectangle` shape that encloses other components (has a parent relationship with other child components) is not semantically meaningful to the way in which Meshery manages these resources. While the `Rectangle` shape might have a parent-child relationship with any number of Meshery-managed components, such a relationship does not implicate any management that Meshery might perform; they are not managed by Meshery.

_[TODO: a visual example is needed here]_

#### Identifiying Non-Semantic Relationships

The `isAnnotation` attribute of a Relationship or Component determines whether the given Relationship or Component represents a management concern for Meshery; whether the given Relationship or Component is sematically meaningful, and whose lifecycle is managed by Meshery.

## Core Constructs of Relationships

- Kinds
- Subtypes
- Selectors

## Kind and Subtypes of Relationships

The combination of `kind` and `subType` uniquely determines the visual paradigm for a given relationship i.e., relationships with the same `kind` and `subType` will share an identical visual representation regardless of the specific components involved in the relationship.

### 1. Edge - Network

This Relationship type configures the networking between one or more components.

**Examples**: Relationships between a Service and a Deployment, or between a Service and a Pod, or between an Ingress and a Service.

- Example 1) Service --> Deployment
- Example 2) Service --> Pod
- Example 3) IngressController --> Ingress --> Service

### 2. Edge - Mount

**Example**: Assignment of PersistentVolumes to Pods via PersistenVolumeClaim.

- Example 1) Pod --> PersistenVolumeClaim --> PersistentVolume

### 3. Edge - Permission

**Example**: The set of Service Accounts that are entitled with the one or more Roles/ClusterRoles bound via Role/ClusterRoleBinding.

- Example 1) ClusterRole --> CluserRoleBinding --> ServiceAccount
- Example 2) Role --> RoleBinding --> ServiceAccount

### 4. Edge - Firewall

Kubernetes Network Policy for controlling ingress and egress traffic from Pod-to-Pod

- Example 1) Pod --> NetworkPolicy --> Pod

### 5. Heirarchical - Inventory

**Example**

- Example 1) (binary and configuration) --> IstioWASMPlugin
- Example 2)  WASMFilter (binary and configuration) --> IstioEnvoyFilter

### 6. Heirarchical - Parent

**Example**:

- Example 1) Any namespaced Kubernetes component --> Kubernetes Namespace

## Selectors in Relationships

In Meshery, a selector is a way to specify which set of components a certain other component should affect or interact with. Selectors provide a flexible and powerful way to manage and orchestrate resources within a under Meshery's management.

Selectors can be applied to various components, enabling a wide range of relationship definitions. Here are some examples:

<table class="table table-dark table-active">
    <tr>
        <th>Model Component</th>
        <th>Relationship Kind</th>
        <th>Relationship SubType</th>
        <th>Model Component</th>
    </tr>
    <tr>
        <td>Kubernetes ConfigMap</td>
        <td>Hierarchical</td>
        <td>Inventory</td>
        <td>Kubernetes Pod</td>
    </tr>
    <tr>
        <td>Kubernetes ConfigMap</td>
        <td>Hierarchical</td>
        <td>Inventory</td>
        <td>Kubernetes Deployment</td>
    </tr>
    <tr>
        <td>Meshery WASMFilter</td>
        <td>Hierarchical</td>
        <td>Inventory</td>
        <td>Istio EnvoyFilter</td>
    </tr>
</table>

The above relationships pairs have hierarchical inventory relationships, and visual paradigm remain consistent across different components. A snippet of the selector backing this relationship is listed below.

 <!-- add images -->

### Types of Relationships

1. **Inventory**
2. **Parent**
3. **Mount**
4. **Firewall**
5. **Permission**
6. **Network**

## Relationship Evaluation

![Meshery Relationship](/assets/img/concepts/logical/relationship-evaluation-flow.svg)

### How Relationships are formed?

1. You can create relationships manually by using the edge handles, bringing related components to close proximity or dragging a component inside other component. It may happen that, you created a relationship from the UI, but the <a href='/concepts/logical/policies)'>Policy Engine</a> rejected or overrode the decision if all the constraints for a particular relationship are not satisfied.

2. Relationships are automatically created when a component's configuration is modified in a way that relationship criteria is satisfied. 
  
{% include/alert.html type="info" title="Explore an example relationship" content="To explore an example of this behavior, see the <a href='https://meshery.io/catalog/deployment/7dd39d30-7b14-4f9f-a66c-06ba3e5000fa.html'>Example Edge-Permission Relationship</a> and follow the steps written in its description." %}

When the relationships are created by the user, almost in all cases the config of the involved components are patched. To see the specific of patching refer [Patch Strategies](#patch-strategies)

The Designs are evaluated by the [Policy Engine]({{site.baseurl}}/concepts/logical/policies) for potential relationships

<!-- Explain how and what configs get patched when relationships are created -->
<!-- Explain real time evaluationof relationships on -->
<!-- 1. Import -->
<!-- 2. When compoennt config is update and it statisfied the condition for the relationship -->

### Patch Strategies

Patches in Meshery relationships utilize strategies and references (mutatorRef/mutatedRef) for the from and to fields. These convey the property path that will be updated as the relationship is created.

# Itemizing your Relationship Definitions in your Meshery deployment

In any given Meshery deployment, you can reference and search the full set of registered relationships (in Meshery's internal registry) in using either of Meshery's client interfaces.

**Meshery UI**

- Visit *Setttings* --> *Registry*

**Meshery CLI**

- Run `mesheryctl relationship list`

<!-- 
```
mesheryctl model import -f [ oci:// | file:// ]`
``` -->