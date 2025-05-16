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

[Relationships](/concepts/logical/relationships) within [Models](/concepts/logical/models) play a crucial role in establishing concrete visualisations of efficient data flow between different components of Meshery. These are used to classify the nature of interaction between one or more interconnected [Components](/concepts/logical/components).

## Overview of Steps to Create Relationships

**Prework:**

1. [Relationship Identification](#relationship-identification)
2. [Relationship Classification](#relationship-visualizations)

**Development:**
3. [Relationship Definition](#relationship-definitions)
4. [Relationship Scopes](#relationship-scopes)

**Postwork:**
5. [Relationship Testing](#relationship-testing)
6. [Relationship Contribution](#relationship-contribution)

## Prework

<a id="relationship-identification"></a>

### 1. Characterize the relationship and any specific constraints

Using your domain expertise, define the qualities of this new relationship. Identify and qualify any specific constraints to be enforced between one or more specific components within the same or different models. Let's take an example to understand this better.

For example, you might know that a Kubernetes `Service` can have a network-based relatinship with a Kubernetes `Pod`. To codify this relationship, you would define the relationship as a `kind: edge` relationship with a `type: network`.

<details close>
<summary>Relationship Example</summary>
<pre><code class="language-yaml highlighter-rouge">
{
  "schemaVersion": "core.meshery.io/v1alpha2",
  "kind": "edge",
  "type": "network",
  "version": "v1.0.0",
  "metadata": {"description": "A relationship that defines network edges between components."}
  selector: [{
      "allow": {
         "from": [
            {
               "kind": "Service",
               "model": "kubernetes"
            }
         ],
         "to": [
            {
               "kind": "Pod",
               "model": "kubernetes"
            }
         ]
      },
      "deny": {}
   }]

</code></pre>
</details>

You might *also* know that this relationship is constrained by the presence of a Kubernetes `Deployment` as the parent of the Kubernetes `Pod`. This constraint would be codified in the relationship definition by including the <code>deny</code> function in your selector.

<details close>
<summary>Example Relationship with Constraints</summary>
<pre><code class="language-yaml highlighter-rouge">
selector: [
   {
      "allow": {
         "from": [
            {
               "kind": "Service",
               "model": "kubernetes"
            }
         ],
         "to": [
            {
               "kind": "Pod",
               "model": "kubernetes"
            }
         ]
      },
      "deny": {
         "from": [
            {
               "kind": "Service",
               "model": "kubernetes"
            }
         ],
         "to": [
            {
               "kind": "Pod",
               "model": "kubernetes",
               "parent": "Deployment"
            }
         ]
      }
   }
]
</code></pre>
</details>

Codify relationships leveraging your domain expertise. In this example, the relationship between the `Service` and `Pod` components in the Kubernetes model is defined with appropriate consideration of surrounding constraints (e.g. presence of a `Deployment`).

#### Understand Relationship Classifications

Relationships can be classified into three main categories:

1. **Hierarchical** relationships involve either an ancestral connection of the components i.e. the creation/ deletion of a Component higher up affects the existence of the Components below in the lineage or a connection that involves the inheritance of features from one Component to the other.

   1. **Parent**: A parent-child relationship implies the requirement of the parent component before the child component can be created. For example, a `Namespace` in Kubernetes can be a parent of `Pods` within that namespace. The namespace must exist before creating pods within it.
   2. **Inventory**: A hierarchical inventory relationship implies the configuration of a(parent) component is patched with the configuration of other (child) component. For example, Wasm filters can inherit features and functionalities from Envoy filters. This can be used to build on existing functionalities provided by Envoy filters and further extend them using Wasm filters. It enables a modular and scalable approach to customize the behavior of the proxy while maintaining a clear hierarchy of features.

2. **Edge** relationships indicate the possibility of traffic flow between two components. They enable communication and interaction between different Components within the system.

   1. **Mount**: This subtype addresses the storage and access possibility between involved components. For example, a `PersistentVolume` can be mounted to a `Pod` to provide persistent storage for the pod's data.
   2. **Network**: This deals with IP addresses and DNS names and provides stable endpoints for communication. For example, a `Service` provides a stable endpoint for accessing multiple replicas of a `Deployment`.
   3. **Firewall**: This acts as an intermediary for communications which include standard networking protocols like TCP and UDP. It can enforce network policies to control traffic between components. For example, a `NetworkPolicy` can be used to manage the traffic flow between different `Pods` in the cluster.
   4. **Permission**: This defines the permissions for components if they can have a possible relationship with other Components. It ensures that only authorized Components can interact with each other. For example, a `Role` can define permissions for Components to access specific resources.

3. **Sibling** relationships represent connections between components that are at the same hierarchical level or share a common parent. Siblings can have the same or similar functionalities or may interact with each other in specific ways. These relationships facilitate communication and cooperation between components that are in the same group or category. For example, a Service and a Pod in Kubernetes are siblings as they share a common parent and are at the same hierarchical level.

<a id="relationship-visualizations"></a>

### 2. Classify relationship type and specify visual representation

Browse and pick the most appropriate visualization for this relationship by using one of the predefined relationship visualizations.

{% include relationships.html %}

Once selected, note the relationship's `kind`, `type`, and `subtype` of your selected visualization. Alternatively, if an existing visualization does not seem appropriate for the relationship, please propose a new visualization at-will. Simply use the whiteboard feature of Meshery's extensions to sketch out the relationship and propose it as a new visualization.

## Development

<a id="relationship-definitions"></a>

### 3. Create a Relationship Definition as a JSON file

Create a relationship definition as a JSON file, placing this new definition file into its respective model folder (see [Contributing to Models](./contributing-models)). Relationship definition files are commonly named  `relationships.yaml` as a convention, however, this name is not required. A model may include any number of relationship definitions. Include the following attributes in your relationship definition:

- `kind`: The genre of relationship (e.g., hierarchical, edge, sibling).
- `type`: The augmentative category of the relationship (e.g., binding, non-binding, inventory).
- `subType`: The specific representative visual paradigm (e.g., parent, mount, network, wallet, badge).
- `selectors`: The scope of the relationship, including the components involved and any constraints. Selector specify to which component(s) the relationship applies or does not apply (think of in terms of the `AND` operators in a query). Selector Sets are used to combine multiple selectors for more granular control over the logic used when matchmaking (establishing a relationship) between components (think of in terms of the `OR` operators in a query).
- `evaluationQuery`: Name of the policy or policies (Open Policy Agent rego file(s)) to invoke for relationship evaluation. Identify an existing OPA policy as the `evaluationQuery` suitable to the relationship. If no policy exists, propose a new policy (rego). *(rarely necessary)* Create a new policy for the evaluation of your relationship using Rego. *This step is only necessary and can typically be skipped. Contact a maintainer if the relationship requires a new policy to evaluate the relationship.*
- `description`: A characterization of the relationship, its purpose, and any constraints or considerations of its application.

{% include alert.html title="Use Existing Relationships as Examples" type="info" content="Browse the <a href='https://github.com/meshery/meshery/tree/master/server/meshmodel'>existing relationships in the Meshery repository</a> to find examples of how to existing relationships, using them as a template. Alternatively, you can review a prior pull request as an example as well, for example <a href='https://github.com/meshery/meshery/pull/9880/files'>PR #9880</a>." %}

<a class="anchorjs-link" id="relationship-scopes"></a>

### 4. Configuring the Scope of Relationships

The extent to which a relationship affects components within a model or beyond a model is defined and controlled using scopes. Scopes exist at two levels in Meshery relationships. 

#### Global Scope

Global scope is defined using the `model` and `version` attributes in the relationship definition.

Relationships can be confined to a specific model, a specific model version, or can be allowed to affect all models. The relationship schema has a `model` and `version` attribute which facilitates this control. For example, if the model is specified as `aws-ec2-controller`, the relationship will work for those components that belong to the `aws-ec2-controller` model.

#### Local Scope

Local scope is defined and controlled via `selectors` attributes in the relationship definition.

Relationship selectors refine the scope of applicability the relationship. It is the details included within the Selector that determines whether there is a match and relationship to be formed. These details include which models and components are involved in the relationship and any constraints in its formation. Selector specify to which component(s) the relationship applies or does not apply (think of in terms of the `AND` operators in a query). Selector Sets are used to combine multiple selectors for more granular control over the logic used when matchmaking (establishing a relationship) between components (think of in terms of the `OR` operators in a query).

<!-- @leecalcote: The following needs rewritten using Selector Sets, @MUzairS15 -->

Selectors are structured as an array wherein each entry comprises a `from` (self) property and a `to` (others) property. The `from` and `to` combined with the `allow` and `deny` properties delineate between components involved in a particular relationship. These entries define the constraints necessary for the existence of a relationship, thus scoping a relationship. Each item in the selector uniquely defines a relation between the components listed. i.e. `from` and `to` fields are evaluated within the context of the selector.

Only the components within the same selector relate to each other via 1:many kind of relation between components listed inside the `from` and `to` field. i.e. Each object inside the `from` relates to each item inside the `two` field within a particular selector.

When defining relationships that involve a large number of combinations between `from` and `to`, selectors sets provide a mechanism to organize and manage these relationships. This prevents the need for crafting complex deny attributes and facilitates easier maintenance. Use of selector sets enhances flexibility and reusability in the definition and configuration of relationships among components.

*Note: When defining Hierarchical relationships, remember that the `from` field represents the child component, while the `to` field represents the parent component.*

<details close>
<summary>Relationship Selector Example</summary>

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

The `selector` defined for the relationship between `WasmFilter` and `EnvoyFilter` components (the first item in the array) is entirely independent of the `selector` identified for the relationship between `ConfigMap` and `Deployment` components. This ensures independence in how these components relate to each other while still permitting similar types of relationships.

This example relationship demonstrates how the `WASMFilter` and `EBPFFilter` components identified in `from` relate to other `EnvoyFilter` `WASMPlugin` components identified in `to` selector. Similarly, `ConfigMap` component identified in the `from` selector corresponds to each `Deployment`, `StatefulSet`, and so on component identified in the `to` selector.

</details>

#### Understanding Relationship Policies and their Evaluation

The `evaluationQuery` property in your relationship definition is used to identity the name of the policy to be used by Meshery's evaluation engine. Meshery embeds Open Policy Agent as it's policy engine

**How should you determine the value for `evaluationQuery`**

All relationship definitions are backed by OPA policies and each relationship depends on their `kind`, `type`, and `subType` in order to be properly evaluated. Which evaluation policy or set of policies used during the evaluation moment is defined by the `evaluationQuery` property, which follows the naming convention of combining each of their  `kind`, `type`, and `subType`  properties with an underscore and the word "_relationship", like so: `kind_type_subtype_relationship`.

So, for example, if you are defining or updating a relationship definition with `kind: edge` and `type: network`, the value for the attribute `evaluationQuery` should be `edge_network_relationship`.

Each policy has a set of evaluation rules defined and the `evaluationQuery` attribute corresponds to the main rule defined inside the policy, during the policy eval the results are collected from this rule.

## Postwork

<a id="relationship-testing"></a>

#### Relationship Authoring Best Practices and Considerations

##### General

1. Use camelCasing as the formatting convention.

##### Scoping

1. To configure a relationship to be applied across models, ensure the `model` property for those relationships is set to `*`, to limit the relationships to a specific model, specify the correct `model`(case sensitive).
1. To configure a relationship to be applied across all versions of a particular model, ensure the `version` property for those relationships is set to `*`, to limit the relationships to a specific version of a model, specify the correct model version.
1. Specify `version` property as a regex to ensure relationships are applied to a subset of versions of a model.

##### Actions

1. If a path `mutatedRef/mutatorRef` contains more than one array path then only the first array position can be specified as `\_` for others explicitly mention them as 0
1. Currently `mutatedRef` doesn’t support having an array.

##### Matching

1. Targets of a Relationship can be specific Components or entire Models.
1. Values for propoerties like `kind`, `version`, and `model` are case-sensitive.
1. Absence of a property in the `selector` is interpretted as the wildcard `*`.
   - For example, a selector with `kind: Pod`, `Model: Kubernetes`, and the absence of the `version` property would be interpretted as `version: *`, which
     means that all the versions of the Kubernetes Pod resource (k8s.io/v1/betav2) will match the selector.
1. The `evaluationQuery` property determines the OPA policy to invoke for relationship evaluation, specify the correct rego query.

##### Conflicts

1. Ensure that the `deny` selectors and `allow` selectors do not conflict with each other i.e. relations are not getting overlapped for `allow` and `deny` selectors.
1. In the event of conflicting Relationship Definitions, the union between them is taken.
   - If we have two Relationships, one from (Component A) to (Component B and Component F), and another
     from (Component A) to (Component B and Component C), then it is similar to having a Relationship
     from Component A to Component B, C and F
1. In the event of an overlapping set of complementary Relationship Definitions, Union.
1. In the event of an overlapping set of conflicting Relationship Definitions, no relationship type (Kind) is inherently more important than the next one, so will not be any case of conflict.

<a class="anchorjs-link" id="relationship-contribution"></a>

#### 4. Contribute your relationship to the Meshery project

Submit a pull request to the Meshery repository with your new relationship definition, so that all users can benefit from the relationship(s) you have defined.

Keeping your relationship definition in a separate file allows for easier management and review of the relationship(s) you have defined.

{% include alert.html title="Keeping your custom Relationships private" type="info" content="Alternatively, if you would like to keep the relatioship definition private, you can bundle your relatinship(s) in a custom model, import the custom model into your Meshery deployment. Your private relationship definition(s) will be registered in your Meshery Server's <a href='/concepts/logical/registry'>registry</a> and available for use within your Meshery deployment." %}

For more information refer - [Model - Construct Models in Meshery](https://docs.google.com/document/d/16z5hA8qVfSq885of9LXFUVvfom-hQXr-6oTD_GgoFmk/edit)
