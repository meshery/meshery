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

<!-- Concepts for which docs needs to be updated: -->
<!-- Scopes - What they mean to contributors/expand on which takes precedence?
1. Which policies get loaded?
   1. What policies are loaded by default?
   2. What happens in conflict?
2. What controls are exposed to model contributors?
3. Are there any global Meshery defaults (can user change them?) 
4. Instructions for Creating a New Connection
5. Instructions for Creating a New Component -->

# Understanding the internals of Meshery's logical object model

Meshery's internal object model is designed to provide a consistent and extensible way of capturing and characterizing the resources under Meshery's management and the capabilities Meshery has at its disposal. Meshery Models serve as the unit of packaging for the object models that define a registered capability or a type of managed infrastructure and their relationships, and details specifics of how to manage them.Models often represent infrastructure and application technologies, however, they are also capable of defining other types of entities like annotations, like shapes (infrastructure ambiguous components). Models are used to define the capabilities of Meshery. _See [Models]({{site.baseurl}}/concepts/logical/models) to learn more about models as a logical concept._

Each model includes a set of entities (in the form of definitions) that Meshery can manage. Models are defined and versioned using on the [Model Schema](https://github.com/meshery/schemas/blob/master/schemas/constructs/openapi/meshmodels.yml). The schema defines the structure of the model, including the entities it contains, their relationships, and the properties they have. The schema also defines the version of the model and the version of the schema itself. _See [Registry]({{site.baseurl}}/concepts/logical/registry) to learn more about Meshery's internal registry and how to use it._

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

# Instructions for creating a new Model

All of Meshery's Models, Components, and Relationships can be found in the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>. This spreadsheet is the source of truth for the definition of Meshery's models. On a daily schedule, the contents of the Meshery Integrations spreadsheet is refreshed.

{% include alert.html type="light" title="Model Source Code" content="See examples of <a href='https://github.com/meshery/meshery/tree/master/server/meshmodel'>Models defined in JSON in meshery/meshery</a>." %}

To add or update a model, follow these steps:

1. **Create a Model Definition.** Open the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>. Create a new row (or comment to suggest a new row) to capture the specific details of your model. As you fill-in model details, referernce each column's notes and comments as instructions and an explanation of their purpose.
1. **Generate Components.** Once you have entered values into the required columns, either execute step 2.a. or 2.b.
   1. Execute the following command to generate components for your model.
{% capture code_content %}$ mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred “${{SPREADSHEET_CRED}}"{% endcapture %}
 {% include code.html code=code_content %}
   1. Ask a maintainer to invoke the [Model Generator workflow](https://github.com/meshery/meshery/actions/workflows/model-generator.yml).
1. **Enhance Component details.** While the default shape for new components is a circle, each component should be considered for its best-fit shape.
1. Review and familiarize with the available set of predefined relationship types. Refer to the Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for a list of possible shapes.
1. Propose a specific shape, best-suited to visually represent the Component. _Example - Deployment as a pentagon._
1. Proposee a specific icon, best-suited to visually represent the Component. _Example - DaemonSet as a skull icon._

{% include alert.html type="info" title="Using Meshery CLI with the Registry (models)" content="Create new and list existing models by using <code>mesheryctl registry</code> to interact with the Meshery Registry and the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>." %}

### Instructions for creating a new Component

See the [Contributing to Components]({{site.baseurl}}/project/contributing/contributing-components) for detailed instructions.

### Instructions for creating a new Relationship

See the [Contributing to Relationships]({{site.baseurl}}/project/contributing/contributing-relationships) for detailed instructions.

{% include alert.html type="info" title="Generating Models does not require Meshery Server" content="Meshery Server is not required to generate models. The Meshery CLI can be used to generate models. Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components." %}


<!-- ### Instructions for Creating a New Connection

### Managed and Unmanaged Connections

Each Meshery Model can contain one more ConnectionDefinitions (files), each Definition representing one Connection, and also, (as a matter of convenience multiple Connections can be described in the same ConnectionDefinition file).

Connections can be:

1. a ConnectionDefinition based Meshery's [Connection Schema](https://github.com/meshery/schemas/) with hand-curated Connection attributes.
2. a custom ConnectionDefinition based Meshery's Connection Schema that references an existing Component within the same Model. -->

## Next Steps

The Meshery team is currently working on the following:

* Extending the model to support additional entities
* Improving the tooling for working with models
* Defining relationships between components and embedding those policies within models

We encourage you to get involved in the development of Meshery Models and to share your feedback.
  
  {% include alert.html type="info" title="Meshery Models are extensible" content="Meshery Models are designed to be extensible, allowing you to define new components as needed. If you have an idea for a new component, please create one and share it with the Meshery community." %}
