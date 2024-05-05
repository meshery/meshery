---
layout: default
title: Relationships
permalink: concepts/logical/relationships
type: concepts
abstract: "Meshery Relationships identify and facilitate genealogy between Components."
language: en
list: include
redirect_from:
  - concepts/relationships
---

**Meshery uses relationships to map how interconnected components interact.** These relationships are defined by a combination of properties: kind, type, and subtype. This allows you to model various connections between components, including:

- **Hierarchical structures:** Parent-child relationships show clear lineage, similar to a family tree (child, parent, grandparent, etc.).
- **Interdependencies:** This captures how components rely on each other to function.
- **Collateral connections:** These describe components that share a common origin but operate independently (siblings, cousins, etc.).
- **Non-genealogical ties:** Relationships like "parent," "sibling," or "binding" can exist regardless of ancestry.

Relationships are further customized by:

- **Selectors:** These specify which components the relationship applies to.
- **Selector sets:** Combine multiple selectors for more granular control.
- **Metadata:** Provides additional descriptive information about the relationship.
- **Optional parameters:** Allow for further tailoring of the relationship behavior.

Relationships define the nature of interaction between interconnected components in Meshery. They represent various types of connections and dependencies between components no matter the genealogy of the relationship such as parent, siblings, binding. Relationships have selectors, metadata, and optional parameters.

{% include/alert.html type="dark" title="Contributor Guide to Meshery Relationships" content="If you want to create a new relationship definition or modify existing relationship definitions, see the <a href='https://docs.meshery.io/project/contributing/contributing-relationships'>Contributing to Meshery Relationships</a> guide." %}

## Types of Relationships

Meshery supports a variety of relationships between components. Relationships are are categorized into different kinds, types, and subtypes, so that can be expressive of the specific manner in which one or more components relate to one another. Each type of relationship can be interpretted by Meshery UI (or other [extensions](/extensibility/extensions)) and mapped to a specific visual paradigm for the given kind relationship. Let's look at some examples of these visual paradigms; let's explore examples of way in which relationships are represented in Meshery.

Here is a list of the different types of relationships that Meshery supports:

1. Edge
   1. Network
   1. Firewall
   1. Binding
      1. Mount
      1. Permission
1. Heirarchical
   1. Inventory
   1. Parent

{% include relationships.html %}
## The Meaning of Relationships

Meshery supports a variety of relationships between components. These relationships are categorized into two types: Semantic and Non-Semantic. Relationships are categorized by whether they are meaningful in terms of how Meshery manages entities - a Semantic relationship - or are simply annotations to aid in the comprehension of you as the designer of the relationship - a Non-Semantic relationship.

### Semantic Relationships

Semantic relationships are those that are meaningful in the context of the application or infrastructure. For example, a `Service` in Kubernetes is semantically related to a `Deployment` or a `Pod`. These relationships are meaningful and are managed by Meshery.

<!-- @iArchitSharma, help, if you would please.

_[TODO: a visual example is needed here]_ -->

### Non-Semantic Relationships

Non-semantic relationships are those that are meaningful to you as a user and your mental representation of your infrastructure and applications, but are not meaningful in terms of how Meshery evaluates the design or manages these relationships and their associated components. Non-sematic relationships are ignored by Meshery's lifecycle management engine. For example, a `Rectangle` shape that encloses other components (has a parent relationship with other child components) is not semantically meaningful to the way in which Meshery manages these resources. While the `Rectangle` shape might have a parent-child relationship with any number of Meshery-managed components, such a relationship does not implicate any management that Meshery might perform; they are not managed by Meshery.

<!-- @iArchitSharma, help, if you would please.

_[TODO: a visual example is needed here]_ -->

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
<details close><summary>Visual Representation of Edge-Network Relationships</summary>
           <br>
           <figure><figcaption>1. Edge - Network: Ingress to Service<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=1f79b0c6-2efe-4ee9-b08c-e1bd07a3926b"> (open in playground)</a></figcaption>
           <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_network_ingress_to_service_relationship.svg"/>
           </figure>
           <figure><figcaption>2. Edge - Network: Service to Pod<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=90a9b4a0-a296-44b5-b1c5-7b1cb4827a77"> (open in playground)</a></figcaption>
           <img alt="Edge - Network: Ingress to Service" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_network_service_to_pod_relationship.svg"/>
           </figure>
           <figure><figcaption>3. Edge - Network: Service to Service<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=4e368e07-5039-400e-b637-96b0241af799"> (open in playground)</a></figcaption>
           <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_network_service_to_service_relationship.svg"/>
           </figure>
           <figure><figcaption>4. Edge - Network: Service to Endpoint<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=ab35416d-7cf7-4540-8b2e-7271ffeadde2"> (open in playground)</a></figcaption>
           <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_network_service_to_endpoints_relationship.svg"/>
           </figure>
           <figure><figcaption>5. Edge - Network: Service to Deployment<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=33742281-428d-4340-b42e-6a0fd4ba1d0a"> (open in playground)</a></figcaption>
           <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_deployment.svg"/>
           </figure>
   </details>

### 2. Edge - Mount

**Example**: Assignment of PersistentVolumes to Pods via PersistenVolumeClaim.

- Example 1) Pod --> PersistenVolumeClaim --> PersistentVolume
<details close><summary>Visual Representation of Edge-Mount Relationship</summary>
           <br>
           <p>1. Edge - Mount: Pod and Persistent volume via Persistent volume claim<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=43d5fdfe-25f8-4c2c-be9d-30861bbc2a08"> (open in playground)</a> </p>
           <figure>
           <img alt="Edge - Mount" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_mount_relationship_pod_persistent_volume.svg"/>
           </figure>
   </details>

### 3. Edge - Permission

