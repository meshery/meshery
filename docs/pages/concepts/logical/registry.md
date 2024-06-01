---
layout: default
title: Registry
permalink: concepts/logical/registry
type: concepts
abstract: Meshery Registry is a database acting as the central repository for all capabilities known to Meshery. These capabilities encompass various entities, including models, components, relationships, and policies.
language: en
display-title: "false"
list: include

---
# Meshery Registry: A Central Hub for Capabilities

The Meshery Registry is a vital component within Meshery, serving as a centralized repository for managing a diverse range of cloud and cloud native resources. It stores and organizes crucial information such as models, categories, components, and relationships, enabling efficient interaction and utilization of these resources within the Meshery ecosystem. You can conveniently [access and manage registry data](#interacting-with-the-meshery-registry) through Meshery UI, and through Meshery CLI ([mesheryctl registry]({{site.baseurl}}/reference/mesheryctl/#meshery-registry-management)).

As the central repository for all capabilities known to Meshery, contains various entities.

<details>
  <summary>Contents of the Registry</summary>
  <br /><br />
  <a href="../models">Models</a>: Blueprints defining configurations for interacting with cloud-native infrastructure. They consist of operations, components, relationships, and policies.
  <ul>
    <li><a href="../components">Components</a>: Reusable building blocks for depicting capabilities defined within models.</li>
    <li><a href="../relationships">Relationships</a>: Define the nature of connections between components within a model, describing how they interact and depend on each other.</li>
    <li><a href="../policies">Policies</a>: Enforce specific rules and governance for system behavior under Meshery's management.</li>
    <li><a href="../connections">Connections</a>: Managed and unmanaged resources that Meshery can interact with.</li>
    <li><a href="../credentials">Credentials</a>: Optionally, included secrets associated with connections contained in a model.</li>
  </ul>
  <br />
</details>

## Key Concepts and Terminology

- **Registry**: a component within Meshery that contains a database of known capabilities.
- **Registrar**: The internal Meshery Server process responsible for managing and maintaining the registry.
- **Registrant** (Entity Source): The source of an entity (e.g., model file, Kubernetes cluster).
Entity (Registree): An individual capability stored in the registry (e.g., model, component).
- **Registrant** *(Host)*: A Meshery Connection responsible for sourcing and registering entities. A registrant can perform registration for their own entities, or a registrant can act as a proxy on behalf of a third-party entity source.
- **Entity** *(registree)* - an entry in the Meshery Registry; e.g. a model, component, relationship, or policy. Sometimes referred to as a capability.
<!-- - **Entity Source**: an entity’s original location from which it was sourced; e.g. (source_uri is used as the flag by Meshery Server to assess whether additional support). The Entity Source should have all the information that Meshery needs to generate the components.   -->

## Models in the Registry

You will find two types of models in the registry: Static and Dynamic.

- **Static Models:** Pre-defined models included with each Meshery release. See the full list of static models.
- **Dynamic Models:** Generated at run-time by connecting Meshery to supported platforms like Kubernetes clusters or cloud providers.

Each Meshery release comes with a built-in set of models automatically registered at Meshery Server boot-time. These built-in models offer a core set of entities for Meshery's supported [integrations](/extensibility/integrations). Once Meshery Server is running, and as it connects to and discovers your infrastructure, *dynamic models* are automatically generated. A given Meshery release may not include all possible models found in your environment, so Meshery automatically generates *and registers* new models and components based on the specific infrastructure Meshery is connected to. Dyanmic models often lack additional metadata, such as descriptions, tags, and relationships, which are typically included in static models.

## Interacting with the Meshery Registry

Use either Meshery UI or CLI to interact with the Registry. Meshery UI offers a user-friendly visual interface for browsing, searching, and managing registry entries. You can easily explore available models, components, and relationships, gaining insights into their properties and connections. Meshery CLI offers commands so that you can register, list, retrieve, update, and delete models, components, and relationships directly from the command line.

### Model Generation

The process of generating a Model (and its entities) is a multi-step process and does not require use of Meshery Server. The process begins with the sourcing of the model information from an authoratitive source: a Registrant. Registrants are responsible for providing all the necessary information to Meshery to generate the model.

#### Using Meshery CLI to Generate Models

Meshery CLI supports the generation of models from a Google Spreadsheet. The Google Spreadsheet should contain a list of model names and source locations from any supported Registrant (e.g. GitHub, Artifact Hub) repositories. The source locations can be a URL to a folder containing Kubernetes CRDs, or to a Helm Chart tar.gz, or an individual Kubernetes Manifest with custom resource definition.

See [`mesheryctl registry generate`](/reference/mesheryctl/registry/generate) for more information.

### Model Registration

Once registered in the Meshery Registry, Models and their entities are available for use within that specific Meshery Server.

Meshery [Adapters]({{ site.baseurl }}/concepts/architecture/adapters) are one example of a Registrant. Registrants are responsible for the registration of entities in the Meshery Registry. Adapters are responsible for the sourcing and registration of entities and the packaging of these enties into one or more models.

#### Using Meshery CLI to Register a Model

```bash
mesheryctl model import -f <path-to-model>
```

#### Using Meshery UI to Register a Model

Visit the Settings --> Registry page and click the "Import" button to import a model.

### Ignoring an Entity

You have control over whether a registered entity (model and all that the model contains) this can be an individual or team-level preference. Use the "Ignore" action to designate whether a given model is allowed to be used within a given Meshery Server deployment. Models that are ignored remain in the Meshery Registry but are not available for use within a given Meshery Server deployment.

