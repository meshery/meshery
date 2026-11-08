# Relationship definition examples

One compact, pedagogical example per canonical `kind` / `type` / `subType` combination found in Meshery models, plus the schema-native sibling encoding.

These are teaching fixtures, not drop-in replacements for in-tree `models/**/relationships/*.json`. In-tree files are mostly `relationships.meshery.io/v1beta2` and carry full model/registrant envelopes. These examples declare `relationships.meshery.io/v1beta3` (the authoring target); a definition that must register on current servers declares `v1beta2` until meshkit#1096 ships in the server - the rest of the document is identical.

See `SKILL.md` for when to pick each combo and how `mutatorRef` / `mutatedRef` pair.

A `model` object inside a selector item is a model *reference* (the schema's `ModelReference`): its outer `version` is the model-definition version, while the nested `model.version` object carries the upstream model version (for example, the Kubernetes release). The `"model": { ... "model": { "version": "" } }` nesting is intentional, not a copy-paste error.

| File | kind | type | subType |
|---|---|---|---|
| `edge-non-binding-reference.json` | edge | non-binding | reference |
| `edge-non-binding-network.json` | edge | non-binding | network |
| `edge-non-binding-firewall.json` | edge | non-binding | firewall |
| `edge-non-binding-permission.json` | edge | non-binding | permission |
| `edge-non-binding-alias.json` | edge | non-binding | alias |
| `edge-non-binding-annotation.json` | edge | non-binding | annotation |
| `edge-non-binding-inventory.json` | edge | non-binding | inventory |
| `edge-binding-permission.json` | edge | binding | permission |
| `edge-binding-mount.json` | edge | binding | mount |
| `edge-binding-network.json` | edge | binding | network |
| `hierarchical-parent-inventory.json` | hierarchical | parent | inventory |
| `hierarchical-parent-alias.json` | hierarchical | parent | alias |
| `hierarchical-parent-wallet.json` | hierarchical | parent | wallet |
| `hierarchical-sibling-matchlabels.json` | hierarchical | sibling | matchlabels |
| `sibling-matchlabels.json` | sibling | matchlabels | tagsets |

`badge` is named in docs as a visual paradigm but has no in-tree `kind`/`type`/`subType` encoding. Do not invent `subType: badge` unless you also propose a visualization and an evaluation policy.
