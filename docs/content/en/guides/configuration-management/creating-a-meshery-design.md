---
title: Creating a Meshery Design
description: Learn how to create a Meshery Design using the built-in Design Configurator in Meshery UI or the mesheryctl CLI.
categories: [configuration]
aliases:
- /guides/configuration-management/working-with-designs/
- /guides/configuration-management/working-with-designs
- /tasks/patterns
---

A Meshery Design is the primary unit of configuration management in Meshery. It is a declarative document that describes the desired state of your infrastructure and applications — the components you want, their configuration, and their relationships. Designs can be deployed, shared, versioned, exported, and imported.

See [Meshery Designs]({{< ref "concepts/logical/designs.md" >}}) for a full description of design capabilities.

## Ways to Create a Design

{{< tabs id="creating-a-meshery-design-tabs" >}}
Meshery UI | fa fa-desktop

## Using the Design Configurator in Meshery UI

The Design Configurator is a built-in tool in Meshery UI. It lets you browse infrastructure categories and models, add components to a design, configure each component through guided forms, and save the resulting design — all without writing YAML by hand.

### Step 1 — Open the Design Configurator

1. Log in to Meshery and go to the **Designs** page (left navigation).
2. Click **+ New Design** (or open an existing design to edit it).  
   The Design Configurator opens with an empty canvas and a component panel on the left.

Opening the configurator requires the **View Designs** permission, the same key the **Designs** page itself requires. Without it, the configurator shows the permission-denied page rather than the canvas. See [Extensibility: Authorization]({{< ref "reference/extensibility/authorization/index.md" >}}).

### Step 2 — Name Your Design

Give your design a meaningful name in the **Design Name** field at the top of the configurator. This name is used when saving, sharing, or deploying the design.

### Step 3 — Add Components

1. In the **Category** dropdown, select the infrastructure category you want to work with (for example, *Kubernetes*, *AWS*, *Prometheus*).
2. In the **Model** dropdown, select the specific model within that category (for example, *Deployment*, *Service*, *ConfigMap* within the Kubernetes category).
3. Click the model or component name to add it to your design. The component appears in the design document on the right.

Repeat this process to add as many components as your design requires.

### Step 4 — Configure Components

Click any component in the design panel to open its configuration form. The form is generated from the component's schema and includes:

- Required fields (highlighted)
- Optional fields with defaults
- Nested sub-properties (expand to configure)

Fill in the fields for your environment. Changes are applied to the design document in real time.

### Step 5 — Save the Design

Click **Save** (floppy disk icon) to save your design. Meshery stores the design in your account. Use **Save As** to create a copy under a new name.

Your saved design appears on the **Designs** page, where you can deploy, export, share, or further edit it.

## Using the Design Configurator to Edit YAML Directly

The Design Configurator also exposes a **code editor** panel alongside the form view. If you prefer to write or paste YAML directly:

1. Open or create a design.
2. Switch to the **YAML/Code** view in the configurator toolbar.
3. Enter valid Meshery Design YAML (following the [Meshery Schemas](https://github.com/meshery/schemas) spec).
4. Click **Save**.

Changes made in the code editor are reflected immediately in the form view, and vice versa.

## Seed Designs

When you start Meshery for the first time, a set of seed designs is available. These cover common Kubernetes patterns and serve as a starting point for exploration.

You can also import community designs from the [Meshery Catalog](https://meshery.io/catalog) or from a Git repository.

<img src="../images/pattern-import.png" width="60%" alt="Importing a design in Meshery UI" />

<!-- tab -->
mesheryctl | fa fa-terminal

## Using mesheryctl

You can also create and manage designs from the command line using `mesheryctl`.

### Import a design from a file

```bash
mesheryctl design import -f your-design.yaml
```

### Apply a design by file

```bash
mesheryctl design apply -f your-design.yaml
```

### Apply an already-imported design by name

```bash
mesheryctl design apply MyDesignName
```

### List saved designs

```bash
mesheryctl design list
```

See the [`mesheryctl design` reference]({{< ref "reference/references/mesheryctl/design/_index.md" >}}) for the full subcommand reference.
{{< /tabs >}}

## Importing and Configuring Sample Designs {#importing-and-configuring-sample-designs}

Users can also import sample designs to their remote provider from this [sample repository](https://github.com/service-mesh-patterns/service-mesh-patterns/tree/master/samples).

<img src=".../images/pattern-import.png" width="50%" />

Once sample designs are imported, you can then edit sample designs or use the pattern configurator to configure them according to your requirements.

<img src="../images/pattern-configure-button.png" width="50%" />

<img src="../images/pattern-configure.png" width="50%" />

## WASM Filters

Meshery can be used for managing WebAssembly Filters through the UI or the CLI.

### Filter Management Through Meshery UI

Like patterns, Meshery also comes with some sample WebAssembly Filters for you to experiment.

<img src="../images/meshery-filters.png" width="50%" />

You can also import these filters manually to your provider from the [wasm-filters](https://github.com/meshery-extensions/wasm-filters) repo.

Meshery's sample application ImageHub will let you test out configuring these filters out-of-the-box.

You can onboard ImageHub to an installed service mesh as shown below.

<img src="../images/image-hub.png" width="50%" />

### Filter Management Through Meshery CLI

You can also manage WASM filters through Meshery's CLI, mesheryctl.

The `mesheryctl filter` command lets you import and configure WebAssembly filters.

For example,

```
mesheryctl exp filter apply -f metrics_collector_bg.wasm
```

If you already have a filter imported into Meshery, you can configure the filter by name.

```
mesheryctl exp filter apply metrics_collector_bg
```
