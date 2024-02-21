---
layout: page
title: Contributing to Models
permalink: project/contributing/contributing-models
redirect_from: project/contributing/contributing-models/
abstract: How to contribute to Meshery Models, Components, Relationships, Policies...
language: en
type: project
category: contributing
list: include
---

Meshery's internal object model is designed to provide a consistent and extensible way of capturing and characterizing the resources under Meshery's management and the capabilities Meshery has at its disposal. Meshery Models serve as the unit of packaging for the object models that define a registered capability or a type of managed infrastructure and their relationships, and details specifics of how to manage them.

Models often represent infrastructure and application technologies, however, they are also capable of defining other types of constructs like annotations, like shapes (infrastructure ambiguous components). Models are used to define the capabilities of Meshery.

_See the Meshery [Registry]({{site.baseurl}}/concepts/logical/registry) to learn more._

[![Model Construct Classification]({{ site.baseurl }}/assets/img/meshmodel/meshmodel-architecture.svg)]({{ site.baseurl }}/assets/img/concepts/meshery-models.png)
_Figure: Model Construct Classification_

## Meshery Constructs and their Lifecycle

This section aids in your understanding of the vernacular of Meshery's internal object model and discusses the difference beteween schemas, definitions, declarations, and instances

The lifecycle of Meshery entities (components, relationships, policies) is represented by the following terms, which are used to describe the various stages of a Model's lifecycle:

<!-- Model Definitions are read-only templates that contain instructions for creating a any given infrasture. A Model Definition is a snapshot or blueprint of the configuration, credentials(s) and dependencies required for an application to run.

Depending upon where they are in their lifecycle, Meshery Models can be referred to differently based on their are comprised of a handful of core constructs. -->

#### Schema

**Schema** _(static)_ **: the skeletal structure representing a logical view of the size, shape, characteristics of a construct.**

The schema represents the skeletal structure of a construct and provides a logical view of its size, shape, and characteristics. It defines the expected properties and attributes of the construct. The schema serves as a blueprint or template for creating instances of the construct. It is a static representation that defines the structure and properties but does not contain specific configuration values.

