---
title: Creating a Meshery Design
description: Learn how to create a Meshery Design using the built-in Design Configurator in Meshery UI or the mesheryctl CLI.
categories: [configuration]
aliases:
- /guides/configuration-management/working-with-designs
- /guides/configuration-management/pattern-management
- /tasks/patterns
- /tasks/pattern-management/
---

A Meshery Design is the primary unit of configuration management in Meshery. It is a declarative document that describes the desired state of your infrastructure and applications: the components you want, their configuration, and their relationships. You can deploy, share, version, export, and import designs.

See [Meshery Designs]({{< ref "concepts/logical/designs.md" >}}) for a full description of design capabilities.

## Ways to Create a Design

{{< tabs id="creating-a-meshery-design-tabs" >}}
Meshery UI | fa fa-desktop

## Using the Design Configurator in Meshery UI

The Design Configurator is built into Meshery UI. It lists the components of a model, gives each one a form generated from its schema, and shows the resulting design as YAML next to the forms, so you don't have to write YAML by hand.

### Step 1: Open the Design Configurator

1. In Meshery UI, go to **Configuration** > **Designs**.
2. Click **Create Design**. To edit an existing design instead, open it from the list.

The configurator opens with a YAML editor. Opening it requires the **View Designs** permission, the same key the **Designs** page requires. Without it, the configurator shows the permission-denied page instead. See [Extensibility: Authorization]({{< ref "reference/extensibility/authorization/index.md" >}}).

### Step 2: Name your design

Enter a name in the **Design Name** field in the toolbar. Meshery uses this name when you save, share, or deploy the design.

### Step 3: Choose a model

1. In the **Select Category** dropdown, choose a category, for example *Orchestration & Management*.
2. In the **Select Model** dropdown, choose a model in that category, for example *Kubernetes*.

The configurator lists a form for each component of the model, such as *Deployment*, *Service*, and *ConfigMap*, next to the YAML editor.

### Step 4: Configure components

Expand a component's form and fill in its fields. Each form is generated from the component's schema, so it marks required fields and shows nested properties you can expand. As you fill in a form, the configurator adds the component to the design and updates the YAML.

Choose another category or model at any time to add components from it to the same design.

### Step 5: Save the design

- For a new design, click the **Save Design as New File** icon.
- For a design you opened from the list, click the **Update Design** icon to save your changes, or **Save Design as New File** to save a copy.

Your saved design appears on the **Designs** page, where you can deploy, export, share, or edit it.

## Editing YAML directly

The YAML editor sits next to the component forms. To write or paste a design by hand, edit the YAML in that editor. It must follow the design schema in [Meshery Schemas](https://github.com/meshery/schemas). Your edits update the design, and changes from the forms appear in the editor.

## Starting from an existing design

Instead of starting from nothing, import a design from the [Meshery Catalog](https://meshery.io/catalog), from a URL, or from a Kubernetes manifest, Helm chart, or Docker Compose file. See [Importing and Exporting Designs]({{< ref "guides/configuration-management/import-export-designs.md" >}}).

<img src="../images/pattern-import.png" width="60%" alt="Importing a design in Meshery UI" />

<!-- tab -->
mesheryctl | fa fa-terminal

## Using mesheryctl

Use `mesheryctl` to bring a design file into Meshery and deploy it from the command line.

### Import a design from a file

```bash
mesheryctl design import -f your-design.yaml -n "My Design"
```

`-f` accepts a local path or a direct URL. The file can be a Meshery design, a Kubernetes manifest, a Helm chart, or a Docker Compose file.

### Deploy a design from a file

```bash
mesheryctl design apply -f your-design.yaml
```

### Deploy a saved design by name

```bash
mesheryctl design apply MyDesignName
```

### List saved designs

```bash
mesheryctl design list
```

See the [`mesheryctl design` reference]({{< ref "reference/references/mesheryctl/design/_index.md" >}}) for the full subcommand reference.
{{< /tabs >}}
