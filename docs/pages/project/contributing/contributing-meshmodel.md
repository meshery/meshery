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

1. `Schema` (static): The schema represents the skeletal structure of a construct and provides a logical view of its size, shape, and characteristics. It defines the expected properties and attributes of the construct. The schema serves as a blueprint or template for creating instances of the construct. It is a static representation that defines the structure and properties but does not contain specific configuration values.

2. `Definition` (static): The definition is an implementation of the schema. It contains specific configurations and values for the construct at hand. The definition provides the actual configuration details for a specific instance of the construct. It is static because it is created based on the schema but does not change once created. The definition is used to instantiate instances of the construct.

3. `Instance` (dynamic): The instance represents a realized construct. It is a dynamic representation that corresponds to a deployed or discovered instance of the construct. The instance is created based on the definition and represents an actual running or deployed version of the construct within the service mesh environment.

If a specific attribute is not provided with a value in the definition, it means that the value for that attribute has to be written or configured per construct. In other words, the absence of a value indicates that the configuration for that attribute is required and specific to each individual construct instance.

[![MeshModel Contruct Classification_](https://raw.githubusercontent.com/meshery/meshery/master/docs/assets/img/meshmodel/meshmodel-architecture.svg)](https://raw.githubusercontent.com/meshery/meshery/master/docs/assets/img/meshmodel/meshmodel-architecture.svg)

_Figure: MeshModel Contruct Classification_

### Contribute to MeshModel Relationships

Relationships within MeshModel play a crucial role in establishing concrete visualisations of efficient data flow between different components of Meshery. These are used to classify the nature of interaction between one or more interconnected Components.

1. Identify the relationship and any specific constraints to be enforced between the two specific components, their models or potentially other components, models, or environmental considerations.
2. Propose a specific visual representation for the relationship. 
3. Visual representation examples:
    - [Hierarchical](https://github.com/meshery/meshery/tree/master/.github/assets/images/hierarchical_relationship.png)
    - [Sibling](https://github.com/meshery/meshery/tree/master/.github/assets/images/sibling_relationship.png)
    - [Binding](https://github.com/meshery/meshery/tree/master/.github/assets/images/binding_realtionship.png)
4. Prospose the appropriate relationship type, using one of the predefined set of relationship types or suggest a new relationship where an existing type does not fit.
5. Create a Relationship Definition (yaml).
6. Create a policy for evaluation of the relationship (rego). See [examples](https://github.com/meshery/meshery/tree/master/server/meshmodel/policies/).
7. Add in Documentation.

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
- [Sibling](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships/sibling.json)


For more information refer - [MeshModel - Construct Models in Meshery](https://docs.google.com/document/d/16z5hA8qVfSq885of9LXFUVvfom-hQXr-6oTD_GgoFmk/edit)