**Example**: The set of Service Accounts that are entitled with the one or more Roles/ClusterRoles bound via Role/ClusterRoleBinding.

- Example 1) ClusterRole --> CluserRoleBinding --> ServiceAccount
- Example 2) Role --> RoleBinding --> ServiceAccount
<details close><summary>Visual Representation of Edge-Permission Relationship</summary>
           <br>
           <figure><figcaption>1. Edge - Permission: Cluster Role to Service Account <a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=7dd39d30-7b14-4f9f-a66c-06ba3e5000fa"> (open in playground)</a></figcaption>
           <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_permission_relationship_cluster_role_service_account.svg"/>
           </figure>
   </details>

### 4. Edge - Firewall

Kubernetes Network Policy for controlling ingress and egress traffic from Pod-to-Pod

- Example 1) Pod --> NetworkPolicy --> Pod
<details close><summary>Visual Representation of Edge-Firewall Relationship</summary>
           <br>
           <figure><figcaption>1. Edge - Firewall: Pod to Pod<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=58fda714-eaa4-490f-b228-b8bcfe3a1e47s"> (open in playground)</a></figcaption>
           <img alt="Edge - Firewall" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/edge_firewall_relationship_pod_to_pod.svg">
           </figure>
   </details>

### 5. Heirarchical - Inventory

**Example**

- Example 1) (binary and configuration) --> IstioWASMPlugin
- Example 2) WASMFilter (binary and configuration) --> IstioEnvoyFilter
<details close><summary>Visual Representation of Hierarchical-Inventory Relationship</summary>
           <figure><br><figcaption>1. Hierarchical - Inventory: Namespace and ConfigMap<a target="_blank" href="https://playground.meshery.io/extension/meshmap?design=21d40e36-8ab7-4f9f-9fed-f6a818510446"> (open in playground)</a></figcaption>
           <img alt="Hierarchical Inventory Relationship" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/hierarchical_inventory_relationship.svg"/>
           </figure>
   </details>

### 6. Heirarchical - Parent

**Example**:

- Example 1) Any namespaced Kubernetes component --> Kubernetes Namespace
<details close><summary>Visual Representation of Hierarchical-Parent Relationship</summary>
           <figure><br><figcaption>1. Hierarchical - Parent: Namespace (Parent) and ConfigMap (child), Role (Child) <a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=6370ffcd-13a6-4a65-b426-30f1e63dc381"> (open in playground)</a></figcaption>
           <img alt="Hierarchical Parent Relationship" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/hierarchical_parent_relationship.svg"/>
           </figure>
   </details>

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

```json
"selector": {
    "allow": {
        "from": [
          {
            "kind": "ConfigMap",
            "model": "kubernetes",
            "patch": {
              "patchStrategy": "replace",
              "mutatorRef": [
                [
                  "name"
                ]
              ],
              "description": "In Kubernetes, ConfigMaps are a versatile resource that can be referenced by various other resources to provide configuration data to applications or other Kubnernetes resources.\n\nBy referencing ConfigMaps in these various contexts, you can centralize and manage configuration data more efficiently, allowing for easier updates, versioning, and maintenance of configurations in a Kubernetes environment."
            }
          }
        ],
        "to": [
          {
            "kind": "Pod",
            "model": "kubernetes",
            "patch": {
              "patchStrategy": "replace",
              "mutatedRef": [
                [
                  "settings",
                  "spec",
                  "containers",
                  "_",
                  "envFrom",
                  "0",
                  "configMapRef",
                  "name"
                ]
              ],
              "description": "ConfigMaps can be referenced in the Pod specification to inject configuration data into the Pod's environment.\n\nThe keys from the ConfigMap will be exposed as environment variables to the container within the Pod."
            }
          }
        ]
    }
}
```

The above snippet defines a selector configuration for allowing relationships between `Kubernetes ConfigMap` and `Kubernetes Pod`.

 <!-- add images -->

## Relationship Evaluation

![Meshery Relationship](/assets/img/concepts/logical/relationship-evaluation-flow.svg)

### How Relationships are formed?

1. You can create relationships manually by using the edge handles, bringing related components to close proximity or dragging a component inside other component. It may happen that, you created a relationship from the UI, but the <a href='/concepts/logical/policies)'>Policy Engine</a> rejected or overrode the decision if all the constraints for a particular relationship are not satisfied.

2. Relationships are automatically created when a component's configuration is modified in a way that relationship criteria is satisfied.

{% include/alert.html type="info" title="Explore an example relationship" content="To explore an example of this behavior, see the <a href='https://meshery.io/catalog/deployment/7dd39d30-7b14-4f9f-a66c-06ba3e5000fa.html'>Example Edge-Permission Relationship</a> and follow the steps written in its description." %}

When the relationships are created by the user, almost in all cases the config of the involved components are patched. To see the specific of patching refer [Patch Strategies](#patch-strategies)

Designs are evaluated by the [Policy Engine]({{site.baseurl}}/concepts/logical/policies) for potential relationships

<!-- Explain how and what configs get patched when relationships are created -->
<!-- Explain real time evaluationof relationships on -->
<!-- 1. Import -->
<!-- 2. When compoennt config is update and it statisfied the condition for the relationship -->

### Patch Strategies

Patches in Meshery relationships utilize strategies and references (mutatorRef/mutatedRef) for the from and to fields. These convey the property path that will be updated as the relationship is created.

# Itemizing your Relationship Definitions in your Meshery deployment

In any given Meshery deployment, you can reference and search the full set of registered relationships (in Meshery's internal registry) in using either of Meshery's client interfaces.

**Meshery UI**

- Visit _Setttings_ --> _Registry_

**Meshery CLI**

- Run `mesheryctl relationship list`

<!--
```
mesheryctl model import -f [ oci:// | file:// ]`
``` -->
