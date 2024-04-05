---
layout: page
title: Contributing to Model Relationships
permalink: project/contributing/contributing-relationships
redirect_from: project/contributing/contributing-relationships/
abstract: How to contribute to Meshery Models Relationships, Policies...
language: en
type: project
category: contributing
list: include
---

Relationships within Model play a crucial role in establishing concrete visualisations of efficient data flow between different components of Meshery. These are used to classify the nature of interaction between one or more interconnected Components.

## Overview of Tasks

1. Identify the relationship and any specific constraints to be enforced between the two specific components, their models, or potentially other components, models, or environmental considerations.
2. Propose the appropriate relationship type, using one of the predefined set of relationship types, or suggest a new relationship where an existing type does not fit.
3. Create a Relationship Definition (yaml), including the following attributes:
4. Identify an existing OPA policy as the `evaluationQuery` suitable to the relationship. If no policy exists, propose a new policy (rego).
5. Submit a pull request to the Meshery repository with the new relationship definition.

### Relationship Visualizations

{% include relationships.html %}

### Defining Relationships: Kinds, Types, Subtypes

Relationships are defined in the `relationships.yaml` file in the Meshery repository. Each relationship definition includes the following attributes:

- `kind`: The genre of relationship (e.g., hierarchical, edge, sibling).
- `type`: The augmentative category of the relationship (e.g., binding, non-binding, inventory).
- `subType`: The specific represenative visual paradigm (e.g., parent, mount, network, wallet, badge).
- `selectors`: The scope of the relationship, including the components involved and any constraints.
- `evaluationQuery`: The OPA policy to invoke for relationship evaluation.
- `documentation`: A description of the relationship, its purpose, and any constraints or considerations.

### Existing Relationships as Examples

1. **Hierarchical** relationships involve either an ancestral connection of the components i.e. the creation/ deletion of a Component higher up affects the existence of the Components below in the lineage or a connection that involves the inheritance of features from one Component to the other.

   1. **Parent**: A parent-child relationship implies the requirement of the parent component before the child component can be created. For example, a "Namespace" in Kubernetes can be a parent of "Pods" within that namespace. The namespace must exist before creating pods within it.
   2. **Inventory**: A hierarchical inventory relationship implies the configuration of a(parent) component is patched with the configuration of other (child) component. For example, Wasm filters can inherit features and functionalities from Envoy filters. This can be used to build on existing functionalities provided by Envoy filters and further extend them using Wasm filters. It enables a modular and scalable approach to customize the behavior of the proxy while maintaining a clear hierarchy of features.

2. **Edge** relationships indicate the possibility of traffic flow between two components. They enable communication and interaction between different Components within the system.

   1. **Mount**: This subtype addresses the storage and access possibility between involved components. For example, a "PersistentVolume" can be mounted to a "Pod" to provide persistent storage for the pod's data.
   2. **Network**: This deals with IP addresses and DNS names and provides stable endpoints for communication. For example, a "Service" provides a stable endpoint for accessing multiple replicas of a "Deployment".
   3. **Firewall**: This acts as an intermediary for communications which include standard networking protocols like TCP and UDP. It can enforce network policies to control traffic between components. For example, a "NetworkPolicy" can be used to manage the traffic flow between different "Pods" in the cluster.
   4. **Permission**: This defines the permissions for components if they can have a possible relationship with other Components. It ensures that only authorized Components can interact with each other. For example, a "Role" can define permissions for Components to access specific resources.

3. **Sibling** relationships represent connections between components that are at the same hierarchical level or share a common parent. Siblings can have the same or similar functionalities or may interact with each other in specific ways. These relationships facilitate communication and cooperation between components that are in the same group or category. For example, a Service and a Pod in Kubernetes are siblings as they share a common parent and are at the same hierarchical level.

### Structure of Selectors

Selectors are structured as an array, wherein each entry comprises a 'from(self)' field and a 'to(other)' field (`[from: [{..}], to: [{..}]]`), delineating the components involved in a particular relationship. These entries define the constraints necessary for the existence of said relationship, thus providing scoping within a relationship.
Each item in the selector uniquely defines a relation between the components listed. i.e. `from` and `to` fields are evaluated within the context of the selector.

Only the components within the same selector relate to each other via 1:many kind of relation between components listed inside the `from` and `to` field. i.e. Each object inside the `from` relates to each item inside the `two` field within a particular selector.

When defining relationships that involve a large number of combinations between `from` and `to`, selectors provide a mechanism to organize and manage these relationships. This prevents the need for crafting complex deny attributes and facilitates easier maintenance.

This arrangement enhances flexibility and reusability in the definition and configuration of relationships among components.

<details close>
<summary>Selector example</summary>

