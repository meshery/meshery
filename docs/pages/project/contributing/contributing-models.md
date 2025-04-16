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

### Understanding the internals of Meshery's logical object model

Meshery uses a logical object model to describe the infrastructure and capabilities it manages in a consistent and extensible way.

At the core of this system are **Meshery Models** — packages that define a specific type of infrastructure, application, or capability. These models include:

- **Components**: Individual parts of a system (e.g., services, databases).
- **Relationships**: How those parts interact.
- **Metadata**: Visual and behavioral traits, such as icons or capabilities.

Models can describe traditional technologies (like Kubernetes workloads), or more abstract entities (like annotations or diagrams).
> Learn more: [What are Meshery Models?]({{site.baseurl}}/concepts/logical/models)

Each model includes a set of entities (in the form of definitions) that Meshery can manage. Models are defined and versioned using on the [Model Schema](https://github.com/meshery/schemas/blob/master/schemas/constructs/openapi/meshmodels.yml). 

The schema defines the structure of the model, including the entities it contains, their relationships, and the properties they have. The schema also defines the version of the model and the version of the schema itself.

> See [Registry]({{site.baseurl}}/concepts/logical/registry) to learn more about Meshery's internal registry and how to use it.

[![Model Entity Classification]({{ site.baseurl }}/assets/img/meshmodel/meshmodel-architecture.svg)]({{ site.baseurl }}/assets/img/concepts/meshery-models.png)
_Figure: Model Entity Classification_

## Meshery Entities and their Lifecycle

This section aids in your understanding of the vernacular of Meshery's internal object model and discusses the difference beteween schemas, definitions, declarations, and instances. The lifecycle of Meshery entities (components, relationships, policies) is represented by the following terms, which are used to describe the various stages of their lifecycle.

### Schema

**Schema** _(static)_ **: the skeletal structure representing a logical view of the size, shape, characteristics of a construct.**

The schema represents the skeletal structure of an entity and provides a logical view of its size, shape, and characteristics. It defines the expected properties and attributes of the entity. The schema serves as a blueprint or template for creating instances of the entity. It is a static representation that defines the structure and properties but does not contain specific configuration values.

{% include alert.html type="info" title="Schema example" content='<details><summary>Component schema excerpt</summary><pre> {
"$id": "https://schemas.meshery.io/component.json",
  "$schema": "<http://json-schema.org/draft-07/schema#>",
"description": "Components are the atomic units for designing infrastructure. Learn more at <https://docs.meshery.io/concepts/components>",
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

</pre></details> See <a href="https://github.com/meshery/schemas">github.com/meshery/schemas</a> for more details.' %}

### Definition

**Definition** _(static)_ **: An implementation of the Schema containing an outline of the specific attributes of a given, unconfigured entity.**

A definition is an implementation of the schema. It contains specific configurations and values for the entity at hand. The definition provides the actual configuration details for a specific instance of the entity. It is static because it is created based on the schema but does not change once created. The definition is used to instantiate declarations of the entity.

{% include alert.html type="info" title="Definition example" content="a generic, unconfigured Kubernetes Pod." %}

### Declaration

**Declaration** _(static)_ **: - A configured entity with detailed intentions of a given Definition.**

{% include alert.html type="info" title="Declaration example" content="NGINX container as a Kubernetes Pod with port 443 and SSL termination." %}

### Instance

**Instance** _(dynamic)_ **: A realized entity (deployed/discovered); An instantiation of the declaration.**

An _instance_ represents a realized entity. An _instance_ is a dynamic representation that corresponds to a deployed or discovered instantiation of a _declaration_. An _instance_ is created based on its corresponding _definition_ and represents an actual running or deployed version of the entity within the environment.

{% include alert.html type="info" title="Instance example" content="NGINX-as234z2 pod running in a cluster as a Kubernetes Pod with port 443 and SSL termination." %}

<!-- ### Importing Generated Models Here -->

### Importing Generated Models

The generated model can be importing using both Mesheryctl and Meshery UI. Read [Importing Models]({{site.baseurl}}/guides/configuration-management/importing-models) for detailed instructions on how to import models.

### Post Model Generation

During model generation, corresponding components are created. Next step is to enrich these component details and define their capabilities and relationships.

1. **Enrich Component Details**
   When a Component is initially generated, a new Component definition is created with default properties (e.g. colors, icons, capabilities, etc.), some of which are inherited from their respective Model.

   - **1.1. Customize Shapes and Colors**

     - Default shape for new components is a circle
     - Consider enriching components' details based on what they represent
     - Reference Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for possible shapes
     - Example: Use a pentagon shape to represent a Deployment
     - Know more about [components shapes and colors](https://docs.meshery.io/extensions/component-shape-guide)

   - **1.2. Customize Icons**

     - Components inherit the icon (colored and white SVGs) of their respective Model by default
     - Propose specific icons best suited to visually represent each component
     - Example: Use a skull icon for a DaemonSet

   - **1.3. Review Capabilities**
     - Review and confirm assigned capabilities
     - Modify capabilities as needed

    See the [Contributing to Components]({{site.baseurl}}/project/contributing/contributing-components) for detailed instructions.

2. **Identify Relationships**

   - **2.1. Review Available Types**
     Review and familiarize yourself with the predefined relationship kinds, types, and subtypes. See ["Relationships logical concepts"]({{ site.baseurl }}/concepts/logical/relationships)

   - **2.2. Map Component Relationships**

     - Identify appropriate relationships for your new components
     - Consider how components relate to others within the same model
     - Consider relationships with components in other models

   - **2.3. Create Definitions**
     Codify the relationships you have identified into a Relationship Definition

    See the [Contributing to Relationships]({{site.baseurl}}/project/contributing/contributing-relationships) for detailed instructions.

## Next Steps

The Meshery team is currently working on the following:

- Extending the model to support additional entities
- Improving the tooling for working with models
- Defining relationships between components and embedding those policies within models

We encourage you to get involved in the development of Meshery Models and to share your feedback.

{% include alert.html type="info" title="Meshery Models are extensible" content="Meshery Models are designed to be extensible, allowing you to define new components as needed. If you have an idea for a new component, please create one and share it with the Meshery community." %}
