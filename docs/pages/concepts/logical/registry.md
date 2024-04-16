---
layout: default
title: Registry
permalink: concepts/logical/registry
type: concepts
abstract: Meshery Registry is a database acting as the central repository for all capabilities known to Meshery. These capabilities encompass various entities, including models, components, relationships, and policies.
language: en
list: include

---
# The Meshery Registry: A Central Hub for Capabilities

The Meshery Registry is a critical component acting as the central repository for all capabilities known to Meshery. These capabilities encompass various entities, including:

- [Models](./models): Blueprints defining configurations for interacting with cloud-native infrastructure. They consist of operations, components, relationships, and policies.
  - [Components](../components): Components are reusable building blocks for depicting capabilities defined within models.
  - [Relationships](../relationships): Define connections and dependencies between components within a model.
  - [Policies](../policies): Enforce specific rules and governance for system behavior under Meshery's management.
  - [Connections](../connections): managed and unmanaged resources that Meshery can interact with.
  - [Credentials](../credentials): Optionally, included secrets associated with connections contained in a model.

## Understanding Key Actors

- **Registry**: a component within Meshery that contains a database of known capabilities.
- **Registrar**: The internal Meshery Server process responsible for managing and maintaining the registry.
- **Registrant** (Entity Source): The source of an entity (e.g., model file, Kubernetes cluster).
Entity (Registree): An individual capability stored in the registry (e.g., model, component).
- **Registrant** *(Host)*: A Meshery Cnnection responsible for sourcing and registering entities. A registrant can perform registration for their own entities, or a registrant can act as a proxy on behalf of a third-party entity source.
- **Entity** *(registree)* - an entry in the Meshery Registry; e.g. a model, component, relationship, or policy. Sometimes referred to as a capability.
<!-- - **Entity Source**: an entity’s original location from which it was sourced; e.g. (source_uri is used as the flag by Meshery Server to assess whether additional support). The Entity Source should have all the information that Meshery needs to generate the components.   -->

### Model Dynamics: Static vs. Dynamic

Each Meshery release comes with a built-in set of models automatically registered during server startup. These models offer a core set of entities for Meshery's supported [integrations](/extensibility/integrations).

**Static Models:** Pre-defined models included with each Meshery release.
**Dynamic Models:** Models generated at runtime by connecting to supported platforms like Kubernetes clusters or cloud providers.

## Lifecycle of Entities in Registry

### Registration

The process of registering an entity in the Meshery Registry is a multi-step process. The process begins with the sourcing of the entity from the Registrant's source. The entity is then registered in the Meshery Registry. The entity is then available for use within the Meshery ecosystem.

Meshery [Adapters]({{ site.baseurl }}/concepts/architecture/adapters) are one example of a Registrant. Registrants are responsible for the registration of entities in the Meshery Registry. Adapters are responsible for the sourcing and registration of entities and the packaging of these enties into one or more models.

#### Using Meshery CLI to Register a Model

```bash
$ mesheryctl model import -f <path-to-model>
```

#### Using Meshery UI to Register a Model

Visit the Settings --> Registry page and click the "Import" button to import a model.

{% include alert.html type="tip" title="Dynamic and Static Models" content="Models can be either dynamic or static. Dynamic models are generated at runtime upon connecting to a supported platform, like a Kubernetes cluster or cloud provider. Static models are pre-defined and are included in each Meshery release." %}

### Ignoring an Entity

You have control over whether a registered entity (model and all that the model contains) this can be an individual or team-level preference. Use the "Ignore" action to designate whether a given model is allowed to be used within a given Meshery Server deployment. Models that are ignored remain in the Meshery Registry but are not available for use within a given Meshery Server deployment.

## Generating Models

Generating a model does not require Meshery Server. The model can be generated from the source_uri. The source_uri is used as the flag by the pattern engine to assess whether additional support is needed. The Entity Source should have all the information that Meshery needs to generate and deploy the components.

