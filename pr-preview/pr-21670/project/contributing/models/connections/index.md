# Contributing to Connections

> How to define a new Connection and register it so Meshery understands it and includes it in the Connection Wizard.

Source: /pr-preview/pr-21670/project/contributing/models/connections/

**Connections are schema-driven.** A Connection's structure, identity, lifecycle, and the forms Meshery renders for it are declared in a **connection definition** that conforms to the [Connection schema](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1beta3/connection) (`connections.meshery.io/v1beta3`) in [`meshery/schemas`](https://github.com/meshery/schemas). Before contributing, familiarize yourself with that schema and read [Contributing to Schemas](/pr-preview/pr-21670/project/contributing/contributing-schemas/) for the development workflow.

This guide explains how to author a connection definition so that Meshery understands a new kind of [Connection](/pr-preview/pr-21670/concepts/logical/connections/) and offers it in the [Connection Wizard](/pr-preview/pr-21670/guides/infrastructure-management/registering-a-connection/). In most cases, **authoring a JSON definition is all you need** - no UI or server code.

## What is a connection definition?

A connection definition is a first-class [Registry](/pr-preview/pr-21670/concepts/logical/registry/) entity, authored per [Model](/pr-preview/pr-21670/concepts/logical/models/) - exactly like [Components](/pr-preview/pr-21670/project/contributing/models/components/) and [Relationships](/pr-preview/pr-21670/project/contributing/models/relationships/). It declares everything Meshery needs to register, render, and manage a kind of Connection:

- its **identity** (`kind`, `type`, `subType`),
- its **lifecycle** (initial `status` and a `transitionMap` of allowed [state transitions](/pr-preview/pr-21670/concepts/logical/connections/#states-and-the-lifecycle-of-connections)),
- the **forms** the Connection Wizard renders (`connectionSchema` and `credentialSchema`), and
- its **visual identity** (`styles`).

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Prerequisite reading</div>


Connection definitions are packaged in the context of a Model. Be sure you understand <a href='/pr-preview/pr-21670/project/contributing/models/models/'>how Models are created and packaged</a> first - without a Model to belong to, your connection definition is homeless.
</div>


## Anatomy of a connection definition

Below is a complete, minimal definition for a hypothetical telemetry backend. It uses the **generic** wizard flow, which is what most contributions should use.



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">{
  &#34;schemaVersion&#34;: &#34;connections.meshery.io/v1beta3&#34;,
  &#34;name&#34;: &#34;Grafana&#34;,
  &#34;description&#34;: &#34;A Grafana instance that brings its dashboards and panels into Meshery.&#34;,
  &#34;kind&#34;: &#34;grafana&#34;,
  &#34;type&#34;: &#34;telemetry&#34;,
  &#34;subType&#34;: &#34;metrics&#34;,
  &#34;status&#34;: &#34;registered&#34;,
  &#34;transitionMap&#34;: {
    &#34;registered&#34;: [
      { &#34;nextState&#34;: &#34;connected&#34;, &#34;description&#34;: &#34;Connect to the Grafana instance.&#34; },
      { &#34;nextState&#34;: &#34;ignored&#34;, &#34;description&#34;: &#34;Keep the registration but do not connect.&#34; }
    ],
    &#34;connected&#34;: [
      { &#34;nextState&#34;: &#34;disconnected&#34;, &#34;description&#34;: &#34;Disconnect the Grafana connection.&#34; },
      { &#34;nextState&#34;: &#34;deleted&#34;, &#34;description&#34;: &#34;Remove the Grafana connection completely.&#34; }
    ]
  },
  &#34;connectionSchema&#34;: {
    &#34;type&#34;: &#34;object&#34;,
    &#34;title&#34;: &#34;Grafana Connection&#34;,
    &#34;required&#34;: [&#34;url&#34;],
    &#34;properties&#34;: {
      &#34;url&#34;: {
        &#34;type&#34;: &#34;string&#34;, &#34;format&#34;: &#34;uri&#34;, &#34;title&#34;: &#34;Grafana Endpoint&#34;,
        &#34;description&#34;: &#34;Base URL of the Grafana instance (e.g. http://grafana.example:3000).&#34;
      },
      &#34;name&#34;: {
        &#34;type&#34;: &#34;string&#34;, &#34;title&#34;: &#34;Connection Name&#34;,
        &#34;description&#34;: &#34;Optional friendly name for this Grafana connection.&#34;
      }
    }
  },
  &#34;credentialSchema&#34;: {
    &#34;type&#34;: &#34;object&#34;,
    &#34;title&#34;: &#34;Grafana Credential&#34;,
    &#34;properties&#34;: {
      &#34;secret&#34;: {
        &#34;type&#34;: &#34;string&#34;, &#34;title&#34;: &#34;API Key or Basic Auth&#34;,
        &#34;description&#34;: &#34;API key, or basic-auth credential formatted as username:password.&#34;
      }
    }
  },
  &#34;styles&#34;: {
    &#34;svgColor&#34;: &#34;&lt;svg ...&gt;...&lt;/svg&gt;&#34;,
    &#34;svgWhite&#34;: &#34;&lt;svg ...&gt;...&lt;/svg&gt;&#34;
  }
}</code>
	</div>
