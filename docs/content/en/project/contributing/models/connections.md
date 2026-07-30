---
title: Contributing to Connections
description: How to define a new Connection and register it so Meshery understands it and includes it in the Connection Wizard.
categories: [contributing]
aliases:
- /project/contributing/contributing-connection-definitions
- /project/contributing/contributing-connections
weight: 15
---

**Connections are schema-driven.** A Connection's structure, identity, lifecycle, and the forms Meshery renders for it are declared in a **connection definition** that conforms to the [Connection schema](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1beta3/connection) (`connections.meshery.io/v1beta3`) in [`meshery/schemas`](https://github.com/meshery/schemas). Before contributing, familiarize yourself with that schema and read [Contributing to Schemas]({{< ref "project/contributing/contributing-schemas.md" >}}) for the development workflow.

This guide explains how to author a connection definition so that Meshery understands a new kind of [Connection]({{< ref "concepts/logical/connections/index.md" >}}) and offers it in the [Connection Wizard]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}). In most cases, **authoring a JSON definition is all you need** - no UI or server code.

## What is a connection definition?

A connection definition is a first-class [Registry]({{< ref "concepts/logical/registry.md" >}}) entity, authored per [Model]({{< ref "concepts/logical/models/index.md" >}}) - exactly like [Components]({{< ref "project/contributing/models/components" >}}) and [Relationships]({{< ref "project/contributing/models/relationships" >}}). It declares everything Meshery needs to register, render, and manage a kind of Connection:

- its **identity** (`kind`, `type`, `subType`),
- its **lifecycle** (initial `status` and a `transitionMap` of allowed [state transitions]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}})),
- the **forms** the Connection Wizard renders (`connectionSchema` and `credentialSchema`), and
- its **visual identity** (`styles`).

{{% alert color="warning" title="Prerequisite reading" %}}
Connection definitions are packaged in the context of a Model. Be sure you understand <a href='{{< ref "project/contributing/models/models" >}}'>how Models are created and packaged</a> first - without a Model to belong to, your connection definition is homeless.
{{% /alert %}}

## Anatomy of a connection definition

Below is a complete, minimal definition for a hypothetical telemetry backend. It uses the **generic** wizard flow, which is what most contributions should use.

{{< code code=`{
  "schemaVersion": "connections.meshery.io/v1beta3",
  "name": "Grafana",
  "description": "A Grafana instance that brings its dashboards and panels into Meshery.",
  "kind": "grafana",
  "type": "telemetry",
  "subType": "metrics",
  "status": "registered",
  "transitionMap": {
    "registered": [
      { "nextState": "connected", "description": "Connect to the Grafana instance." },
      { "nextState": "ignored", "description": "Keep the registration but do not connect." }
    ],
    "connected": [
      { "nextState": "disconnected", "description": "Disconnect the Grafana connection." },
      { "nextState": "deleted", "description": "Remove the Grafana connection completely." }
    ]
  },
  "connectionSchema": {
    "type": "object",
    "title": "Grafana Connection",
    "required": ["url"],
    "properties": {
      "url": {
        "type": "string", "format": "uri", "title": "Grafana Endpoint",
        "description": "Base URL of the Grafana instance (e.g. http://grafana.example:3000)."
      },
      "name": {
        "type": "string", "title": "Connection Name",
        "description": "Optional friendly name for this Grafana connection."
      }
    }
  },
  "credentialSchema": {
    "type": "object",
    "title": "Grafana Credential",
    "properties": {
      "secret": {
        "type": "string", "title": "API Key or Basic Auth",
        "description": "API key, or basic-auth credential formatted as username:password."
      }
    }
  },
  "styles": {
    "svgColor": "<svg ...>...</svg>",
    "svgWhite": "<svg ...>...</svg>"
  }
}` >}}

### Identity: `kind`, `type`, `subType`

These three fields identify the Connection and determine how the wizard treats it:

