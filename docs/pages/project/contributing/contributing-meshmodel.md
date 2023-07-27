---
layout: page
title: Contributing to MeshModel
permalink: project/contributing/contributing-meshmodel
redirect_from: project/contributing/contributing-meshmodel/
description: How to contribute to MeshModel
language: en
type: project
category: contributing
---

MeshModel serves as a foundational element in the Meshery ecosystem, representing the interconnectedness of various components inside the infrastructure of an architecture.

In the context of MeshModel, the `Core Constructs` of MeshModel are represented in three forms:

1.  `Schema` (static): The schema represents the skeletal structure of a construct and provides a logical view of its size, shape, and characteristics. It defines the expected properties and attributes of the construct. The schema serves as a blueprint or template for creating instances of the construct. It is a static representation that defines the structure and properties but does not contain specific configuration values.

2.  `Definition` (static): The definition is an implementation of the schema. It contains specific configurations and values for the construct at hand. The definition provides the actual configuration details for a specific instance of the construct. It is static because it is created based on the schema but does not change once created. The definition is used to instantiate instances of the construct.

        Things to Keep in Mind while creating RelationshipDefinitions:

        a. Relationships are defined only between Components.
        b. The values for Kind, Version  and Model are case-sensitive
        c. The convention is to use camel-casing for Kind and SubType values.
        d. Absence of a field means in the selector means “*” (or wildcard).
            - If we have a selector with {Kind: Pod, Model: Kubernetes}, the absence of the Version field here
              means that all the versions of the Kubernetes Pod resource (e.g. k8s.io/v1/betav2) will match.
        e. In the event of conflicting Relationship Definitions, union between them is taken.
            - If we have two Relationships, one from (Component A) to (Component B and Component F), and another
              from (Component A) to (Component B and Component C), then it is similar to having a Relationship
              from Component A to Component B, C and F
        f. In the event of an overlapping set of complementary Relationship Definitions, Union.
        g. In the event of an overlapping set of conflicting  Relationship Definitions:
            - No relationship type (Kind) is inherently more important than the next one, so will not be any case
              of conflict

3.  `Instance` (dynamic): The instance represents a realized construct. It is a dynamic representation that corresponds to a deployed or discovered instance of the construct. The instance is created based on the definition and represents an actual running or deployed version of the construct within the service mesh environment.

If a specific attribute is not provided with a value in the definition, it means that the value for that attribute has to be written or configured per construct. In other words, the absence of a value indicates that the configuration for that attribute is required and specific to each individual construct instance.

[![MeshModel Contruct Classification]({{ site.baseurl }}/assets/img/meshmodel/meshmodel-architecture.svg)]({{ site.baseurl }}/assets/img/meshmodel/meshmodel-architecture.svg)

_Figure: MeshModel Contruct Classification_

### Contribute to MeshModel Relationships

Relationships within MeshModel play a crucial role in establishing concrete visualisations of efficient data flow between different components of Meshery. These are used to classify the nature of interaction between one or more interconnected Components.

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

<<<<<<< HEAD
**Existing Relationships, their Definitions and their Subtypes**

