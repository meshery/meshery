---
title: Contributing to Models Quick Start
description: A no-fluff guide to creating your own Meshery Models quickly.
categories: [contributing]
aliases: [project/contributing/contributing-models-quick-start]
weight: -5
---

**Models follow Meshery's schema-driven development approach.** Model, Component, and Relationship definitions are validated against schemas in [`Model schema`](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1beta1/model). See [Contributing to Schemas]({{< ref "project/contributing/contributing-schemas.md" >}}) for details.

[Meshery Models]({{< ref "concepts/logical/models/index.md" >}}) are a way to represent the architecture of a system or application. Models are defined in JSON and can be used to visualize the components and relationships between them. This guide will walk you through the process of creating a new model.

[Meshery Components]({{< ref "concepts/logical/components.md" >}}) are the building blocks of a model. Each component represents a different part of the system or application. Components can be anything from a database to a microservice to a server. Relationships define how components interact with each other. For example, a database component might have a relationship with a microservice component that represents the microservice's dependency on the database.

### Creating your first Meshery Model

To get started quickly, check out the [Creating Models]({{< ref "guides/configuration-management/creating-models/index.md" >}})  guide for step-by-step instructions.

### Contributing a Model Definition

1. Fork the [meshery/meshery.io](https://github.com/meshery/meshery.io) repository.
1. Create a new branch in your fork of the meshery/meshery.io repository.
1. Add your model definition to the `collections/_models` directory.
1. Create a pull request to the meshery/meshery.io repository.
1. Once your pull request is merged, your model will be available in the next Meshery release.

### Next Steps

To learn more about how to contribute and customize your models, check out the full [Contributing to Models]({{< ref "project/contributing/models" >}}) documentation. We encourage you to get involved in the development of Meshery Models and to share your feedback!

{{% alert color="info" title="Meshery Models are Extensible" %}}
Meshery Models are designed to be extensible, allowing you to define new components as needed. If you have an idea for a new component, please create one and share it with the Meshery community.
{{% /alert %}}
