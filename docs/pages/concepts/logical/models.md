---
layout: default
title: Models
permalink: concepts/logical/models
type: concepts
abstract: Meshery uses a set of resource models to define concrete boundaries to ensure extensible and sustainable management.
language: en
list: include
redirect_from:
- concepts/models
---

Models are units of packaging for Meshery's logical object representation. Models (packages) are versioned and bundle any number of components, relationships, policies, connections and credentials. Meshery's internal object model is designed to provide a consistent and extensible way of capturing and characterizing the resources (applications, services, and infrastructure) under Meshery's management. Models serve as an exportable package (OCI-compatible image) to define managed infrastructure, their relationships, and details specifics of how they are to be managed.

## Key aspects of Meshery Models

* **Logical Representation:** Models offer a logical representation of entities under Meshery's management, abstracting away system-specific details.
* **Packaging:** Model packages can be imported and exported as OCI-compatible images, ensuring portability and encapsulation of intellectual property.
* **Components:** Components within a Model represent the capabilities of the underlying platform and are used by operators to build applications.
* **Relationships:** Relationships define interactions and dependencies between components within a model.
* **Policies:** Policies govern the behavior and constraints of components and relationships, ensuring adherence to desired operational practices.
* **Extensibility:** Models are highly extensible, allowing you to define custom components, relationships, and policies. 

{% include alert.html title="Creating your own models" type="light" content="If you would like to create your own, augment existing models, or contribute new models, please refer to the <a href='/project/contributing/contributing-models'>Contributing to Models</a> guide." %}


See [Contributing to Models](/project/contributing/contributing-models).

As a cloud native manager that provides a comprehensive set of tools for managing multi-cloud and cloud native applications and infrastructure, Meshery needs object models that can granularly characterize a wide range of systems under management from simple applications to complex microservices architectures and their infrastructure. This document describes the Meshery Model and its constituent set of entities  used to represent and manage cloud and cloud native systems.

Using Meshery Models, users can define and manage cloud-native systems, including applications, services, and infrastructure. Models provide a way to represent the desired state of the system, and they can be used to deploy, monitor, and manage cloud-native applications. Models are designed to be extensible, allowing users to define new constructs as needed. They are also machine-readable, allowing them to be used by automation tools.

### Model Portability and Your Intellectual Property

Each model (package) can be imported and exported from Meshery Server as OCI-compatible images. This makes models portable, abstracts their complexity, and encapsulates intellectual property that you might have created while designing your own models. Model packages are versioned and bundle any number of components, relationships, policies, connections, and credentials.

Every entity type is implemented by a [registrant]({{site.baseurl}}/concepts/logical/registry). Without registrants and models, Meshery can't manage any kind of infrastructure.

Most registrants configure a specific infrastructure platform (either cloud or self-hosted). Registrants vary in their capabilities. Capabilities come in the form of components, some of which represent infrastructure to be managed others of which represent functional additional functionality to augment Meshery's behavior and deepen its feature set. Some registrants offer infrastructure-specific orchestration. For example the Meshery Adapter for Istio offers integration with each of Istio's addons. Others registrants offer specific services. For example, the Meshery Adapter for Nighthawk offers load generation and service performance characteriazation.

## Design Principles

Meshery Models adhere to several design principles, including establishing a set of constructs, clearly defining construct boundaries, allowing construct extension, reusing existing models, being user-centric, and simplifying complex networking and distributed systems. Meshery Models is designed to meet the following goals:

* **Comprehensive:** The model should be able to represent a wide range of cloud and cloud native resources.
* **Extensible:** The model should be extensible, allowing new constructs to be added as needed.
* **User-centric:** The model should be easy for users to understand and use.
* **Machine-readable:** The model should be machine-readable, allowing it to be used by automation tools.

## Models as the Unit of Packaging

Each model package can be imported and exported from the system as OCI-compatible images, making them portable (a design goal), abstracting their complexity, and encapsulating potential intellectual property that users might have invested into their models. Model packages are versioned and bundle any number of components, relationships, policies, connections and credentials. For example:

{% include alert.html title="Model Packaging" type="info" content="Model constructs can be packaged and exported as OCI-compatible images. This makes them portable and allows them to be shared between different environments." %}

![Meshery Models]({{ site.baseurl }}/assets/img/concepts/meshery-models.png)
_Figure: Model Construct Classification_

## Key aspects and characteristics

You might not fully comprehend the Meshery Models figure above. The following analogy offers an alternative viewpoint from which to comprehend the Capabilities Registry (where Models are imported and their capabilities registered for use).  

![Meshery Models Analogy]({{ site.baseurl }}/assets/img/concepts/meshery-models-analogy.svg)
_Figure: Registrar's Office and Meshery Models Analogy_

### Models

Models introduce various core constructs that form the foundation of the model. Some of the core constructs mentioned in the document include Components, Designs, Policies, and Relationships. Models having the same `name` and `version` attributes are considered duplicates.

### Component

[Component]({{site.baseurl}}/concepts/logical/components) represent entities in the Meshery ecosystem, exposing capabilities of the underlying platform. They can be registered, created, and used by users and operators. Components have definitions, instances, and associated metadata. Components having the same `kind`, `apiVersion` and `model.name` attributes are considered duplicates.

### Policy

[Policies]({{site.baseurl}}/concepts/logical/policies) includes constructs for managing metrics, defining actions, and specifying color properties of components or designs. These constructs help in monitoring, controlling, and visualizing different aspects of the Meshery ecosystem.

### Relationships

[Relationships]({{site.baseurl}}/concepts/logical/relationships) define the nature of interaction between interconnected components in Model. They represent various types of connections and dependencies between components, such as hierarchical, network, or default relationships. Relationships have selectors, metadata, and optional parameters.

#### Evaluation of Relationships

Meshery provides a relationship evaluation algorithm to ensure desired behavior enforcement. [Policies](policies) can be applied to components and relationships, defining rules and actions based on predefined conditions.

{% include alert.html title="Model Schema" type="info" content="Model constructs are defined using a schema language called Cue. Cue is a powerful and expressive language that is well-suited for defining cloud-native constructs." %}

### Designs

[Designs]({{site.baseurl}}/concepts/logical/designs) are deployable units in Meshery that describe the desired infrastructure. They consist of components and patterns, allowing users to define and configure the behavior of their cloud-native applications. A design is a collection of components and patterns that represent a desired state of infrastructure. Designs are used to deploy and manage cloud-native systems.

### Patterns

[Patterns](patterns) are reusable configuration that can be applied to components or designs. They define best practices for configuring and operating cloud-native application functionality. Patterns can be applied to components or designs, and they are read-only.

{% include alert.html title="Metadata" type="info" content="Metadata provide additional details about a component in Meshery. They offer specific functionality or characteristics associated with a component, enhancing its capabilities. Metadata can be attached to components to customize their behavior." %}

## Versioning

![Versioning Models](/assets/img/concepts/logical/versioning-meshery-entities.svg)


## Summary

A Model serves as a unit of packaging for Meshery's logical object representation. Models encompass a structured framework that defines the components, relationships, policies, and other elements that constitute a system or application under Meshery's management. Models offer a way to encapsulate the complexities of diverse systems and provide a standardized representation for management and interaction.
