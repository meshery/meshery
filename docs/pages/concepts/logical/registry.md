---
layout: enhanced
title: Registry
permalink: concepts/logical/registry
type: concepts
abstract: Registry is a component within Meshery that contains a database of known capabilities.
language: en
list: include

---
The Meshery Registry is a component within Meshery that contains a database of known capabilities. Capabilities include registered models that many contains operations, components, relationships, policies, and more. The Registry is a key component of Meshery's core functionality.

Each Meshery release includes a predefined set of [models](./models) that are registered upon boot of Meshery Server in the Registry. The Meshery Registry is a key component of Meshery's core functionality.

## Glossary of Core Concepts

- **Registry** - a component within Meshery that contains a database of known capabilities.
- **Registrar** - the Meshery Server process responsible for the process of record keeping.
- **Entity** *(registree)* - an capability (an entity) that is found in the Meshery Registry; e.g. a model, component, relationship, policy.Sometimes referred to as a capability.
- **Entity Source** - an entity’s original location from which it was sourced; e.g. (source_uri is used as the flag by th pattern engine to assess whether additional support). The Entity Source should have all the information that Meshery needs to generate the components.  
- **Registrant** *(host)* - A Connection responsible for the sourcing and registration of entities. Sometimes Registrants perform registration for their own components, while other times they act on behalf of the component source.

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

