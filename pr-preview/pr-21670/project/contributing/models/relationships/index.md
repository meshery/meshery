# Contributing to Relationships

> How to contribute to Meshery Models Relationships, Policies...

Source: /pr-preview/pr-21670/project/contributing/models/relationships/

**Relationships follow a schema-defined structure.** The [Relationship schema](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1beta3/relationship) (`relationships.meshery.io/v1beta3`) is the single source of truth for how relationships between components are expressed. The in-tree corpus under `models/**/relationships/` is still mostly `v1beta2`. The version shapes are compatible - Meshery Server bridges registered definitions to the `v1beta2` shape for its policy engine - but registration on current releases accepts only `v1beta2`/`v1alpha3` documents ([meshkit#1096](https://github.com/meshery/meshkit/pull/1096) adds `v1beta3`), so a definition that must register today declares `v1beta2`. Refer to the schema when defining new relationship types or selectors. See [Contributing to Schemas](/pr-preview/pr-21670/project/contributing/contributing-schemas/) for details.

Coding agents: use the `gen-relationship` skill in `.agents/skills/gen-relationship/` (one example of every canonical `kind` / `type` / `subType`).

[Relationships](/pr-preview/pr-21670/concepts/logical/relationships/) within [Models](/pr-preview/pr-21670/concepts/logical/models/) classify how [Components](/pr-preview/pr-21670/concepts/logical/components/) relate and whether they affect each other.

## Overview of Steps to Create Relationships

**Prework:**

1. [Relationship Identification](#relationship-identification)
2. [Relationship Classification](#relationship-visualizations)

**Development:**

3. [Relationship Definition](#relationship-definitions)
4. [Relationship Scopes](#relationship-scopes)

**Postwork:**

5. [Relationship Authoring Best Practices and Considerations](#relationship-authoring-best-practices-and-considerations)
6. [Relationship Contribution](#relationship-contribution)

## Prework

<a id="relationship-identification"></a>

### 1. Characterize the relationship and any specific constraints

Using your domain expertise, define the qualities of this new relationship. Identify and qualify any specific constraints to be enforced between one or more specific components within the same or different models.

For example, a Kubernetes `Service` can have a network relationship with a Kubernetes `Deployment`. That is `kind: edge`, `type: non-binding`, `subType: network`.

<details close>
<summary>Relationship Example</summary>
<pre><code class="language-json highlighter-rouge">
{
  "id": "00000000-0000-0000-0000-000000000000",
  "schemaVersion": "relationships.meshery.io/v1beta3",
  "version": "v1.0.0",
  "kind": "edge",
  "type": "non-binding",
  "subType": "network",
  "status": "enabled",
  "evaluationQuery": "",
  "metadata": {
    "description": "A Service selects Pods of a Deployment."
  },
  "model": {
    "id": "00000000-0000-0000-0000-000000000000",
    "name": "kubernetes",
    "version": "v1.0.0",
    "displayName": "kubernetes",
    "registrant": { "kind": "github" },
    "model": { "version": "" }
  },
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "kind": "Service",
            "model": { "name": "kubernetes" }
          }
        ],
        "to": [
          {
            "kind": "Deployment",
            "model": { "name": "kubernetes" }
          }
        ]
      },
      "deny": {
        "from": [],
        "to": []
      }
    }
  ]
}
</code></pre>
</details>

You might *also* know that this relationship should not form when the destination is not a workload the Service can select. Encode that with `deny` in the same selector-set item.

Codify relationships using your domain expertise. The `kind`, `type`, and `subType` together pick both the visual paradigm and the evaluation policy.

#### Kind, type, and subType

`kind` is a schema enum: `hierarchical`, `edge`, or `sibling`. `type` and `subType` are open strings. Use an established combination rather than inventing one.

| kind | type | subType | Meaning |
|---|---|---|---|
| `hierarchical` | `parent` | `inventory` | Parent contains/scopes children. Parent identity is patched onto the child: the Namespace's name lands in each namespaced resource's `metadata.namespace`. |
| `hierarchical` | `parent` | `alias` | Child is a nested object inside the parent (Container inside Pod). |
| `hierarchical` | `parent` | `wallet` | Child configuration is held/patched into the parent (WASMFilter → EnvoyFilter). |
| `hierarchical` | `sibling` | `matchlabels` | **In-tree tagsets encoding.** Components that share labels. Schema also allows `kind: sibling`; do not mix encodings in one model. |
| `edge` | `non-binding` | `reference` | Logical name/id pointer (Deployment → ConfigMap). |
| `edge` | `non-binding` | `network` | Documented network selection without provisioning an attachment (Service → Deployment). |
| `edge` | `binding` | `network` | Connecting provisions or rewrites network identity. Rare. |
| `edge` | `non-binding` | `firewall` | Policy that allows or denies traffic (NetworkPolicy → Pod). |
| `edge` | `binding` | `permission` | Assigns identities (Role → ServiceAccount). |
| `edge` | `non-binding` | `permission` | Mentions a role or identity without binding it. |
| `edge` | `binding` | `mount` | Storage or device is attached (PVC → Pod). |
| `edge` | `non-binding` | `annotation` | Designer-only line. Set `metadata.isAnnotation: true`. No patch. |
| `edge` | `non-binding` | `alias` | Named stand-in, not nested ownership. |
| `edge` | `non-binding` | `inventory` | Rare peer index/list. Prefer hierarchical parent inventory for containment. |

`badge` is a visual paradigm only; there is no in-tree encoding. Propose a visualization before introducing `subType: badge`.

**Hierarchical `from` / `to`:** `from` is the child, `to` is the parent.

- Inventory: parent **mutates** the child.
- Alias and wallet: child **mutates** the parent.

**Binding vs non-binding:** `binding` means forming the relationship assigns, mounts, or entitles. `non-binding` means the link is real (selector, name, policy match) but does not itself provision the attachment.

<a id="relationship-visualizations"></a>

### 2. Classify relationship type and specify visual representation

Browse and pick the most appropriate visualization for this relationship by using one of the predefined relationship visualizations.

<style>
    figcaption {
        margin: 1rem auto;
    }
</style>
<details open>
    <summary>Example Visual Representations</summary>
    <details><summary>Kind: Hierarchical</summary>
        <figure>
            <figcaption>subType: Parent | Namespace (Parent) and ConfigMap (child), Role (Child)<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=6370ffcd-13a6-4a65-b426-30f1e63dc381"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Hierarchical - Parent: Namespace to other components" src="/pr-preview/pr-21670/concepts/logical/relationships/images/hierarchical_parent_relationship.svg"/>
        <br>
        <figure>
            <figcaption>subType: Inventory | Namespace and ConfigMap<a target="_blank" href="https://playground.meshery.io/extension/meshmap?design=21d40e36-8ab7-4f9f-9fed-f6a818510446"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Hierarchical - Parent: Namespace to other components" src="/pr-preview/pr-21670/concepts/logical/relationships/images/hierarchical_inventory_relationship.svg"/>
    </details>
    <details><summary>Kind: Edge</summary>
        <figure>
            <figcaption>Type: `Non-Binding`, subType: `Permission`: Cluster Role with Cluster Role Binding to Service Account<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=7dd39d30-7b14-4f9f-a66c-06ba3e5000fa"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Binding" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_permission_relationship_cluster_role_service_account.svg"/>
        <figure>
            <figcaption>Type: `Binding`, subType: `Mount`: Pod to Persistent volume via Persistent volume claim<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=43d5fdfe-25f8-4c2c-be9d-30861bbc2a08"> (open in playground)</a> </figcaption>
        </figure>
        <img alt="Edge - Mount" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_mount_relationship_pod_persistent_volume.svg"/>
        <br>
        <figure>
            <figcaption>Type: `Non-Binding`, subType: `Network`: Ingress to Service<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=1f79b0c6-2efe-4ee9-b08c-e1bd07a3926b"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Edge - Network" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_network_ingress_to_service_relationship.svg"/>
        <br>
        <figure>
            <figcaption>Type: `Non-Binding`, subType: `Network`: Service to Pod<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=90a9b4a0-a296-44b5-b1c5-7b1cb4827a77"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Edge - Network: Ingress to Service" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_network_service_to_pod_relationship.svg"/>
        <br>
        <figure>
            <figcaption>Type: `Non-Binding`, subType: `Network`: Service to Service<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=4e368e07-5039-400e-b637-96b0241af799"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Edge - Network" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_network_service_to_service_relationship.svg"/>
        <br>
        <figure>
            <figcaption>Type: `Non-Binding`, subType: `Network`: Service to Endpoint<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=ab35416d-7cf7-4540-8b2e-7271ffeadde2"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Edge - Network" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_network_service_to_endpoints_relationship.svg"/>
        <br>
        <figure>
            <figcaption>Type: `Non-Binding`, subType: `Network`: Service to Deployment<a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=33742281-428d-4340-b42e-6a0fd4ba1d0a"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Edge - Network" src="/pr-preview/pr-21670/concepts/logical/relationships/images/network_edge_relationship_service_deployment.svg"/>
        <br>
        <figure>
            <figcaption><code>type:non-binding</code>, subType: `Network`: Network Policy (Pod to Pod) <a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=58fda714-eaa4-490f-b228-b8bcfe3a1e47s"> (open in playground)</a></figcaption>
        </figure>
        <img alt="Edge - Network Policy" src="/pr-preview/pr-21670/concepts/logical/relationships/images/edge_firewall_relationship_pod_to_pod.svg" />
    </details>
</details>


Once selected, note the relationship's `kind`, `type`, and `subType`. If an existing visualization does not seem appropriate, propose a new one. Use the whiteboard feature of Meshery's extensions to sketch the relationship.

## Development

<a id="relationship-definitions"></a>

### 3. Create a Relationship Definition as a JSON file

Create a relationship definition as a JSON file, placing this new definition file into its respective model folder (see [Contributing to Models](/pr-preview/pr-21670/project/contributing/models/)). A model may include any number of relationship definitions. Filename convention: `{kind}-{type}-{subType}-<suffix>.json`.

Include:

- `schemaVersion`: `relationships.meshery.io/v1beta3` is the authoring target, but declare `v1beta2` for definitions that must register on current servers (see the note at the top of this page). Keep `v1beta2` when refining an existing in-tree file unless you are deliberately migrating it.
- `kind`: The genre of relationship (`hierarchical`, `edge`, `sibling`).
- `type`: The augmentative category (`parent`, `binding`, `non-binding`, `sibling`, …).
- `subType`: The specific visual paradigm (`inventory`, `mount`, `network`, `wallet`, `reference`, `matchlabels`, …).
- `selectors`: The scope of the relationship. One selector-set item is an OR. Inside an item, every `from` entry relates to every `to` entry - a cross-product (AND).
- `evaluationQuery`: Deprecated. Set it to `""` as every in-tree definition does; the evaluation engine enters through the fixed `data.relationship_evaluation_policy` package and dispatches on `kind`/`type`/`subType`.
- `metadata.description`: A characterization of the relationship, its purpose, and constraints.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Use Existing Relationships as Examples</div>


Browse the <a href='https://github.com/meshery/meshery/tree/master/models'>existing relationships in the Meshery repository</a> and the pedagogical examples in <code>.agents/skills/gen-relationship/examples/</code>. For a prior pull request as a template, see <a href='https://github.com/meshery/meshery/pull/9880/files'>PR #9880</a>.
</div>


<a id="relationship-scopes"></a>

### 4. Configuring the Scope of Relationships

The extent to which a relationship affects components within a model or beyond a model is defined and controlled using scopes.

#### Global Scope

Global scope is defined using the `model` attribute in the relationship definition.

Relationships can be confined to a specific model or allowed to affect all models. For example, if the model is specified as `aws-ec2-controller`, the relationship will work for those components that belong to the `aws-ec2-controller` model.

#### Local Scope

Local scope is defined and controlled via `selectors` in the relationship definition.

Relationship selectors refine applicability. Selector details determine whether there is a match: which models and components are involved, and any constraints. Selector items combine with AND. Selector **sets** (the `selectors` array) combine with OR.

Selectors are an array. Each entry has `allow` and optional `deny`, each with `from` and `to`. Only components inside the same selector-set item relate to each other: each object in `from` relates to each object in `to`.

When many `from`/`to` combinations would otherwise force a complicated `deny`, split them into additional selector-set items.

*Note: When defining Hierarchical relationships, the `from` field represents the child component, while the `to` field represents the parent component.*

#### Actions: mutatorRef and mutatedRef

Patches copy values from one component to another when the selector matches. Both fields are nested arrays of string path segments (`string[][]`). Sequence length must match: index `i` of `mutatorRef` is copied onto index `i` of `mutatedRef`.

```json
"mutatorRef": [["config", "url"], ["config", "name"]],
"mutatedRef": [["configPatch", "value"], ["name"]]
```

`[config, url]` is patched onto `[configPatch, value]`; `[config, name]` onto `[name]`.

| Field | Role |
|---|---|
| `mutatorRef` | **Source.** JSON path of the value to read. |
| `mutatedRef` | **Sink.** Path segments of the field to patch. |
| `patchStrategy` | How to apply. Schema enum: `merge`, `strategic`, `add`, `remove`, `replace`, `copy`, `move`, `test`. The in-tree corpus and the evaluation engine use `replace` exclusively; default to `replace` unless you need different semantics. |

Paths are relative to the Meshery component document (`configuration`, `displayName`, `component.kind`, …), not the raw Kubernetes YAML root. `_` may mark only the first array position in a path; later arrays need an explicit index. Omit `patch` when the relationship only matches (tagsets, annotation).

<details close>
<summary>Relationship Selector Example</summary>

<pre><code class="language-json highlighter-rouge">
"selectors": [
  {
    "allow": {
      "from": [
        {
          "kind": "WASMFilter",
          "model": { "name": "meshery-core" },
          "patch": {
            "patchStrategy": "replace",
            "mutatorRef": [
              ["configuration", "config"]
            ]
          }
        }
      ],
      "to": [
        {
          "kind": "EnvoyFilter",
          "model": { "name": "istio-base" },
          "patch": {
            "patchStrategy": "replace",
            "mutatedRef": [
              ["configuration", "spec", "configPatches", "_", "patch", "value"]
            ]
          }
        }
      ]
    },
    "deny": {
      "from": [],
      "to": []
    }
  },
  {
    "allow": {
      "from": [
        {
          "kind": "ConfigMap",
          "model": { "name": "kubernetes" },
          "patch": {
            "patchStrategy": "replace",
            "mutatorRef": [
              ["configuration", "metadata", "name"]
            ]
          }
        }
      ],
      "to": [
        {
          "kind": "Deployment",
          "model": { "name": "kubernetes" },
          "patch": {
            "patchStrategy": "replace",
            "mutatedRef": [
              ["configuration", "spec", "template", "spec", "containers", "_", "envFrom", "0", "configMapRef", "name"]
            ]
          }
        }
      ]
    },
    "deny": {
      "from": [],
      "to": []
    }
  }
]
</code></pre>
<br/>

The first selector-set item (WASMFilter → EnvoyFilter) is independent of the second (ConfigMap → Deployment). Use separate items when the pairs should not cross-match.

</details>

The WASMFilter example is hierarchical parent **wallet** (child config patched into the parent). The ConfigMap example is edge **reference** (a name pointer), not hierarchical inventory.

#### Understanding Relationship Policies and their Evaluation

Meshery evaluates designs with Open Policy Agent. The engine enters through the fixed policy package `data.relationship_evaluation_policy` and dispatches on each definition's `kind`, `type`, and `subType`. The policies live under `models/meshery-core/<version>/<definition-version>/policies/`.

**What value should `evaluationQuery` carry?**

An empty string. The property is deprecated (see the schema's deprecation notice) and the current engine ignores it; every in-tree definition sets `""`. The historical per-relationship rule name was `{kind}_{subType}_relationship` — no `{type}` segment — and any value you do set must be a valid Rego identifier (letters, digits, underscores), so a hyphenated type such as `non-binding` can never appear in one.

## Postwork

<a id="relationship-authoring-best-practices-and-considerations"></a>
<a id="relationship-testing"></a>

### 5. Relationship Authoring Best Practices and Considerations

#### General

1. Use camelCase on the wire (`subType`, `mutatorRef`, `schemaVersion`).
2. Author against `relationships.meshery.io/v1beta3` (declare `v1beta2` while current servers require it - see the top of this page). Do not use `core.meshery.io/v1alpha2` or `v1beta1`.

#### Scoping

1. To apply a relationship across models, set the selector `model.name` to `*`. To limit it to one model, specify that model name (case sensitive).
2. Absence of a selector property is interpreted as the wildcard `*`.
3. Values for `kind`, `version`, and `model` are case-sensitive.

#### Actions

1. If a `mutatedRef` / `mutatorRef` path contains more than one array, only the first array position may be `_`; later arrays must be an explicit index (`0`, `1`, …).
2. `mutatedRef` currently does not support patching an array value itself.
3. Pair every mutator path with a mutated path. Do not swap source and sink.
4. Verify each path against the component schema in `models/<model>/.../components/<Kind>.json`.

#### Matching

1. Targets of a Relationship can be specific Components or entire Models.
2. Leave `evaluationQuery` empty (`""`); evaluation dispatches on `kind`, `type`, and `subType`.
3. `metadata.isAnnotation: true` means Meshery must not evaluate or patch the relationship.

#### Conflicts

1. Ensure `deny` selectors and `allow` selectors do not overlap for the same pair.
2. In the event of conflicting Relationship Definitions, the union between them is taken.
   - If we have two Relationships, one from (Component A) to (Component B and Component F), and another from (Component A) to (Component B and Component C), then it is similar to having a Relationship from Component A to Component B, C and F.
3. No relationship kind is inherently more important than another.

#### Schema Conformance

Every `mutatorRef`/`mutatedRef` path rooted at `configuration.` must resolve against the JSON schema of the component it addresses (`models/<model>/<version>/v1.0.0/components/<Kind>.json`). For a selector item belonging to the relationship's own model, that is the schema shipped in the **same** model version directory as the definition; for a selector item that references a **different** model - as a cross-model edge does - it is the schema in that model's **newest** version directory. A path that names a field the component schema does not define is written by the evaluation engine but never reaches the rendered resource - the defect behind [#21482](https://github.com/meshery/meshery/issues/21482), where an Ingress relationship patched the pre-1.22 `backend.serviceName` shape into a `networking.k8s.io/v1` component.

- **Exemptions.** The `configuration.metadata` subtree and paths rooted at `displayName` or `component` are not checked: component schemas describe the resource's `spec`, not its ObjectMeta or the Meshery component envelope.
- **Where it runs.** `server/policies/relationship_schema_conformance_test.go`, executed by the policies test workflow on every change under `models/**`.
- **When it fails, fix the path.** The `knownUnresolvedMutationPaths` allowlist exists only for known pre-existing defects, and cites the open follow-up issue tracking their repair ([#21490](https://github.com/meshery/meshery/issues/21490)). An allowlisted entry must keep failing, so repairing one of those definitions means deleting its line in the same pull request.

<a id="relationship-contribution"></a>

### 6. Contribute your relationship to the Meshery project

Submit a pull request to the Meshery repository with your new relationship definition, so that all users can benefit from the relationship(s) you have defined.

Keeping your relationship definition in a separate file allows for easier management and review of the relationship(s) you have defined.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Keeping your custom Relationships private</div>


Alternatively, if you would like to keep the relationship definition private, you can bundle your relationship(s) in a custom model and import the custom model into your Meshery deployment. Your private relationship definition(s) will be registered in your Meshery Server's <a href='/pr-preview/pr-21670/concepts/logical/registry/'>registry</a> and available for use within your Meshery deployment.
</div>


For more information refer - [Model - Construct Models in Meshery](https://docs.google.com/document/d/16z5hA8qVfSq885of9LXFUVvfom-hQXr-6oTD_GgoFmk/edit)
