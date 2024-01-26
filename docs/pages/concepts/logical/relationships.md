---
layout: enhanced
title: Relationships
permalink: concepts/logical/relationships
type: concepts
abstract: "Meshery Relationships identify and facilitate genealogy between Components."
language: en
list: include
---

[Relationships](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships) define the nature of interaction between interconnected components in Meshery. They represent various types of connections and dependencies between components no matter the genealogy of the relationship such as hierarchical,  relationships. Relationships have selectors, metadata, and optional parameters.

{% include/alert.html type="info" title="Contributing a new Relationship" content="<a href='https://docs.meshery.io/project/contributing/contributing-models#contribute-to-meshmodel-relationships'>contributing to relationships</a>" %}

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

## Selectors in Relationships

Selectors in Meshery are a powerful tool for defining relationships between different components. Selectors are organized in an array, providing a flexible way to specify match conditions and  corresponding actions. Multiple sets of selectors may be defined in order to precisely match and not match specific components and their attributes.

### Structure of Selectors

Selectors consist of an array, and each item in the array contains the following fields:

- `[from, to]`: Defines the relationship between components.


### Example Selectors

Selectors can be applied to various components, enabling a wide range of relationship definitions. Here are some examples:

1. **WASMFilter - EnvoyFilter**
2. **ConfigMap - Pod**
3. **ConfigMap - Deployment**

These pairs establish hierarchical inventory relationships, and the policy and visual paradigm remain consistent across different components. To avoid ambiguity when defining relationships, selectors act as internal scopes within a relationship. This means that the `from` and `to` fields are evaluated within the context of the selector. For example, the following selector:

```yaml
selector: [
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
    }
]
```

### Types of Relationships
1. **Inventory** 
2. **Mount** 
3. **Firewall** 
4. **Permission**
5. **Network**

### Hierarchical Relationships

When defining relationships that involve a large number of combinations between from and to, selectors provide a mechanism to organize and manage these relationships hierarchically. This prevents the need for crafting complex deny attributes and facilitates easier maintenance.

### Patch Strategies

Patches in Meshery relationships utilize strategies and references (mutatorRef/mutatedRef) for the from and to fields. These convey the property path that will be updated as the relationship is created.

{% include/alert.html type="info" title="Future Feature" %}

```
mesheryctl model import -f [ oci:// | file:// ]`
```
