---
name: 💡 General Feature Request
about: Suggest an enhancement to Meshery.
title: 'Adopt schema-backed `Validate()` in Meshery consumers; add thin `Validate() error` adapters to entity helper files'
labels: 'kind/enhancement'
assignees: ''
---

### Current Behavior

The `meshery/schemas` package generates strongly-typed Go structs from OpenAPI v3 schemas using `oapi-codegen`, but the generated code discards all validation constraints — required fields, enum membership, string formats — declared in the source schemas.

For example, `relationship.yaml` defines:

```yaml
RelationshipDefinition:
  required:
    - schemaVersion
    - version
    - model
    - kind
    - type
    - subType
    - id
  properties:
    kind:
      type: string
      enum: [hierarchical, edge, sibling]
```

The generated `relationship.go` produces a `RelationshipDefinitionKind` string alias with three constants, but nothing prevents a value of `"foo"` from being unmarshalled and used without error.

Every Go consumer that needs to validate these constraints — the CLI, the server, external importers — must independently re-implement them. This has already caused bugs: in [meshery/meshery#17923](https://github.com/meshery/meshery/pull/17923), the `model validate` command initially missed the `kind` enum check that `relationship validate` had, because each command owned its own copy of the rules.

**What changed:** `meshery/meshkit#932` merged and added `github.com/meshery/meshkit/schema`, a reusable validation package that validates Meshery documents against the actual embedded schemas from `github.com/meshery/schemas`. This package uses `kin-openapi` for validation, embeds the schema assets so no filesystem access is needed at runtime, and returns structured field-level `Violation` diagnostics. The canonical Go validator now exists. What remains is adopting it in consumers and optionally exposing an ergonomic typed API on entity structs.

### Desired Behavior

**In Meshery consumers (e.g. `mesheryctl`):** validation paths should delegate to `github.com/meshery/meshkit/schema` instead of re-implementing field checks locally. Structured per-field diagnostics should be surfaced in CLI output.

```go
// Before (mesheryctl today): hand-written field checks per command
if rel.Kind == "" {
    violations = append(violations, "kind is required")
}
// ... many more checks, duplicated across commands

// After: one call to the canonical validator
if err := meshkitschema.Validate(data); err != nil {
    if details, ok := meshkitschema.ValidationDetailsFromError(err); ok {
        for _, v := range details.Violations {
            // v.InstancePath, v.Keyword, v.Message are available
        }
    }
    return err
}
```

**In `meshery/schemas` (optional, ergonomic layer):** each entity definition type may gain a thin `Validate() error` method in its existing `*_helper.go` file. This method delegates to `meshkit/schema` rather than duplicating rules:

```go
// relationship_helper.go
func (r RelationshipDefinition) Validate() error {
    return meshkitschema.Default().ValidateAny(meshkitschema.Ref{
        Type: meshkitschema.TypeRelationship,
    }, r)
}

// Usage
rel := relationship.RelationshipDefinition{Kind: "foo"}
if err := rel.Validate(); err != nil {
    // err wraps structured ValidationDetails with per-field Violations
    if details, ok := meshkitschema.ValidationDetailsFromError(err); ok {
        // details.Violations[0].InstancePath == "/kind"
        // details.Violations[0].Keyword == "enum"
    }
}
```

The canonical list of constraints stays in the schema YAML. `Validate()` helper methods never duplicate schema rules; they only provide a typed call site.

### Implementation

A full design document covering all alternatives and trade-offs is available here: [design-validate-methods-schemas.md](https://github.com/meshery/meshery/blob/master/docs/design-validate-methods-schemas.md)

**Phase 1 — Adopt `meshkit/schema` in Meshery consumers:**

1. Upgrade `meshery/meshery` to a MeshKit release that includes `meshery/meshkit#932`.
2. Replace inline validation logic in `mesheryctl/internal/cli/root/relationships/validate.go` with `meshkitschema.Validate(data)`.
3. Replace inline validation logic in `mesheryctl/internal/cli/root/model/validate.go` for models, components, and relationships with `meshkitschema.Validate(data)`.
4. Add an adapter that converts `meshkitschema.ValidationDetails` violations into the existing CLI output model.
5. Decide explicitly whether CLI validation is strict schema conformance or a draft-friendly mode that tolerates server-generated fields such as `id`.
6. Add tests covering both success and failure paths, asserting on structured `Violation` fields rather than string-matching hand-written messages.

**Phase 2 — Optional `Validate()` adapters in `meshery/schemas`:**

Add thin `Validate() error` to the helper files for each primary entity type, delegating to `meshkit/schema.Default().ValidateAny(...)`:

| File | Entity | `meshkitschema.Ref` |
|---|---|---|
| `models/v1alpha3/relationship/relationship_helper.go` | `RelationshipDefinition` | `Ref{Type: TypeRelationship}` |
| `models/v1beta1/model/model_helper.go` | `ModelDefinition` | `Ref{Type: TypeModel}` |
| `models/v1beta1/component/component_helper.go` | `ComponentDefinition` | `Ref{Type: TypeComponent}` |

Prerequisites before Phase 2:
- Confirm the `meshery/schemas → meshkit/schema` import direction does not create a package cycle.
- Confirm the transitive dependency additions to `meshery/schemas` are acceptable.
- Ensure `Validate()` returns the original MeshKit error unchanged so callers can extract `ValidationDetails` via `meshkitschema.ValidationDetailsFromError(err)`.

**Phase 3 — Optional codegen or interface refinement (future):**

If thin helper methods prove broadly valuable across many schema types, investigate generating them automatically. If validation becomes a universal contract, evaluate whether `Validate() error` should be added to the `entity.Entity` interface in MeshKit — this is a semver-major interface change and must be coordinated separately.

### Acceptance Tests

**Phase 1 — MeshKit schema adoption in consumers:**
- [ ] A valid relationship document file passes `mesheryctl relationship validate` with no errors
- [ ] A relationship document with `kind: invalid` produces output that identifies the failing field (not a generic error message)
- [ ] A relationship document missing a required field produces output that names the missing field
- [ ] `mesheryctl model validate` produces equivalently structured output for model, component, and relationship violations
- [ ] The strict-versus-lenient behavior for server-generated fields such as `id` is explicitly tested according to the product decision recorded in the design document

**Phase 2 — `Validate()` adapters in `meshery/schemas`:**
- [ ] `rel.Validate()` on a fully schema-conformant `RelationshipDefinition` returns `nil`
- [ ] Same for `ModelDefinition` and `ComponentDefinition`
- [ ] `rel.Validate()` on a `RelationshipDefinition{Kind: "foo"}` returns a non-nil error from which `meshkitschema.ValidationDetailsFromError` recovers a `Violation` with `InstancePath: "/kind"` and `Keyword: "enum"`
- [ ] `rel.Validate()` on a `RelationshipDefinition` with `Version: ""` returns a non-nil error mentioning the missing `version` field
- [ ] Equivalent negative-path coverage for each required field of `ModelDefinition` and `ComponentDefinition`
- [ ] `go test ./...` passes in the `meshery/schemas` repository with the new helper methods

**Integration:**
- [ ] At least one consumer (`mesheryctl`) delegates to the canonical validator rather than inline field checks

**Note on scope:** Phase 2 helper methods must not contain hand-written re-implementations of schema rules. Any deviation from schema-backed validation is a prior CLI bug, not a feature. `id` handling is a product decision (see design document Section 8.1) and must be addressed explicitly, not silently.

---
### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://meshery.io/community/handbook)
- 🛠 [Meshery Build & Release Strategy](https://docs.meshery.io/project/contributing/build-and-release)
- 📚 [Instructions for contributing to documentation](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow)
   - Meshery documentation [site](https://docs.meshery.io/) and [source](https://github.com/meshery/meshery/tree/master/docs)
- 🎨 Wireframes and designs for Meshery UI in [Figma](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](http://discuss.meshery.io) and [Community Slack](https://slack.meshery.io)