</pre>


### Identity: `kind`, `type`, `subType`

These three fields identify the Connection and determine how the wizard treats it:

- **`kind`** - the genre of Connection (e.g. `grafana`, `prometheus`, `kubernetes`). The wizard groups credentials and renders icons by `kind`.
- **`type`** - a broad classification: `platform`, `telemetry`, `source`, `collaboration`, and so on.
- **`subType`** - a finer classification: `orchestration`, `metrics`, `git`, `registry`, `chat`, and so on.

Together they let the UI target a specific Connection with a [custom wizard extension](#advanced-customizing-the-wizard) when the generic flow is not enough. For reference, the definitions Meshery ships with:

| Connection   | `kind`        | `type`      | `subType`       | initial `status` |
| ------------ | ------------- | ----------- | --------------- | ---------------- |
| Kubernetes   | `kubernetes`  | `platform`  | `orchestration` | `discovered`     |
| Grafana      | `grafana`     | `telemetry` | `metrics`       | `registered`     |
| Prometheus   | `prometheus`  | `telemetry` | `metrics`       | `registered`     |
| Artifact Hub | `artifacthub` | `source`    | `registry`      | `registered`     |
| GitHub       | `github`      | `source`    | `git`           | `registered`     |

### Lifecycle: `status` and `transitionMap`

`status` is the state a freshly created Connection starts in. Manually registered Connections typically start at `registered`; resources that Meshery discovers (like Kubernetes) start at `discovered`.

`transitionMap` declares the **state machine** for the Connection: for each state, the list of states it may move to, each with a human-readable `description` shown as a confirmation prompt in the UI. Use only the canonical states - `discovered`, `registered`, `connected`, `disconnected`, `ignored`, `maintenance`, `deleted`, and `not found` - and keep the transitions consistent with their documented meanings. See [States and the Lifecycle of Connections](/pr-preview/pr-21670/concepts/logical/connections/#states-and-the-lifecycle-of-connections) for what each state means and which transitions make sense.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">The transition map drives the UI</div>


The set of transitions Meshery offers a user for a given Connection comes directly from its definition's `transitionMap`. If a transition is not declared, it is not offered. Model the lifecycle deliberately.
</div>


### Forms: `connectionSchema` and `credentialSchema`

Both are [JSON Schemas](https://json-schema.org/). The Connection Wizard renders them directly with [`react-jsonschema-form`](https://rjsf-team.github.io/react-jsonschema-form/docs/) - the same library used for [Component](/pr-preview/pr-21670/project/contributing/models/components/) forms:

- **`connectionSchema`** becomes the **Configure Connection** step. Mark the fields a Connection cannot exist without (such as `url`) as `required`. A `name` property, if present, is used as the Connection's display name.
- **`credentialSchema`** becomes the **Associate Credential** step. **Omit it** for a Connection that needs no secret - the wizard then skips the credential step entirely.

Because these schemas live on the definition, the wizard needs no per-kind UI code to render them. Adding a property to the schema adds a field to the form.

#### How a credential's `secret` is persisted - and read

The canonical form, declared in [`meshery/schemas`](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1beta1/credential/forms), is a top-level `name` plus a `secret` object holding the kind-specific fields. **The persisted `secret` object *is* the payload.**

Four shapes exist in stored data in total - the canonical one plus three others - and Meshery reads all of them. Two of the three non-canonical shapes are still written today: Meshery writes the Kubernetes shape when it imports a kubeconfig, and its credential form writes the double-nested wrapper.

| Shape | Stored `secret` | Where the payload is |
| --- | --- | --- |
| Canonical | `{"grafanaURL": "...", "grafanaAPIKey": "..."}` | the object itself |
| Kubernetes | `{"auth": {...}, "cluster": {...}}` | the object itself |
| Legacy double-nested | `{"credentialName": "x", "secret": {...}}` | one level down |
| Legacy string | `{"secret": "<token>"}` | a bare string |

Legacy rows are never rewritten - tolerance is what keeps them working. Two mirrored helpers own the whole decision, and every read site goes through them rather than reaching into the map:

- **Go** - `models.CredentialPayload` (the object carrying the credential's fields) and `models.CredentialAuthSecret` (the string auth material) in `server/models/credential_secret.go`.
- **TypeScript** - `resolveCredentialPayload` and `resolveCredentialAuthSecret` in `ui/utils/credentialSecret.ts`.

Ambiguity resolves toward the canonical shape: an object is only unwrapped when it consists of nothing but the legacy wrapper keys (`credentialName`, `name`, `secret`), so a canonical payload that happens to carry its own `secret` field is left alone. The one shape this cannot distinguish is a payload whose *only* fields are `name` and `secret` - no canonical form has that shape, so keep a new credential kind's payload away from it.

What must stay in step between the two files is the **resolution rules**, not the return types: the helpers deliberately differ on the legacy string shape, where the TypeScript `resolveCredentialPayload` returns the bare string and the Go `CredentialPayload` returns nil because it is typed to a map (Go callers use `CredentialAuthSecret`). The auth-key list is currently grafana only; a credential kind whose canonical form holds string auth material under a new property must be added to `canonicalAuthSecretKeys` **and** `CANONICAL_AUTH_SECRET_KEYS`.

The helpers resolve the wrapper, not the field names inside it. A Kubernetes credential created through the credential form describes its cluster with `clusterName`/`clusterServerURL` rather than the kubeconfig-style `cluster` block `K8sContextFromConnection` reads, so it yields an auth block but no cluster; which side moves is tracked in [meshery/meshery#21336](https://github.com/meshery/meshery/issues/21336).

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">The registration payload's `secret` is a string</div>


The server rehydrates `credentialSecret.secret` into `PromCred`/`GrafanaCred`, whose `secret` field is a plain string. Handing it an object fails the `register` (verify) step outright, which is why the wizard sends `resolveCredentialAuthSecret(...)` rather than the payload object.
</div>


### Visual identity: `styles`

`styles` carries inline SVG markup for the kind's icon: `svgColor` (for light backgrounds), `svgWhite` (for dark backgrounds), and optionally `svgComplete`. Follow the same icon conventions as [Components](/pr-preview/pr-21670/project/contributing/models/components/).

### Optional `metadata`

Two optional `metadata` keys tune wizard behavior:

- **`metadata.flow`** - force a wizard flow: `generic` (the default for every kind except Kubernetes) or `kubernetes`. You rarely need to set this.
- **`metadata.docsURL`** - a documentation link surfaced for the kind in the wizard. Defaults to the [Connections](/pr-preview/pr-21670/concepts/logical/connections/) concept page.

## Where the definition lives

Place the definition as a JSON file in a `connections/` folder inside its Model, alongside that Model's `components/` and `relationships/`:

```
models/<model>/<version>/connections/<Name>Connection.json
```

For example, the shipped definitions live under [`models/meshery-core/.../connections/`](https://github.com/meshery/meshery/tree/master/models) as `KubernetesConnection.json`, `GrafanaConnection.json`, `PrometheusConnection.json`, `ArtifactHubConnection.json`, and `GitHubConnection.json`. A Model may include any number of connection definitions. Use these existing files as templates.

## How the definition is registered and consumed

1. **Registration.** On Model registration, Meshery registers the connection definition into the [Registry](/pr-preview/pr-21670/concepts/logical/registry/) under its Model and registrant - the same path as Components and Relationships. A Model (carrying a registrant) is required. Definitions can also be managed over the registry API:

   | Method   | Endpoint                                          | Purpose                                |
   | -------- | ------------------------------------------------- | -------------------------------------- |
   | `GET`    | `/api/registry/connections`                     | List connection definitions            |
   | `GET`    | `/api/registry/connections/{id}`                | Fetch one definition                   |
   | `POST`   | `/api/registry/connections`                     | Register a definition (needs a Model)  |
   | `PUT`    | `/api/registry/connections/{id}`                | Update a definition                    |
   | `DELETE` | `/api/registry/connections/{id}`                | Remove a definition                    |

2. **Consumption.** The [Connection Wizard](/pr-preview/pr-21670/guides/infrastructure-management/registering-a-connection/) lists every registered definition as a creatable kind and renders its `connectionSchema` and `credentialSchema` as wizard steps. **Register your definition and it appears in the wizard automatically** - no UI changes required.

3. **Lifecycle.** Registration drives the default connection state machine, which implements these transitions: `discovered → registered | ignored`, `registered → connected | ignored`, `connected → disconnected | deleted`, and `disconnected → connected | deleted`. Connecting persists the Connection and its credential. **No server code is required** - add a kind-specific action only when the kind needs a reachability probe before it may advance to `connected` (Grafana and Prometheus verify their endpoint; Kubernetes has a bespoke machine altogether). Keep the [`transitionMap`](#lifecycle-status-and-transitionmap) you author in step with this machine: it is the copy the UI shows for each transition.

4. **Seeding.** Most Connections are created by a user through the wizard. A few, though, describe something Meshery itself uses, and Meshery seeds those for itself at boot so they are present out of the box - see [System-owned Connections](#system-owned-connections) below.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Verify it appears</div>


After registering, open the Connection Wizard (**Connections → Create Connection**) and confirm your kind is listed with its icon, that the Configure and Associate Credential steps render your schemas, and that creating a Connection drives the states you declared in the `transitionMap`.
</div>


## System-owned Connections

Meshery generates Models and Components from sources it reaches on its own: it pulls packages from [Artifact Hub](/pr-preview/pr-21670/guides/configuration-management/importing-models/) and reads manifests from GitHub repositories. Those sources are Connections too, and Meshery **seeds them for itself during boot-time seed-data initialization** so a user finds them already present rather than having to create them by hand.

Seeding runs at the end of model registration (`SeedComponents` → `SeedConnections` in `server/models/`), which is also when every path that rebuilds the registry - server start, database reset, hard reset - picks it up. Because the two reset paths drop every table before rebuilding, they re-migrate the **full** system-table set from one shared list (`models.SystemDatabaseModels`, applied via `AutoMigrateSystemTables`) rather than a hand-maintained subset, so the tables a seeded Connection and the Connections page read from - including `environment_connection_mappings` - always exist after a reset. It is driven entirely by the registry: the set of Connections seeded is derived from the registered connection definitions, never from a list of kinds written into the server. A definition is seeded when **both** of the following hold:

1. **Meshery already sources content through the kind** - it holds a *registrant* Connection of that kind that owns registered **Models**, the host those Models, and the Components and Relationships beneath them, are registered under. This is what separates Artifact Hub and GitHub from Kubernetes, Grafana and Prometheus: the latter describe resources a *user* brings, and would be meaningless as empty, endpoint-less rows. Registered Models are required rather than merely any registered entity, because registering a connection definition through `POST /api/registry/connections` creates a registrant of whatever kind the request body names; requiring Models keeps the rule out of reach of request input.
2. **The kind works anonymously** - its `credentialSchema` marks nothing as `required`. A seeded Connection is owned by the system and carries no credential, so a kind that cannot be used without one is never seeded.

A seeded Connection takes its `name`, `type` and `subType` from its definition, which is the authoritative identity for the kind. Three properties are deliberately left alone:

- **`status`** is not re-asserted. The definition's `status` is the state a *new* Connection starts in; once the Connection exists its state belongs to the state machine, so seeding never undoes a user who connected or ignored it.
- **A credential you attach is preserved.** A kind is only seedable because its `credentialSchema` marks nothing as `required`, which is exactly what makes attaching one optional rather than impossible - an Artifact Hub API key, or a GitHub token that raises your API rate limit. Seeding never clears `credentialId`, so a credential you add through the wizard survives every subsequent restart.
- **A user's own Connections are never touched.** Seeding only ever writes the registrant rows Meshery created itself, so a Connection you created of the same kind is left exactly as you left it.

Seeding is idempotent across restarts, and it **never adds a Connection row**. It only ever rewrites rows that model registration already created, and it writes one only when that row does not already match its definition, so a restarted server performs no seeding work at all.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">A kind may already hold more than one registrant</div>


The `registrant` block is hand-authored into every `model.json`, and the shipped blocks are not uniform - some carry a `user_id`, some omit it. A registrant's id is a hash of its content, so those spellings register as *separate* rows and a kind such as `artifacthub` can already hold several. Seeding stamps the definition's identity onto exactly **one** canonical registrant per kind - it reuses whichever registrant already carries that identity, falling back to the lowest id only when none does yet, so the same row is picked on every boot even as a later models release adds more registrant spellings - and leaves every sibling exactly as registration wrote it. Deduplicating the registrant rows themselves would have to repoint `registries.registrant_id` and is tracked separately in [meshery/meshery#20950](https://github.com/meshery/meshery/issues/20950).
</div>


<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Seeding is not a default for your kind</div>


Authoring a definition does **not** get your kind seeded, and it should not: a Connection to *your* Grafana or *your* cluster is yours to create. Seeding exists only for the sources Meshery itself reads from anonymously.
</div>


## Advanced: customizing the wizard

The generic flow - choose → configure → credential → review → done, all derived from your schemas - covers most Connections. When a Connection needs bespoke steps (Kubernetes, for example, imports clusters from a kubeconfig and offers a MeshSync deployment-mode step), register a **connection extension** in the Meshery UI at `ui/components/connections/wizard/registry.ts`.

An extension matches a Connection by `kind` (optionally narrowed by `type`/`subType`) and may override any step; the most specific match wins, and any step left unset falls back to the generic default:

- `detailsStep`, `credentialStep`, `registerStep`, `receiptStep` - override individual steps. Set `credentialStep: null` to remove the credential step (Kubernetes carries its credential inline as a kubeconfig).
- `postConfigSteps` - extra steps appended after registration; these also drive the wizard's **configure** mode for an already-registered Connection.

Each step implements a small contract - an `id` and `label`, a `Component` to render, and optional `canProceed`, `onNext`, `nextLabel`, and `hidden` hooks. Use the existing `kubernetesExtension` as a worked example, and see [Contributing to the Meshery UI](/pr-preview/pr-21670/project/contributing/ui/) for building and testing UI changes.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Prefer schemas over code</div>


Reach for a custom extension only when the generic, schema-driven flow genuinely cannot express what your Connection needs. A definition-only contribution is easier to review, ships without a UI release, and stays consistent with every other Connection in the wizard.
</div>


## Authoring best practices

1. Use `camelCase` for property names, matching the rest of Meshery's schemas.
2. Keep the `transitionMap` consistent with the [documented Connection states](/pr-preview/pr-21670/concepts/logical/connections/#states-and-the-lifecycle-of-connections); do not invent states.
3. Mark only genuinely required fields as `required` in `connectionSchema`, and omit `credentialSchema` when no secret is needed.
4. Start from a [shipped definition](https://github.com/meshery/meshery/tree/master/models) rather than from scratch.
5. Provide both `svgColor` and `svgWhite` icons so the kind renders well on light and dark backgrounds.

## Contribute your Connection

Submit a pull request to the [Meshery repository](https://github.com/meshery/meshery) adding your connection definition to its Model's `connections/` folder, so every user benefits from the new Connection kind. Follow the [contribution gitflow](/pr-preview/pr-21670/project/contributing/contributing-gitflow/) and sign off your commits.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Keeping your Connection private</div>


Prefer to keep a Connection definition private? Bundle it in a custom [Model](/pr-preview/pr-21670/concepts/logical/models/) and [import that Model](/pr-preview/pr-21670/guides/configuration-management/importing-models/) into your Meshery deployment. Your definition is registered in your Meshery Server's [Registry](/pr-preview/pr-21670/concepts/logical/registry/) and offered in your Connection Wizard, without being published upstream.
</div>


<div class="alert alert-dark" role="alert">
  <h4 class="alert-heading">Discussion Forum</h4>
  <p>Don't find an answer to your question here? Ask on the <a href="https://discuss.meshery.io/">Discussion Forum</a>.</p>
</div>
