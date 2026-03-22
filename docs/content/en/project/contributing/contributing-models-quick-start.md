---
title: "Contributing to Models Quick Start"
description: "A no-fluff guide to creating your own Meshery Models quickly."
type: "project"
layout: "project"
---

**Models follow Meshery's schema-driven development approach.** Model, Component, and Relationship definitions are validated against schemas in the [Meshery Schemas repository](https://github.com/meshery/schemas).

[Meshery Models](/concepts/logical/models) are the primary way to represent cloud-native infrastructure. They are defined in JSON and used by Kanvas for visual orchestration.

### Creating your first Meshery Model

To get started quickly, you can generate models directly from existing Kubernetes CRDs or Helm charts using `mesheryctl`.

### Contributing a Model Definition

1. **Fork and Clone:** Fork the [meshery/meshery](https://github.com/meshery/meshery) repository.
2. **Generate Definitions:** Use the Meshery Registry tool or `mesheryctl model import` to generate your model's JSON definitions.
3. **Organize:** Place your new model metadata and component definitions in the `install/kubernetes/bootstrap/models` directory.
4. **Validate:** Ensure your definitions pass the [v1beta1 schema validation](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1beta1/model).
5. **Submit PR:** Create a pull request to the `master` branch of the `meshery/meshery` repository.

### Next Steps

Once merged, your model will be automatically registered and available for use in the **Kanvas Designer** and **Meshery Catalog**.

{{% alert color="info" title="Meshery Models are Extensible" %}}
Meshery Models are designed to be extensible. If you have an idea for a new component or a complex relationship (like ACR-to-AKS), please share it with the community in #meshery-mentorship!
{{% /alert %}}