- **`kind`** - the genre of Connection (e.g. `grafana`, `prometheus`, `kubernetes`). The wizard groups credentials and renders icons by `kind`.
- **`type`** - a broad classification: `platform`, `telemetry`, `collaboration`, and so on.
- **`subType`** - a finer classification: `orchestration`, `metrics`, `git`, `chat`, and so on.

Together they let the UI target a specific Connection with a [custom wizard extension](#advanced-customizing-the-wizard) when the generic flow is not enough. For reference, the definitions Meshery ships with:

| Connection | `kind`       | `type`      | `subType`       | initial `status` |
| ---------- | ------------ | ----------- | --------------- | ---------------- |
| Kubernetes | `kubernetes` | `platform`  | `orchestration` | `discovered`     |
| Grafana    | `grafana`    | `telemetry` | `metrics`       | `registered`     |
| Prometheus | `prometheus` | `telemetry` | `metrics`       | `registered`     |

### Lifecycle: `status` and `transitionMap`

`status` is the state a freshly created Connection starts in. Manually registered Connections typically start at `registered`; resources that Meshery discovers (like Kubernetes) start at `discovered`.

`transitionMap` declares the **state machine** for the Connection: for each state, the list of states it may move to, each with a human-readable `description` shown as a confirmation prompt in the UI. Use only the canonical states - `discovered`, `registered`, `connected`, `disconnected`, `ignored`, `maintenance`, `deleted`, and `not found` - and keep the transitions consistent with their documented meanings. See [States and the Lifecycle of Connections]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}) for what each state means and which transitions make sense.

{{% alert color="info" title="The transition map drives the UI" %}}
The set of transitions Meshery offers a user for a given Connection comes directly from its definition's `transitionMap`. If a transition is not declared, it is not offered. Model the lifecycle deliberately.
{{% /alert %}}

### Forms: `connectionSchema` and `credentialSchema`

