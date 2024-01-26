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
The combination of `kind` and `subType` uniquely defines the visual paradigm for that relationship. i.e. For all relations that have same `kind` and `subType` have identical visual representation irrespective of components involved.

### 1. Edge - Network
Configures the networking.

**Example**
- Service --> Deployment,
- Service --> Pod,
- Ingress Controller --> Ingress --> Service

### 2. Edge - Mount
**Example**

Assignment of Pods to PersistentVolumes via PersistenVolumeClaim .

Pod --> Persisten Volume Claim --> Persistent Volume

### 3. Edge - Permission
**Example**

The set of Service Accounts that are entitled with the one or more Roles/ClusterRoles bound via Role/ClusterRoleBinding.

Cluster Role --> Cluser Role Binding --> Service Account

### 4. Edge - Firewall

#### Example
Kubernetes Network Policy for controlling ingress and egress traffic from Pod-to-Pod

Pod --> NetworkPolicy --> Pod

### 5. Heirarchical - Inventory

**Example**

WASM Filter (binary and configuration) --> Istio WASM Plugin
WASM Filter (Binary and Configuration) --> Istio Envoy Filter

### 6. Heirarchical - Parent

**Example**

Any namespaced Kubernetes component --> Kubernetes Namespace

## Selectors in Relationships

### Structure of Selectors
<!-- Define allow and deny -->

Selectors are structured as an array, wherein each entry comprises a 'from(self)' field and a 'to(other)' field (`[from: [{..}], to: [{..}]]`), delineating the components involved in a particular relationship. These entries define the constraints necessary for the existence of said relationship, thus providing scoping within a relationship. 
Each item in the selector, uniquely defines a relation between the components listed. i.e. `from` and `to` fields are evaluated within the context of the selector.

This arrangement enhances flexibility and reusability in the definition and configuration of relationships among components.

### Example Selector

Selectors can be applied to various components, enabling a wide range of relationship definitions. Here are some examples:
1. **ConfigMap - Pod**
2. **ConfigMap - Deployment**
3. **WASMFilter - EnvoyFilter**

The above pairs have hierarchical inventory relationships, and visual paradigm remain consistent across different components.
A snippet of the selector backing this relationship is listed below.
<details open>
<summary>
 <b>Selector</b>
</summary>


_https://github.com/meshery/meshery/blob/master/server/meshmodel/kubernetes/relationships/hierarchical_inventory.json_

```json
selector: [
  {
    "allow": {
      "from": [
        {
          "kind": "WASMFilter",
          "model": "istio-base",
          "patch": {
            "patchStrategy": "replace",
            "mutatorRef": [
              [
                "settings",
                "config"
              ],
              [
                "settings",
                "wasm-filter"
              ]
            ],
            "description": "WASM filter configuration to be applied to Envoy Filter."
          }
        }
      ],
      "to": [
        {
          "kind": "EnvoyFilter",
          "model": "istio-base",
          "patch": {
            "patchStrategy": "replace",
            "mutatedRef": [
              [
                "settings",
                "configPatches",
                "_",
                "patch",
                "value"
              ]
            ],
            "description": "Receive the WASM filter configuration."
          }
        }
      ]
    },
    "deny": {
      ...
    }
  },
  {
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
          "kind": "Deployment",
          "model": "kubernetes",
          "patch": {
            "patchStrategy": "replace",
            "mutatedRef": [
              [
                "spec",
                "containers",
                "_",
                "envFrom",
                "configMapRef",
                "name"
              ]
            ],
            "description": "Deployments can reference ConfigMaps to inject configuration data into the Pods they manage. This is useful for maintaining consistent configuration across replica sets.\n\nThe keys from the ConfigMap will be exposed as environment variables to the containers within the pods managed by the Deployment."
          }
        }
      ]
    },
    "deny": {
      ...
    }
  }
]
```
The `selector` defined for the relationship between `WasmFilter` and `EnvoyFilter` (the first item in the array) is entirely distinct from the `selector` defined for the relationship between `ConfigMap` and `Deployment`. This ensures independence in how these components relate to each other while still permitting similar types of relationships.
</details>

 <!-- add images -->
### Types of Relationships
1. **Inventory** 
2. **Mount** 
3. **Firewall** 
4. **Permission**
5. **Network**

### Hierarchical Relationships

When defining relationships that involve a large number of combinations between from and to, selectors provide a mechanism to organize and manage these relationships hierarchically. This prevents the need for crafting complex deny attributes and facilitates easier maintenance.

<!-- Explain how and what configs get patched when relationships are created -->
<!-- Explain real time evaluationof relationships on -->
<!-- 1. Import -->
<!-- 2. When compoennt config is update and it statisfied the condition for the relationship -->

### Patch Strategies

Patches in Meshery relationships utilize strategies and references (mutatorRef/mutatedRef) for the from and to fields. These convey the property path that will be updated as the relationship is created.

{% include/alert.html type="info" title="Future Feature" %}

```
mesheryctl model import -f [ oci:// | file:// ]`
```

## Meshery Registry

You can reference and search the full set of registered relationships in Meshery UI --> Setttings --> Registry
