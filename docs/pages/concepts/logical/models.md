---
layout: default
title: Meshery Models
permalink: concepts/models
type: concepts
abstract: Meshery uses a set of resource models to define concrete boundaries to ensure extensible and sustainable management.
language: en
list: include
---

Meshery's internal object model is designed to provide a consistent and extensible way of capturing and characterizing the resources under Meshery's management. Meshery Models serve as an internal model within Meshery, defining constructs and their relationships, and enabling consistent communication between components.

![Meshery Models]({{ site.baseurl }}/assets/img/concepts/meshery-models.png)

## Key aspects and characteristics

You might not fully comprehend the Meshery Models figure above. The following analogy offers an alternative viewpoint from which to comprehend the Capabilities Registry (where Models are imported and their capabilities registered for use).  

![Meshery Models Analogy]({{ site.baseurl }}/assets/img/concepts/meshery-models-analogy.svg)

## Design Principles

Meshery Models adheres to several design principles, including establishing a set of constructs, clearly defining construct boundaries, allowing construct extension, reusing existing models, being user-centric, and simplifying complex networking and distributed systems.

### Models

[Models](https://github.com/meshery/meshery/tree/master/server/meshmodel) introduce various core constructs that form the foundation of the model. Some of the core constructs mentioned in the document include Components, Designs, Policies, and Relationships. Models having the same `name` and `version` attributes are considered duplicates.

#### Models as the Unit of Packaging

MeshModel supports packaging constructs as OCI-compatible images, making them portable and encapsulating intellectual property. Model packages can include multiple MeshModel constructs, facilitating reusability and versioning.

### Component

[Components](components) represent entities in the Meshery ecosystem, exposing capabilities of the underlying platform. They can be registered, created, and used by users and operators. Components have definitions, instances, and associated metadata. Components having the same `kind`, `apiVersion` and `model.name` attributes are considered duplicates.See [Relationships in GitHub](https://github.com/meshery/meshery/tree/master/server/meshmodel/components) for more information.

### Designs

[Designs](designs) are deployable units in Meshery that describe the desired infrastructure. They consist of components and patterns, allowing users to define and configure the behavior of their cloud-native applications.

### Patterns

Patterns augment the operational behavior of a deployed instance of a design. They define best practices for configuring and operating cloud-native application functionality. Patterns can be applied to components or designs, and they are read-only.

### Metadata

Metadata provide additional details about a component in Meshery. They offer specific functionality or characteristics associated with a component, enhancing its capabilities. Metadata can be attached to components to customize their behavior.

### Policy

[Policy](policy) includes constructs for managing metrics, defining actions, and specifying color properties of components or designs. These constructs help in monitoring, controlling, and visualizing different aspects of the Meshery ecosystem. See [Policies in GitHub](https://github.com/meshery/meshery/tree/master/server/meshmodel/policies) for more information.

### Relationships

[Relationships](relationships) define the nature of interaction between interconnected components in MeshModel. They represent various types of connections and dependencies between components, such as hierarchical, network, or default relationships. Relationships have selectors, metadata, and optional parameters. See [Relationships in GitHub](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships) for more information.

#### Evaluation of Relationships

Meshery provides a relationship evaluation algorithm to ensure desired behavior enforcement. [Policies]() can be applied to components and relationships, defining rules and actions based on predefined conditions.



