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
2. What policies are loaded by default?
3. What happens in conflict?
4. What controls are exposed to model contributors?
5. Are there any Global meshery defaults (can user change them?) -->

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

</pre></details> See <a href="https://github.com/meshery/schemas">github.com/meshery/schemas</a> for more details.' %}

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

# Instructions for Creating a New Model

All of Meshery's Models, Components, and Relationships can be found in the Meshery Integrations spreadsheet. This spreadsheet is the source of truth for the definition of Meshery's models. On a daily schedule, the contents of the Meshery Integrations spreadsheet is 
{% include alert.html type="light" title="Model Source Code" content="See examples of <a href='https://github.com/meshery/meshery/tree/master/server/meshmodel'>Models defined in JSON in meshery/meshery</a>." %} 

To add or update a model, follow these steps:

1. **Create a Model Definition.** Open the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>. Create a new row (or comment to suggest a new row) to capture the specific details of your model. As you fill-in model details, referernce each column's notes and comments as instructions and an explanation of their purpose.
2. **Generate Components.** Once you have entered values into the required columns, execute the following command to generate components for your model.

{% capture code_content %}$ mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred “${{SPREADSHEET_CRED}}"{% endcapture %}
 {% include code.html code=code_content %}

1. **Enhance Component details.** While the default shape for new components is a circle, each component should be considered for its best-fit shape.
  1. Review and familiarize with the available set of predefined relationship types. Refer to the Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for a list of possible shapes. 
  2. Propose a specific shape, best-suited to visually represent the Component. _Example - Deployment as a pentagon._
  3. Proposee a specific icon, best-suited to visually represent the Component. _Example - DaemonSet as a skull icon._

{% include alert.html type="info" title="Using Meshery CLI with the Registry (models)" content="Create new and list existing models by using <code>mesheryctl registry</code> to interact with the Meshery Registry and the <a href='https://docs.google.com/spreadsheets/d/1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw/edit#'>Meshery Integrations spreadsheet</a>." %} 

### Instructions for Relationships

See the [Contributing to Relationships]({{site.baseurl}}/project/contributing/contributing-relationships) page.

1. Identify the relationship and any specific constraints to be enforced between one or more specific components within the same or different models.
1. Propose a specific visual representation for the relationship. See list of visualizations on [Visualizing Relationships](https://docs.meshery.io/project/contributing/contributing-relationships#relationship-visualizations)
1. Prospose the appropriate relationship type, using one of the predefined set of relationship types or suggest a new relationship where an existing type does not fit.
1. Create a Relationship Definition (yaml).
1. (Rarely necessary) Create a policy for evaluation of the relationship (rego). _This step is only necessary and can typically be skipped. Contact a maintainer if the relationship requires a new policy to evaluate the relationship._
1. Review a prior pull request as an example of how to define a Relationships. For example, see [PR #9880](https://github.com/meshery/meshery/pull/9880/files)

### Generating Models does not require Meshery Server

Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components. For adding a model, the link to the CRDs for the specific model needs to be added to the Meshery Integration Spreadsheet. On adding the link to the Meshery Integration Spreadsheet, the model generator automatically registers the specific model with Meshery.
