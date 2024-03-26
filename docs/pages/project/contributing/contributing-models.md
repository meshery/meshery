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

Model and Component generation logic is MeshKit. `mesheryctl` and Meshery Server both utilize MeshKit’s libraries for ongoing programmatic generation of models and components. For adding a model, the link to the CRDs for the specific model needs to be added to the Meshery Integration Spreadsheet. On adding the link to the Meshery Integration Spreadsheet, the model generator automatically registers the specific model with Meshery.