1. `Hierarchical` relationships involve either an ancestral connection of the components i.e. the creation/ deletion of a Component higher up affects the existence of the Components below in the lineage or a connection which involves the inheritence of features from one Component to the other.  
- [Parent](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/hierarchical_parent.json) - A parent-child relationship implies the requirement of the parent component before the child component can be created. For example, a "Namespace" can be a parent of "Pods" within that namespace. The namespace must exist before creating pods within it.
- [Inventory](https://github.com/meshery/meshery/blob/master/server/meshmodel/relationships/hierarchical_inv_wasm_filters.json) - Wasm filters can inherit features and functionalities from Envoy filters. This can be used to build on existing functionalities provided by Envoy filters and further extend them using Wasm filters. It enables a modular and scalable approach to customize the behavior of the proxy while maintaining a clear hierarchy of features.

2. `Edge` relationships indicate the possibility of traffic flow between two components. They enable communication and interaction between different Components within the system.
- [Mount](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/mount_edge.json) - This subtype addresses the storage and access possibility between involved components. For example, a "PersistentVolume" can be mounted to a "Pod" to provide persistent storage for the pod's data.
- [Network](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/network_edge.json) - This deals with IP addresses and DNS names and provides stable endpoints for communication. For example, a "Service" provides a stable endpoint for accessing multiple replicas of a "Deployment."
- [Firewall](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/network_policy_edge.json) - This acts as intermediary for communications which include standard networking protocols like TCP and UDP. It can enforce network policies to control traffic between components. For example, a "NetworkPolicy" can be used to manage the traffic flow between different "Pods" in the cluster.
- [Permission](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/permission_edge.json) - This defines the permissions for components if they can have a possible relationship with other Components. It ensures that only authorized Components can interact with each other. For example, a "Role" can define permissions for Components to access specific resources.

3. `Sibling` relationships represent connections between components that are at the same hierarchical level or share a common parent. Siblings can have the same or similar functionalities or may interact with each other in specific ways. These relationships facilitate communication and cooperation between components that are in the same group or category.
=======
**Existing Relationships and their Subtypes**

1. Hierarchical
- [Inventory](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/hierarchical_inv_wasm_filters.json)
- [Parent](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/hierarchical_parent.json)

2. Edge
- [Mount](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/mount_edge.json)
- [Network](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/network_edge.json)
- [Firewall](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/network_policy_edge.json)
- [Permission](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/permission_edge.json)

3. Sibling
>>>>>>> 8ab616978 (add screenshots in contrib docs)
- [Sibling](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/sibling.json)


<details open>
<summary>See all Visual Representations</summary>
    <details close><summary>Hierarchical</summary>
    <figure><br><figcaption>Hierarchical Parent</figcaption>
<<<<<<< HEAD
    <img alt="Hierarchical Parent" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/hierachical_relationship_namespace_others.png"/>
=======
    <img alt="Hierarchical Parent" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/hierarchical_relationship.png"/>
>>>>>>> 8ab616978 (add screenshots in contrib docs)
    </figure>
    </details>

    <details close><summary>Sibling</summary>
    <figure><br><figcaption>Sibling</figcaption>
    <img alt=Sibling src="{{ site.baseurl }}/assets/img/meshmodel/relationships/sibling_relationship.png"/>
    </figure>
    </details>

    <details close><summary>Binding</summary>
    <figure><br><figcaption>Binding</figcaption>
    <img alt=Binding src="{{ site.baseurl }}/assets/img/meshmodel/relationships/binding_relationship.png"/>
    </figure>
    </details>

    <details close><summary>Edge</summary>
    <figure><br><figcaption>Mount Edge</figcaption>
    <img alt="Mount Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/mount_edge_relationship.png"/>
    </figure>

    <br>
    <figure><figcaption>Network Edge</figcaption>
    <img alt="Network Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_ingress_service.png"/>
    <img alt="Network Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_pod.png"/>
    <img alt="Network Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_service.png"/>
    <img alt="Network Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_endpoints.png"/>
    <img alt="Network Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_edge_relationship_service_deployment.png"/>
    </figure>

    <br>
    <figure><figcaption>Permission Edge</figcaption>
    <img alt="Permission Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_role_service.png"/>
    <img alt="Permission Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_role_pod.png"/>
    <img alt="Permission Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_role_deployment.png"/>
    <img alt="Permission Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_clusterrole_pod.png"/>
    <img alt="Permission Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_clusterrole_service.png"/>
    <img alt="Permission Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/permission_edge_relationship_clusterrole_deployment.png"/>
    </figure>

    <br>
    <figure><figcaption>Network Policy Edge</figcaption>
    <img alt="Network Policy Edge" src="{{ site.baseurl }}/assets/img/meshmodel/relationships/network_policy_edge_relationship.png">
    </figure>

    </details>

</details>

For more information refer - [MeshModel - Construct Models in Meshery](https://docs.google.com/document/d/16z5hA8qVfSq885of9LXFUVvfom-hQXr-6oTD_GgoFmk/edit)

