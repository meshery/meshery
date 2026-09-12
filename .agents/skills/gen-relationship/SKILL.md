---
name: gen-relationship
description: Create and refine schema-backed Meshery relationship definitions (models/**/relationships/*.json). Use for "add a relationship", "fix this relationship JSON", classifying kind/type/subType, or reviewing selectors and mutatorRef/mutatedRef patches between model components.
---

# Create and refine Meshery relationship definitions

Relationships describe how components in the same or different models relate, and whether they affect each other (semantic) or only the designer's mental model (non-semantic).

This skill covers **creating new definitions** and **refining existing ones**. Invoke it for `/gen-relationship`, "add a relationship", "fix this relationship JSON", or "what kind/type/subType should this be?".

## Source of truth

**[meshery/schemas](https://github.com/meshery/schemas) is SSOT.** Do not invent fields, enums, or path shapes.

| What | Where |
|---|---|
| Relationship object | `schemas/constructs/v1beta3/relationship/relationship.yaml` |
| Selectors, patch, mutator/mutated | `schemas/constructs/v1beta3/relationship/api.yml` |
| Starter document | `schemas/constructs/v1beta3/relationship/templates/relationship_template.json` |
| Human docs | [Concepts: Relationships](https://docs.meshery.io/concepts/logical/relationships), [Contributing to Relationships](https://docs.meshery.io/project/contributing/contributing-relationships) |
| In-tree corpus | `models/<model>/<model-version>/<definition-version>/relationships/` |

`v1beta3` is the authoring target; in-tree model files are still mostly `v1beta2` (a few `v1alpha3`). The version shapes are compatible - Meshery Server bridges registered definitions to the `v1beta2` shape for the policy engine (`server/models/pattern/utils/relationship_version_bridge.go`) - **but the registration gate accepts only `v1beta2`/`v1alpha3` documents until [meshkit#1096](https://github.com/meshery/meshkit/pull/1096) ships in the server's meshkit**. A definition that must register on current servers (in-tree seeding, `mesheryctl model import`) declares `v1beta2`; flip to `v1beta3` once the gate accepts it. When **refining** an existing file, keep its `schemaVersion` unless you are deliberately migrating it.

`kind` is a schema enum: `hierarchical` | `edge` | `sibling`. `type` and `subType` are open strings. The combinations below are the ones Meshery currently visualizes and evaluates. A new `subType` needs a visual paradigm (whiteboard a proposal in Kanvas) and an evaluation policy that understands it.

Required object fields (v1beta3): `schemaVersion`, `version`, `model`, `kind`, `type`, `subType`, `id`. In practice every useful definition also carries `status`, `metadata` (with `description`), `selectors`, and an empty `evaluationQuery` — start from the registry template or an [example](examples/), not the bare minimum.

## When to use which combo

`kind` + `type` + `subType` together pick the visual paradigm and the evaluation policy. Copy an existing combo; do not remix. Open only the example file for the combo you picked - the others add nothing for your case.

| kind | type | subType | Meaning | Example | File |
|---|---|---|---|---|---|
| `edge` | `non-binding` | `reference` | Logical name/id pointer; no traffic, mount, or RBAC | Deployment → ConfigMap | [examples/edge-non-binding-reference.json](examples/edge-non-binding-reference.json) |
| `edge` | `non-binding` | `network` | Documented L3/L4/L7 selection, no provisioned attachment | Service → Deployment | [examples/edge-non-binding-network.json](examples/edge-non-binding-network.json) |
| `edge` | `non-binding` | `firewall` | Policy that allows/denies traffic between peers | NetworkPolicy → Pod | [examples/edge-non-binding-firewall.json](examples/edge-non-binding-firewall.json) |
| `edge` | `non-binding` | `permission` | Mentions an identity or role without binding it | IAM Policy → Role | [examples/edge-non-binding-permission.json](examples/edge-non-binding-permission.json) |
| `edge` | `non-binding` | `alias` | Named stand-in, not nested ownership | Secret → Secret | [examples/edge-non-binding-alias.json](examples/edge-non-binding-alias.json) |
| `edge` | `non-binding` | `annotation` | Designer-only line; `metadata.isAnnotation: true`; no patch | Shape → Shape | [examples/edge-non-binding-annotation.json](examples/edge-non-binding-annotation.json) |
| `edge` | `non-binding` | `inventory` | Rare. Peer index/list, not parent-child containment | Cluster → Database | [examples/edge-non-binding-inventory.json](examples/edge-non-binding-inventory.json) |
| `edge` | `binding` | `permission` | Assigns identities (Role/RoleBinding/ServiceAccount) | Role → ServiceAccount via RoleBinding | [examples/edge-binding-permission.json](examples/edge-binding-permission.json) |
| `edge` | `binding` | `mount` | Storage or device is attached | PVC → Pod | [examples/edge-binding-mount.json](examples/edge-binding-mount.json) |
| `edge` | `binding` | `network` | Rare. Connecting provisions/rewrites network identity | Listener → Domain | [examples/edge-binding-network.json](examples/edge-binding-network.json) |
| `hierarchical` | `parent` | `inventory` | Parent scopes/contains children; parent identity patched onto child | `*` → Namespace | [examples/hierarchical-parent-inventory.json](examples/hierarchical-parent-inventory.json) |
| `hierarchical` | `parent` | `alias` | Child is a nested object inside the parent | Container → Pod | [examples/hierarchical-parent-alias.json](examples/hierarchical-parent-alias.json) |
| `hierarchical` | `parent` | `wallet` | Child config is held/patched into the parent | WASMFilter → EnvoyFilter | [examples/hierarchical-parent-wallet.json](examples/hierarchical-parent-wallet.json) |
| `hierarchical` | `sibling` | `matchlabels` | **In-tree tagsets encoding.** Shared labels; no patch | `*` ↔ `*` | [examples/hierarchical-sibling-matchlabels.json](examples/hierarchical-sibling-matchlabels.json) |
| `sibling` | `matchlabels` | `tagsets` | **Schema-native tagsets.** `kind` enum value `sibling` | `*` ↔ `*` | [examples/sibling-matchlabels.json](examples/sibling-matchlabels.json) |

`badge` appears in contributing docs as a visual paradigm only. There is no in-tree encoding. Do not invent `subType: badge`.

### Sibling encoding caveat

The schema `kind` enum is `hierarchical | edge | sibling`. Kubernetes in-tree files use `kind: hierarchical`, `type: sibling`, `subType: matchlabels` (see `models/kubernetes/.../relationships/sibling-tagsets.json`).

- **Refining kubernetes (or copying its pattern):** keep `hierarchical` / `sibling` / `matchlabels`.
- **New model, standardizing on the schema kind:** `kind: sibling` as in `examples/sibling-matchlabels.json`.
- Never mix both encodings in one model.

### Binding vs non-binding

- `binding`: forming the relationship assigns, mounts, or entitles (lifecycle effect).
- `non-binding`: the link is real (selector, name reference, policy match) but does not itself provision the attachment.

Most Service → workload edges are `non-binding` / `network`. Mount and RBAC assignment are `binding`.

### Semantic vs annotation

`metadata.isAnnotation: true` (or a shape/annotation component) means Meshery must **not** evaluate or patch. Edge-annotation is the usual encoding. Do not add `mutatorRef`/`mutatedRef` on annotation relationships.

## Hierarchical `from` / `to`

**`from` is the child. `to` is the parent.** Always.

| subType | Data flow |
|---|---|
| `inventory` | Parent **mutates** the child (Namespace name → child's `metadata.namespace`) |
| `alias` | Child **mutates** the parent (Container spec → Pod `spec.containers._`) |
| `wallet` | Child **mutates** the parent (filter config → EnvoyFilter patch) |

The previous gen-relationship example that put Namespace in `from` was wrong.

## Actions: `mutatorRef` and `mutatedRef`

Defined in `api.yml` as **nested arrays of string path segments** (`string[][]`). Sequence length on both sides must match: index `i` of `mutatorRef` patches onto index `i` of `mutatedRef`.

```json
"mutatorRef": [["config", "url"], ["config", "name"]],
"mutatedRef": [["configPatch", "value"], ["name"]]
```

`[config, url]` is copied onto `[configPatch, value]`; `[config, name]` onto `[name]`.

| Field | Role |
|---|---|
| `mutatorRef` | **Source.** JSON path of the value to read. |
| `mutatedRef` | **Sink.** Path segments of the field to patch. |
| `patchStrategy` | How to apply. Schema enum: `merge`, `strategic`, `add`, `remove`, `replace`, `copy`, `move`, `test`. The in-tree corpus and the evaluation engine use `replace` exclusively; default to `replace` unless you need different semantics. (`replace` joined the v1beta3 enum in [meshery/schemas#1166](https://github.com/meshery/schemas/pull/1166).) |

Paths are relative to the **component document** as Meshery stores it (`configuration`, `displayName`, `component.kind`, …), not the raw Kubernetes YAML root.

Path rules from contributing docs (still true):

1. `_` may mark **only the first** array position in a path. Later arrays must be an explicit index (`0`, `1`, …).
2. `mutatedRef` currently does not support an array as the patched value itself.
3. Pair every mutator path with a mutated path. Do not leave an unpaired side.
4. Verify each path against the **component schema** in `models/<model>/.../components/<Kind>.json` (or the schemas construct). Do not guess.

Some selectors use a nested `match` block (`match.from` / `match.to` / `match.refs`) instead of, or in addition to, `patch`. Tagsets use `match.refs` on labels. Permission/mount often use `match` with `kind: self` plus `mutatorRef`/`mutatedRef` on the match item. When refining, copy the structure already used by that combo in kubernetes.

Omit `patch` entirely when the relationship only matches (tagsets, annotation).

## Selectors

`selectors` (plural) is an array of selector-set items. Each item has `allow` and optional `deny`, each with `from` and `to`.

- One selector-set item = one independent matchmaking rule (OR across items).
- Inside one item, every `from` entry relates to every `to` entry - a cross-product, not 1:many.
- `deny` subtracts matches. `allow` and `deny` must not describe the same pair.
- Missing selector fields imply `*` (wildcard). Values for `kind`, `model`, `version` are case-sensitive.
- `model.name: "*"` on a selector item lets the relationship span models; the relationship-level `model` still names the owning model.

Wrong names you will see in stale docs: `selector` (singular), `core.meshery.io/v1alpha2`. Use `selectors` and `relationships.meshery.io/v1beta3`.

## `evaluationQuery`

Deprecated (the schema says so). Set it to `""`, as every one of the ~3,860 in-tree definitions does. The evaluation engine enters through the fixed policy package `data.relationship_evaluation_policy` and dispatches on `kind`/`type`/`subType`; it never reads this field.

If legacy tooling forces you to name a rule, the historical default is `{kind}_{subType}_relationship` (lowercase, from `GetDefaultEvaluationQuery()` in meshery/schemas) — note there is no `{type}` segment, and a Rego rule name allows only letters, digits, and underscores, so a hyphenated type such as `non-binding` can never appear in one.

## Create a new definition

1. Identify the two (or more) component kinds and their models. Confirm they exist under `models/`.
2. Decide whether they **affect** each other (semantic) or only the diagram (annotation).
3. Pick `kind` / `type` / `subType` from the table. If none fit, stop and propose a visualization; do not invent a combo.
4. Set `from` / `to` (hierarchical: child in `from`, parent in `to`).
5. If values should flow, add paired `mutatorRef` / `mutatedRef` after reading the component JSON schemas. If they should only match, use `match.refs` or omit patch.
6. Set `evaluationQuery: ""`, `status: enabled`, and the `schemaVersion` the registration gate accepts (`v1beta2` today; `v1beta3` once meshkit#1096 ships - see Source of truth).
7. Write one JSON file per relationship (or a small cohesive set) under `models/<model>/<model-version>/<def-version>/relationships/`. Filename convention: `{kind}-{type}-{subType}-<suffix>.json`.
8. Compare against the matching file in [examples/](examples/) and against a kubernetes in-tree neighbour of the same combo.

## Refine an existing definition

1. Open the in-tree file. Note `schemaVersion`, `kind`, `type`, `subType`. Do not "upgrade" those three unless the classification was wrong.
2. Keep `from`/`to` direction. Hierarchical files that already have child-in-`from` must stay that way.
3. Widen or narrow `allow`/`deny` kinds rather than cloning a near-duplicate relationship, unless the patch paths truly differ.
4. When adding a component kind to `from` or `to`, copy an existing selector item and change `kind` plus patch paths. Re-verify paths.
5. Do not swap `mutatorRef` and `mutatedRef`. Source vs sink is the usual bug.
6. Conflicting overlapping definitions union. Complementary overlaps union. Do not add a second relationship that restates the same pair with a different kind.
7. Leave nil UUIDs (`00000000-0000-0000-0000-000000000000`) as-is; the registry assigns real ids.

## Common mistakes

- `schemaVersion` `v1beta1` or `core.meshery.io/v1alpha2` — wrong. New: `relationships.meshery.io/v1beta3`.
- Putting the parent in `from` for hierarchical inventory.
- Using `type: network` with no `subType`, or treating `inventory` as a `kind`.
- `selector` instead of `selectors`.
- Flat paths (`"spec.containers[0].name"`) instead of nested arrays (`[["configuration","spec","containers","_","name"]]`).
- Unpaired mutator/mutated sequences.
- Guessing JSON paths instead of reading the component schema.
- Copying v1beta2 boilerplate blindly into a v1beta3 file (or the reverse) without checking required fields.
- Mixing `kind: sibling` with in-tree `kind: hierarchical, type: sibling` in the same model.

## Validate and register

Machine-check every definition before shipping it; do not stop at eyeballing.

1. `python3 -m json.tool <file>` — parses.
2. Package it: `mesheryctl model init <model> --version <ver>` scaffolds `<model>/<ver>/{model.json,components/,relationships/}`; drop your files in `relationships/` (copy the in-tree `model.json`), then `mesheryctl model build <model>/<ver> --path .` produces an OCI tar and fails on malformed documents.
3. Register it: `mesheryctl model import -f <model>-<ver>.tar` against a running server.
4. **Verify registration yourself — never trust the import summary.** Query `/api/meshmodels/models/<model>/relationships` (or `mesheryctl relationship view <model>`) and confirm your definitions appear. Until [meshery/schemas#1169](https://github.com/meshery/schemas/pull/1169) ships, importing relationships into a model the server already knows silently orphans them (`model_id` = nil UUID) while the summary reports success.
5. Prove behavior with an evaluation: POST a design whose component configurations already satisfy the relationship (matching names/paths) to `/api/meshmodels/relationships/evaluate` and confirm your definition appears in the response with `status: approved`.

The strict document validator (meshkit `schema.Validate`) currently rejects corpus-conventional files for reasons unrelated to your authoring (see [meshery/schemas#1167](https://github.com/meshery/schemas/issues/1167)); `model build` + `import` + evaluation are the gates that matter today.

## Output

Return the full relationship JSON (or a unified diff against the existing file). State the `kind` / `type` / `subType` you chose and why, name the mutator and mutated paths, and cite the schema files you used. If the combo is new, stop and say so.