Both are [JSON Schemas](https://json-schema.org/). The Connection Wizard renders them directly with [`react-jsonschema-form`](https://rjsf-team.github.io/react-jsonschema-form/docs/) - the same library used for [Component]({{< ref "project/contributing/models/components" >}}) forms:

- **`connectionSchema`** becomes the **Configure Connection** step. Mark the fields a Connection cannot exist without (such as `url`) as `required`. A `name` property, if present, is used as the Connection's display name.
- **`credentialSchema`** becomes the **Associate Credential** step. **Omit it** for a Connection that needs no secret - the wizard then skips the credential step entirely.

Because these schemas live on the definition, the wizard needs no per-kind UI code to render them. Adding a property to the schema adds a field to the form.

### Visual identity: `styles`

`styles` carries inline SVG markup for the kind's icon: `svgColor` (for light backgrounds), `svgWhite` (for dark backgrounds), and optionally `svgComplete`. Follow the same icon conventions as [Components]({{< ref "project/contributing/models/components" >}}).

### Optional `metadata`

Two optional `metadata` keys tune wizard behavior:

- **`metadata.flow`** - force a wizard flow: `generic` (the default for every kind except Kubernetes) or `kubernetes`. You rarely need to set this.
- **`metadata.docsURL`** - a documentation link surfaced for the kind in the wizard. Defaults to the [Connections]({{< ref "concepts/logical/connections/index.md" >}}) concept page.

## Where the definition lives

Place the definition as a JSON file in a `connections/` folder inside its Model, alongside that Model's `components/` and `relationships/`:

```
models/<model>/<version>/connections/<Name>Connection.json
```

For example, the shipped definitions live under [`models/meshery-core/.../connections/`](https://github.com/meshery/meshery/tree/master/models) as `KubernetesConnection.json`, `GrafanaConnection.json`, and `PrometheusConnection.json`. A Model may include any number of connection definitions. Use these existing files as templates.

## How the definition is registered and consumed

1. **Registration.** On Model registration, Meshery registers the connection definition into the [Registry]({{< ref "concepts/logical/registry.md" >}}) under its Model and registrant - the same path as Components and Relationships. A Model (carrying a registrant) is required. Definitions can also be managed over the registry API:

   | Method   | Endpoint                                          | Purpose                                |
   | -------- | ------------------------------------------------- | -------------------------------------- |
   | `GET`    | `/api/registry/connections`                     | List connection definitions            |
   | `GET`    | `/api/registry/connections/{id}`                | Fetch one definition                   |
   | `POST`   | `/api/registry/connections`                     | Register a definition (needs a Model)  |
   | `PUT`    | `/api/registry/connections/{id}`                | Update a definition                    |
   | `DELETE` | `/api/registry/connections/{id}`                | Remove a definition                    |

2. **Consumption.** The [Connection Wizard]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) lists every registered definition as a creatable kind and renders its `connectionSchema` and `credentialSchema` as wizard steps. **Register your definition and it appears in the wizard automatically** - no UI changes required.

{{% alert color="info" title="Verify it appears" %}}
After registering, open the Connection Wizard (**Connections → Create Connection**) and confirm your kind is listed with its icon, that the Configure and Associate Credential steps render your schemas, and that creating a Connection drives the states you declared in the `transitionMap`.
{{% /alert %}}

## Advanced: customizing the wizard

The generic flow - choose → configure → credential → review → done, all derived from your schemas - covers most Connections. When a Connection needs bespoke steps (Kubernetes, for example, imports clusters from a kubeconfig and offers a MeshSync deployment-mode step), register a **connection extension** in the Meshery UI at `ui/components/connections/wizard/registry.ts`.

An extension matches a Connection by `kind` (optionally narrowed by `type`/`subType`) and may override any step; the most specific match wins, and any step left unset falls back to the generic default:

- `detailsStep`, `credentialStep`, `registerStep`, `receiptStep` - override individual steps. Set `credentialStep: null` to remove the credential step (Kubernetes carries its credential inline as a kubeconfig).
- `postConfigSteps` - extra steps appended after registration; these also drive the wizard's **configure** mode for an already-registered Connection.

Each step implements a small contract - an `id` and `label`, a `Component` to render, and optional `canProceed`, `onNext`, `nextLabel`, and `hidden` hooks. Use the existing `kubernetesExtension` as a worked example, and see [Contributing to the Meshery UI]({{< ref "project/contributing/ui" >}}) for building and testing UI changes.

{{% alert color="warning" title="Prefer schemas over code" %}}
Reach for a custom extension only when the generic, schema-driven flow genuinely cannot express what your Connection needs. A definition-only contribution is easier to review, ships without a UI release, and stays consistent with every other Connection in the wizard.
{{% /alert %}}

## Authoring best practices

1. Use `camelCase` for property names, matching the rest of Meshery's schemas.
2. Keep the `transitionMap` consistent with the [documented Connection states]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}); do not invent states.
3. Mark only genuinely required fields as `required` in `connectionSchema`, and omit `credentialSchema` when no secret is needed.
4. Start from a [shipped definition](https://github.com/meshery/meshery/tree/master/models) rather than from scratch.
5. Provide both `svgColor` and `svgWhite` icons so the kind renders well on light and dark backgrounds.

## Contribute your Connection

Submit a pull request to the [Meshery repository](https://github.com/meshery/meshery) adding your connection definition to its Model's `connections/` folder, so every user benefits from the new Connection kind. Follow the [contribution gitflow]({{< ref "project/contributing/contributing-gitflow.md" >}}) and sign off your commits.

{{% alert color="info" title="Keeping your Connection private" %}}
Prefer to keep a Connection definition private? Bundle it in a custom [Model]({{< ref "concepts/logical/models/index.md" >}}) and [import that Model]({{< ref "guides/configuration-management/importing-models/index.md" >}}) into your Meshery deployment. Your definition is registered in your Meshery Server's [Registry]({{< ref "concepts/logical/registry.md" >}}) and offered in your Connection Wizard, without being published upstream.
{{% /alert %}}

{{< discuss >}}