{% highlight yaml %}
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
            },
            {
               "kind": "EBPFFilter",
   .....
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
            },
            {
               "kind": "WASMPlugin",
   ....
            }
   ...
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
                  "description": "In Kubernetes, ConfigMaps are a versatile resource that can be referenced by various other resources to provide configuration data to applications or other Kubernetes resources.\n\nBy referencing ConfigMaps in these various contexts, you can centralize and manage configuration data more efficiently, allowing for easier updates, versioning, and maintenance of configurations in a Kubernetes environment."
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
            },
            {
               "kind": "StatefulSets",
               "model": "kubernetes",
               "patch": {
   ....
               }
            }
   ...
         ]
      },
      "deny": {
   ...
      }
   }
]
{% endhighlight %}

The `selector` defined for the relationship between `WasmFilter` and `EnvoyFilter` (the first item in the array) is entirely independent of the `selector` defined for the relationship between `ConfigMap` and `Deployment`. This ensures independence in how these components relate to each other while still permitting similar types of relationships.

The above relation shows `WASMFilter` and `EBPFFilter` defined inside `from` relates to each component defined inside `to` `(EnvoyFilter, WASMPlugin...)`.
Similarly, `ConfigMap` defined inside `from` relates to each component defined inside `to` `(Deployment, StatefulSet,...)`

</details>

#### What is `evaluationQuery` attribute and how to determine the value for `evaluationQuery` inside a relationship definition?

As all relationship definitions are backed by OPA policies and the relationships depending upon their Kind and Subtype needs to be evaluated with respective policies, the policy to invoke for evaluation is determined by the property `evaluationQuery`, which follows the convention as `kind_subtype_relationship`.

```
Eg: If you are defining/updating a relationship definition with kind: Edge and subType: Network, the value for the attribute `evaluationQuery` should be edge_network_relationship.
```

Each policy has a set of rules defined and the `evaluationQuery` attribute corresponds to the main rule defined inside the policy, during the policy eval the results are collected from this rule.

### Configuring the scopes of the relationship definitions

The extent to which a relationship affects components within a model or beyond a model is defined and controlled using scopes. Scopes are defined using the `model` and `version` attributes within the relationship schema.

#### Global Scope

Relationships can be confined to a specific model, a specific model version, or can be allowed to affect all models. The relationship schema has a `model` and `version` attribute which facilitates this control. For example, if the model is specified as `aws-ec2-controller`, the relationship will work for those components that belong to the `aws-ec2-controller` model.

#### Local Scope

Scope is defined and controlled via the `selectors` [Selectors](#structure-of-selectors) attribute in the relationships.

### Best practices for defining new relationships

1. Ensure that the `deny` selectors and `allow` selectors do not conflict with each other i.e. relations are not getting overlapped for `allow` and `deny` selectors.
2. To configure a relationship to be applied across models, ensure the `model` property for those relationships is set to `*`, to limit the relationships to a specific model, specify the correct `model`(case sensitive).
3. To configure a relationship to be applied across all versions of a particular model, ensure the `version` property for those relationships is set to `*`, to limit the relationships to a specific version of a model, specify the correct model version.
4. Support for specifying version property as a regex to ensure relationships are applied to a subset of versions of a model is coming soon.
5. The `evaluationQuery` property determines the OPA policy to invoke for relationship evaluation, specify the correct rego query.
6. If a path `mutatedRef/mutatorRef` contains more than one array path then only the first array position can be specified as \_ for others explicitly mention them as 0
7. Currently `mutatedRef` doesn’t support having an array

#### Things to keep in mind while defining relationships

1. Targets of a Relationship can be specific Components or entire Models.
1. The values for Kind, Version, and Model are case-sensitive.
1. The convention is to use camel-casing for Kind and SubType values.
1. Absence of a field means in the selector means “\*” (or wildcard).
   - If we have a selector with {Kind: Pod, Model: Kubernetes}, the absence of the Version field here
     means that all the versions of the Kubernetes Pod resource (e.g. k8s.io/v1/betav2) will match.
1. In the event of conflicting Relationship Definitions, the union between them is taken.
   - If we have two Relationships, one from (Component A) to (Component B and Component F), and another
     from (Component A) to (Component B and Component C), then it is similar to having a Relationship
     from Component A to Component B, C and F
1. In the event of an overlapping set of complementary Relationship Definitions, Union.
1. In the event of an overlapping set of conflicting Relationship Definitions:
   - No relationship type (Kind) is inherently more important than the next one, so will not be any case
     of conflict

For more information refer - [Model - Construct Models in Meshery](https://docs.google.com/document/d/16z5hA8qVfSq885of9LXFUVvfom-hQXr-6oTD_GgoFmk/edit)