{% include alert.html type="info" title="Schema example" content='<details><summary>Component schema excerpt</summary><pre> {
"$id": "https://schemas.meshery.io/component.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
"description": "Components are the atomic units for designing infrastructure. Learn more at https://docs.meshery.io/concepts/components",
"required": [
"apiVersion",
"kind",
"schema",
"model"
],
"additionalProperties": false,
"type": "object",
"properties": {
"apiVersion": {
"type": "string",
"description": "API Version of the component."
},
"kind": {
"type": "string",
"description": "Kind of the component."
.
.
.

</pre></details> See <a href="https://github.com/meshery/schema">github.com/meshery/schemas</a> for more details.' %}

#### Definition

**Definition** _(static)_ **: An implementation of the Schema containing an outline of the specific attributes of a given, unconfigured construct.**

A definition is an implementation of the schema. It contains specific configurations and values for the construct at hand. The definition provides the actual configuration details for a specific instance of the construct. It is static because it is created based on the schema but does not change once created. The definition is used to instantiate instances of the construct.

{% include alert.html type="info" title="Definition example" content="a generic, unconfigured Kubernetes Pod." %}

#### Declaration

**Declaration** _(static)_ **: - A configured construct with detailed intentions of a given Definition.**

{% include alert.html type="info" title="Declaration example" content="NGINX container as a Kubernetes Pod with port 443 and SSL termination." %}

#### Instance

**Instance** _(dynamic)_ **: A realized construct (deployed/discovered); An instantiation of the declaration.**

An _instance_ represents a realized construct. An _instance_ is a dynamic representation that corresponds to a deployed or discovered instantiation of a _declaration_. An _instance_ is created based on its corresponding _definition_ and represents an actual running or deployed version of the construct within the environment.

{% include alert.html type="info" title="Instance example" content="NGINX-as234z2 pod running in a cluster as a Kubernetes Pod with port 443 and SSL termination." %}

### Generating Models does not require Meshery Server

Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components.

## Contributing to Model Relationships

Relationships within Model play a crucial role in establishing concrete visualisations of efficient data flow between different components of Meshery. These are used to classify the nature of interaction between one or more interconnected Components.

1. Identify the relationship and any specific constraints to be enforced between the two specific components, their models or potentially other components, models, or environmental considerations.
2. Propose a specific visual representation for the relationship.
3. Visual representation examples:

   - [Hierarchical]({{ site.baseurl }}/assets/img/meshmodel/relationships/hierarchical_relationship.png)
   - [Sibling]({{ site.baseurl }}/assets/img/meshmodel/relationships/sibling_relationship.png)
   - [Binding]({{ site.baseurl }}/assets/img/meshmodel/relationships/binding_relationship.png)
   - [Edge]({{ site.baseurl }}/assets/img/meshmodel/relationships/mount_edge_relationship.png)

4. Prospose the appropriate relationship type, using one of the predefined set of relationship types or suggest a new relationship where an existing type does not fit.
5. Create a Relationship Definition (yaml).
6. Create a policy for evaluation of the relationship (rego). See [examples](https://github.com/meshery/meshery/tree/master/server/meshmodel/policies/).
7. Add in Documentation.

### Existing Relationships, their Definitions and their Subtypes

1. `Hierarchical` relationships involve either an ancestral connection of the components i.e. the creation/ deletion of a Component higher up affects the existence of the Components below in the lineage or a connection which involves the inheritence of features from one Component to the other.

- [Parent](https://github.com/meshery/meshery/tree/master/server/meshmodel/kubernetes/relationships/) - A parent-child relationship implies the requirement of the parent component before the child component can be created. For example, a "Namespace" can be a parent of "Pods" within that namespace. The namespace must exist before creating pods within it.
- [Inventory](https://github.com/meshery/meshery/blob/master/server/meshmodel/kubernetes/relationships/) - Wasm filters can inherit features and functionalities from Envoy filters. This can be used to build on existing functionalities provided by Envoy filters and further extend them using Wasm filters. It enables a modular and scalable approach to customize the behavior of the proxy while maintaining a clear hierarchy of features.

1. `Edge` relationships indicate the possibility of traffic flow between two components. They enable communication and interaction between different Components within the system.

- [Mount](https://github.com/meshery/meshery/tree/master/server/meshmodel/kubernetes/relationships/) - This subtype addresses the storage and access possibility between involved components. For example, a "PersistentVolume" can be mounted to a "Pod" to provide persistent storage for the pod's data.
- [Network](https://github.com/meshery/meshery/tree/master/server/meshmodel/kubernetes/relationships/) - This deals with IP addresses and DNS names and provides stable endpoints for communication. For example, a "Service" provides a stable endpoint for accessing multiple replicas of a "Deployment."
- [Firewall](https://github.com/meshery/meshery/tree/master/server/meshmodel/kubernetes/relationships/) - This acts as intermediary for communications which include standard networking protocols like TCP and UDP. It can enforce network policies to control traffic between components. For example, a "NetworkPolicy" can be used to manage the traffic flow between different "Pods" in the cluster.
- [Permission](https://github.com/meshery/meshery/tree/master/server/meshmodel/kubernetes/relationships/) - This defines the permissions for components if they can have a possible relationship with other Components. It ensures that only authorized Components can interact with each other. For example, a "Role" can define permissions for Components to access specific resources.

1. `Sibling` relationships represent connections between components that are at the same hierarchical level or share a common parent. Siblings can have the same or similar functionalities or may interact with each other in specific ways. These relationships facilitate communication and cooperation between components that are in the same group or category.

- [Sibling](https://github.com/meshery/meshery/blob/master/server/meshmodel/kubernetes/relationships/sibling_matchlabels.json)

### Structure of Selectors

Selectors are structured as an array, wherein each entry comprises a 'from(self)' field and a 'to(other)' field (`[from: [{..}], to: [{..}]]`), delineating the components involved in a particular relationship. These entries define the constraints necessary for the existence of said relationship, thus providing scoping within a relationship.
Each item in the selector, uniquely defines a relation between the components listed. i.e. `from` and `to` fields are evaluated within the context of the selector.

Only the components within the same selector relates to each other via 1:many kind of relation between components listed inside the `from` and `to` field. i.e. Each object inside the `from` relates to each item inside the `two` field within a particular selector.

When defining relationships that involve a large number of combinations between `from` and `to`, selectors provide a mechanism to organize and manage these relationships. This prevents the need for crafting complex deny attributes and facilitates easier maintenance.

This arrangement enhances flexibility and reusability in the definition and configuration of relationships among components.

<details closed>
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
"kind" : "WASMPlugin",
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

The `selector` defined for the relationship between `WasmFilter` and `EnvoyFilter` (the first item in the array) is entirely independent from the `selector` defined for the relationship between `ConfigMap` and `Deployment`. This ensures independence in how these components relate to each other while still permitting similar types of relationships.

The above relation shows `WASMFilter` and `EBPFFilter` defined inside `from` relates to each component defined inside `to` `(EnvoyFilter, WASMPlugin...)`.
Similarly, `ConfigMap` defined inside `from` relates to each component defined inside `to` `(Deployment, StatefulSet,...)`

</details>

#### What is `evaluationQuery` attribute and how to determine value for `evaluationQuery` inside relationship definition?

As all relationship definitions are backed by OPA policies and the relationships depending upon their Kind and Subtype needs to be evaluated with respective policies, the policy to invoke for evaluation is determined by the property `evaluationQuery`, which follows the convention as `kind_subtype_relationship`.

```
Eg: If you are defining/updating a relationship definition with kind: Edge and subType: Network, the value for the attribute `evaluationQuery` should be edge_network_relationship.
```

Each policy has set of rules defined and the `evaluationQuery` attribute corresponds to the main rule defined inside the policy, during the policy eval the results are collected from this rule.

### Configuring the scopes of the relationship definitions

The extent to which a relationship affects components within a model or beyond a model is defined and controlled using scopes. Scopes are defined using the `model` and `version` attributes within the relationship schema.

#### Global Scope

Relationships can be confined to specific model, a specific model version, or can be allowed affects all models. The relationship schema has a `model` and `version` attribute which facilitates this control. For example, if the model is specified as `aws-ec2-controller`, the relationship will work for those components which belongs to the `aws-ec2-controller` model.

#### Local Scope

Scope is defined and controlled via the `selectors` [Selectors](#structure-of-selectors) attribute in the relationships.

### Best practices for defining new relationships

1. Ensure that the `deny` selectors and `allow` selectors do not conflict each other i.e. relations are not getting overlapped for `allow` and `deny` selectors.
2. To configure a relationship to be applied across models, ensure `model` property for those relationships are set to `*`, to limit the relationships to specific model, specify correct `model`(case sensitive).
3. To configure a relationship to be applied across all versions of a particular model, ensure `version` property for those relationships are set to `*`, to limit the relationships to specific version of a model, specify correct model version.
4. Support for specifying version property as a regex to ensure relationships are applied to a subset of versions of a model is coming soon.
5. The `evaluationQuery` property determines the OPA policy to invoke for relationship evaluation, specify correct rego query. To understand what query to specify [refer](#what-is-evaluationquery-attribute-and-how-to-determine-value-for-evaluationquery-inside-relationship-definition).
6. If a path `mutatedRef/mutatorRef` contains more than one array path then only first array positin can be sprciifced as \_ for others explicity meniton them as 0
7. Currently `mutatedRef` doesn’t supoort having an aaray

#### Things to keep in mind while defining relationships

1. Targets of a Relationship can be specific Components or entire Models.
1. The values for Kind, Version and Model are case-sensitive.
1. The convention is to use camel-casing for Kind and SubType values.
1. Absence of a field means in the selector means “\*” (or wildcard).
   - If we have a selector with {Kind: Pod, Model: Kubernetes}, the absence of the Version field here
     means that all the versions of the Kubernetes Pod resource (e.g. k8s.io/v1/betav2) will match.
1. In the event of conflicting Relationship Definitions, union between them is taken.
   - If we have two Relationships, one from (Component A) to (Component B and Component F), and another
     from (Component A) to (Component B and Component C), then it is similar to having a Relationship
     from Component A to Component B, C and F
1. In the event of an overlapping set of complementary Relationship Definitions, Union.
1. In the event of an overlapping set of conflicting Relationship Definitions:
   - No relationship type (Kind) is inherently more important than the next one, so will not be any case
     of conflict

## Relationship Types

1. Edge
   1. Network
   1. Firewall
   1. Binding
      1. Mount
      1. Permission
1. Heirarchical
   1. Inventory
   1. Parent

<details open>
<summary>Example Visual Representations</summary>
    <details close><summary>Hierarchical</summary>
    <figure><br><figcaption>Hierarchical - Parent</figcaption>
    <img alt="Hierarchical - Parent: Namespace to other components" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/hierachical_relationship_namespace_others.png"/>
    </figure>
    </details>
    <details close><summary>Sibling</summary>
    <figure><br><figcaption>Hierarchical - Sibling: Matching Label Selectors</figcaption>
    <img alt=Sibling src="{{ site.baseurl }}/assets/img/meshmodel/relationships/sibling_relationship.png"/>
    </figure>
    </details>
    <details close><summary>Binding</summary>
    <figure><br><figcaption>Hierarchical - Binding: Cluster Role with Cluster Role Binding to ConfigMap</figcaption>
    <img alt=Binding src="{{ site.baseurl }}/assets/img/meshmodel/relationships/binding_relationship.png"/>
    </figure>
    </details>
    <details close><summary>Edge</summary>
    <figure><br><figcaption>Edge - Mount</figcaption>
    <img alt="Edge - Mount" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/mount_edge_relationship.png"/>
    </figure>
    <br>
    <figure><figcaption>Edge - Network: Ingress to Service</figcaption>
    <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_ingress_service.png"/>
    <figure><figcaption>Edge - Network: Service to Pod</figcaption>
    <img alt="Edge - Network: Ingress to Service" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_pod.png"/>
    <figure><figcaption>Edge - Network: Service to Service</figcaption>
    <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_service.png"/>
    <figure><figcaption>Edge - Network: Service to Endpoint</figcaption>
    <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_endpoints.png"/>
    <figure><figcaption>Edge - Network: Service to Deployment</figcaption>
    <img alt="Edge - Network" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_deployment.png"/>
    </figure>
    <br>
    <figure><figcaption>Edge - Permission</figcaption>
    <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_role_service.png"/>
    <figure><figcaption>Edge - Permission: Role to Service</figcaption>
    <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_role_pod.png"/>
    <figure><figcaption>Edge - Permission: Role to Deployment</figcaption>
    <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_role_deployment.png"/>
    <figure><figcaption>Edge - Permission: Cluster Role to Pod</figcaption>
    <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_clusterrole_pod.png"/>
    <figure><figcaption>Edge - Permission: Cluster Role to Service</figcaption>
    <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_clusterrole_service.png"/>
    <figure><figcaption>Edge - Permission: Cluster Role to Deployment</figcaption>
    <img alt="Edge - Permission" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_clusterrole_deployment.png"/>
    </figure>
    <br>
    <figure><figcaption>Edge - Network Policy</figcaption>
    <img alt="Edge - Network Policy" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_policy_edge_relationship.png">
    </figure>
    </details>
</details>

For more information refer - [Model - Construct Models in Meshery](https://docs.google.com/document/d/16z5hA8qVfSq885of9LXFUVvfom-hQXr-6oTD_GgoFmk/edit)

