# Design: Schema-Backed `Validate()` Methods on meshery/schemas Entity Types

**Status:** Draft
**Author:** Meshery CLI Working Group
**Date:** 2026-03-13
**Related PRs:** [meshery/meshery#17923](https://github.com/meshery/meshery/pull/17923), [meshery/meshkit#932](https://github.com/meshery/meshkit/pull/932)
**Target repos:** [meshery/meshery](https://github.com/meshery/meshery), [meshery/meshkit](https://github.com/meshery/meshkit), and optionally [meshery/schemas](https://github.com/meshery/schemas)

**Revision note:** This document was updated after `meshery/meshkit#932` merged. The original draft assumed that schema-driven validation in Go required either hand-written rule duplication or a heavy runtime schema-loading path. That assumption is no longer accurate.

---

## Executive Summary

The core problem still exists: `meshery/schemas` generates typed Go structs from OpenAPI v3, but the generated types do not themselves enforce required fields, enum membership, or other schema constraints. In Meshery today, that has led to hand-written validation logic in multiple consumers, including `mesheryctl`, with predictable drift.

What changed materially is the solution space. `meshery/meshkit#932` added a reusable `schema` package that validates Meshery documents against embedded schemas from `github.com/meshery/schemas`, with built-in schema registration, `schemaVersion`-based detection, compiled-schema caching, and structured validation errors. In other words, Meshery now has an authoritative runtime validator that uses the actual schemas and does not depend on filesystem access at runtime.

That changes the recommended design:

1. `github.com/meshery/meshkit/schema` should become the canonical validation engine for Meshery consumers such as `mesheryctl` and the server.
2. If ergonomic `Validate()` methods are still desired on typed structs in `meshery/schemas`, they should be thin adapters over the MeshKit validator, not hand-written re-implementations of schema rules.

This preserves the original goal of a shared validation surface while avoiding schema drift and taking advantage of the new MeshKit capability.

---

## 1. Problem Statement

### 1.1 The Constraint Gap in Generated Types

Meshery's schemas are defined in OpenAPI v3 YAML. Consider the `RelationshipDefinition` schema:

```yaml
# schemas/constructs/v1alpha3/relationship/relationship.yaml
RelationshipDefinition:
  type: object
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
      enum:
        - hierarchical
        - edge
        - sibling
```

`oapi-codegen` translates this into typed Go definitions, but the runtime constraints do not survive as executable validation logic:

```go
// models/v1alpha3/relationship/relationship.go  (auto-generated, do NOT edit)
type RelationshipDefinitionKind string

const (
    Edge         RelationshipDefinitionKind = "edge"
    Hierarchical RelationshipDefinitionKind = "hierarchical"
    Sibling      RelationshipDefinitionKind = "sibling"
)

type RelationshipDefinition struct {
    Id               uuid.UUID                  `json:"id"            yaml:"id"`
    Kind             RelationshipDefinitionKind `json:"kind"          yaml:"kind"`
    RelationshipType string                     `json:"type"          yaml:"type"`
    SchemaVersion    string                     `json:"schemaVersion" yaml:"schemaVersion"`
    SubType          string                     `json:"subType"       yaml:"subType"`
    Version          string                     `json:"version"       yaml:"version"`
    // ...
}
```

A `RelationshipDefinition` with `Kind: "foo"` and missing required fields is still a valid Go value at compile time and after unmarshal. The schema remains authoritative, but the generated type does not expose the schema as a callable validator.

### 1.2 Where This Surfaces in Meshery Today

Any code path that accepts untrusted or user-authored documents has to validate them. In Meshery today that still means duplicated, consumer-owned logic:

| Consumer | What it checks today | File |
| --- | --- | --- |
| `mesheryctl relationship validate` | `schemaVersion`, `kind` required + enum, `type`, `subType`, `version` | `mesheryctl/.../relationships/validate.go` |
| `mesheryctl model validate` | model name/version, component kind/version, relationship kind enum and required fields | `mesheryctl/.../model/validate.go` |
| Meshery server | partial validation, varies by path | spread across handlers and registration flow |
| External integrators | re-implement constraints independently | n/a |

The motivating review feedback on `meshery/meshery#17923` was exactly this kind of drift: one CLI path enforced the relationship-kind enum while another initially missed it.

### 1.3 What Changed with MeshKit `schema`

The original draft framed the problem as "there is no shared, importable definition of validity in Go." That is no longer true.

After `meshery/meshkit#932`, MeshKit now provides a shared runtime validator for Meshery documents. The remaining problem is narrower and more architectural:

1. Meshery consumers in this repository do not yet use the canonical validator.
2. `meshery/schemas` still does not expose an ergonomic typed `Validate()` surface on generated structs.
3. Product decisions such as whether CLI validation should be strict about server-generated fields are still unresolved.

The question is no longer whether Meshery can validate against the real schemas. It can. The question is how Meshery should expose and adopt that capability.

---

## 2. Current State in Detail

### 2.1 The Helper File Pattern in `meshery/schemas`

The schemas repository already uses `*_helper.go` files for hand-written methods attached to generated types:

```text
models/
├── v1alpha3/relationship/
│   ├── relationship.go          # Auto-generated - do NOT edit
│   └── relationship_helper.go   # Hand-written
├── v1beta1/model/
│   ├── model.go                 # Auto-generated
│   └── model_helper.go          # Hand-written
├── v1beta1/component/
│   ├── component.go             # Auto-generated
│   └── component_helper.go      # Hand-written
```

Those helper files already carry methods such as `TableName()`, `Type()`, `GenerateID()`, `Create(...)`, and serialization helpers. The documented pattern explicitly allows adding validation or business-logic helpers. So a `Validate()` method is still a natural API shape for `meshery/schemas` if the project wants one.

### 2.2 MeshKit's New Embedded Schema Validator

`meshery/meshkit#932` adds `github.com/meshery/meshkit/schema`, a reusable validation package built around embedded schemas from `github.com/meshery/schemas`.

At a high level, the package now provides APIs such as:

```go
func Validate(data []byte) error
func ValidateWithRef(ref Ref, data []byte) error
func ValidateAs(documentType DocumentType, data []byte) error
func DecodeAndValidate[T any](data []byte) (T, error)

func (v *Validator) ValidateAny(ref Ref, value any) error
func ValidationDetailsFromError(err error) (ValidationDetails, bool)
```

Key properties of this package matter directly to this design:

- It validates against the actual embedded schemas from `github.com/meshery/schemas`, not a hand-maintained rewrite.
- It registers built-in mappings for core Meshery document types and can auto-detect the schema from `schemaVersion`.
- It caches compiled schemas and guards concurrent compilation with `singleflight`.
- It returns structured validation information via MeshKit's `ErrorV2`, including `instancePath`, `schemaPath`, `keyword`, and `message` per violation.
- It does not require schema files to be present on disk at runtime.
- It uses `kin-openapi` rather than the older CUE-based validation path described in the original draft.

This directly eliminates the biggest argument in favor of hand-written validation logic: Meshery now has a shared, schema-driven implementation that is already packaged for Go consumers.

### 2.3 What Gap Remains After `meshkit/schema`

The new MeshKit package does not automatically solve every ergonomics question:

- It validates document bytes or normalized `any` values; it does not automatically attach `Validate()` to every generated Go struct.
- Consumers still need to decide how to render structured violations into CLI output or API responses.
- Strict schema conformance may reject author-authored drafts that omit fields the server normally generates, such as `id`.
- The current Meshery repository still depends on an older MeshKit release and therefore does not yet consume this package.

The capability gap is now about API shape and rollout, not about the absence of a canonical validator.

---

## 3. Proposed Solution

### 3.1 Make `meshkit/schema` the Canonical Validation Engine

All consumer-facing validation in Meshery should treat `github.com/meshery/meshkit/schema` as the source of truth.

That means `mesheryctl`, server-side document ingestion paths, and any future importers should validate incoming bytes or typed values through MeshKit's schema package instead of manually checking fields.

Example consumer-side usage:

```go
if err := meshkitschema.Validate(data); err != nil {
    if details, ok := meshkitschema.ValidationDetailsFromError(err); ok {
        for _, violation := range details.Violations {
            // Convert structured violations into CLI or API output.
        }
    }
    return err
}
```

This ensures every consumer executes the same schema logic against the same embedded artifacts.

### 3.2 If `Validate()` Exists on Typed Structs, It Should Delegate

If the project still wants an ergonomic method on `meshery/schemas` types, that method should be a thin wrapper over MeshKit's validator rather than a hand-written field checker.

Example shape for `relationship_helper.go`:

```go
// SchemaVersion is the canonical schemaVersion string for this package's
// generated types, derived from schemas/constructs/v1alpha3/relationship/templates/.
const SchemaVersion = "relationships.meshery.io/v1alpha3"

func (r RelationshipDefinition) Validate() error {
    return meshkitschema.Default().ValidateAny(meshkitschema.Ref{
        SchemaVersion: SchemaVersion,
        Type:          meshkitschema.TypeRelationship,
    }, r)
}
```

Important properties of this approach:

- The helper method stays ergonomic for callers.
- The actual validation remains schema-backed.
- New required fields, enums, patterns, and nested constraints are picked up automatically when the schema changes.
- The helper method does not become a second source of truth.

### 3.3 Preserve Structured Errors End-to-End

The return type should remain `error`, but the concrete error should be the MeshKit validation error so callers can recover structured details with `ValidationDetailsFromError(err)`.

The original draft proposed flattening violations into joined strings or inventing a new local `ValidationErrors` type. That is no longer necessary as a baseline design. MeshKit already supplies structured field-level diagnostics.

Guideline:

- Library and helper layers should preserve the original MeshKit error.
- Presentation layers such as `mesheryctl` may format those violations into human-friendly messages.
- Callers should not parse `err.Error()` when `ValidationDetailsFromError` is available.

### 3.4 Scope of Coverage

The immediate adoption targets in Meshery should be the current validation-heavy CLI paths:

| Area | Immediate target |
| --- | --- |
| `mesheryctl relationship validate` | Replace inline relationship checks with `meshkit/schema` |
| `mesheryctl model validate` | Replace inline model/component/relationship checks with `meshkit/schema` |
| Server registration/import paths | Move toward `meshkit/schema` as documents enter the system |

If optional helper methods are added in `meshery/schemas`, the first candidates remain:

| Type | Notes |
| --- | --- |
| `RelationshipDefinition` (v1alpha3) | Highest current duplication across CLI paths |
| `ModelDefinition` (v1beta1) | Used in `mesheryctl model validate` |
| `ComponentDefinition` (v1beta1) | Same |
| `DesignDefinition` (v1beta1) | Already supported by MeshKit's new schema package |

Unlike the original hand-written proposal, this approach is not limited to top-level required fields and enum membership. It inherits whatever the schema actually declares, including nested object rules, `const`, `pattern`, and future additions.

---

## 4. Alternatives Considered

### 4.1 Use MeshKit Directly, Without Adding `Validate()` to `meshery/schemas`

**How it works:** Consumers call `meshkit/schema` directly on bytes or typed values. No changes are made in `meshery/schemas`.

**Advantages:**

- Immediate use of the canonical validator with no new API surface in `meshery/schemas`
- No schema drift
- Structured validation errors already available
- No need to reason about helper-file dependency layering yet

**Disadvantages:**

- Callers that already hold typed structs do not get a convenient `Validate()` method
- Each consumer still needs a small adapter layer to turn `ValidationDetails` into its own UX

**Verdict:** This should be the default Phase 1 rollout regardless of whether helper methods are ever added.

### 4.2 Thin `Validate()` Wrappers in `meshery/schemas`

**How it works:** Helper methods are added to generated types, but they simply delegate to `meshkit/schema.ValidateAny` with explicit refs.

**Advantages:**

- Preserves the ergonomic API shape of `obj.Validate()`
- Keeps MeshKit as the single validation engine
- Avoids hand-written duplication while still improving typed-API usability

**Disadvantages:**

- Adds explicit dependency from versioned schema model packages to `meshkit/schema`
- Requires confirming there is no Go package-cycle issue and that the transitive dependency increase is acceptable

**Verdict:** Good optional Phase 2 if the ergonomics matter enough to justify the coupling.

### 4.3 Hand-Written `Validate()` Methods That Re-Encode Rules

**How it works:** The original design: helper files manually check required fields and enums with local Go code.

**Advantages:**

- Minimal runtime machinery inside `meshery/schemas`
- Full control over error wording and any intentional leniency

**Disadvantages:**

- Reintroduces schema drift immediately
- Misses nested constraints and future schema additions unless maintained by hand
- Duplicates capability that MeshKit now already provides

**Verdict:** No longer recommended as the canonical design. It is only defensible if dependency layering blocks wrapper-based delegation and the project still insists on a typed method inside `meshery/schemas`.

### 4.4 Generated Wrappers or Generated Validation

**How it works:** Extend code generation to emit `Validate()` methods or wrapper stubs from schema metadata.

**Advantages:**

- Reduces boilerplate if many schema types eventually gain helper methods
- Can still target MeshKit as the runtime engine instead of generating hand-written rule logic

**Disadvantages:**

- Requires codegen-template maintenance
- Premature until the desired public API shape is settled

**Verdict:** A reasonable future optimization after the runtime design is proven in real consumers.

### 4.5 Do Nothing / Leave Validation to Each Consumer

**Verdict:** Still the weakest option. Meshery already has evidence of drift in duplicated validation logic, and `meshkit/schema` now removes most of the old justification for leaving the problem unsolved.

---

## 5. Critical Analysis

### 5.1 The Original Draft's Baseline Assumptions Are Now Outdated

The original draft treated schema-driven validation as expensive, awkward, or operationally brittle in Go because it assumed runtime schema loading, CUE compilation, or custom validator wiring. `meshery/meshkit#932` materially changes that baseline. Any updated design must start from the fact that Meshery now has a first-class schema validator with embedded assets and structured diagnostics.

### 5.2 The Hard Question Is Now Validation Profile, Not Validation Mechanics

If the schema says `id` is required, the MeshKit validator will enforce it. That is correct for strict schema conformance, but it may not be the UX Meshery wants for hand-authored documents.

This means the important product decision is no longer "should validation be hand-written or schema-backed?" It is:

- Should `mesheryctl validate` mean strict canonical conformance?
- Or should it support a draft-authoring mode that tolerates server-populated fields and perhaps normalizes documents before validation?

That decision must be made explicitly. It should not be hidden inside a supposedly canonical `Validate()` method.

### 5.3 Structured Errors Are a Material Improvement and Should Not Be Thrown Away

MeshKit now returns structured `ValidationDetails` with per-field `Violation` entries. That is strictly better than a semicolon-joined string for servers, APIs, tests, and rich CLIs.

Flattening those violations too early would discard:

- the exact failing field path
- the schema keyword (`required`, `enum`, `pattern`, `const`, etc.)
- the schema location used for validation

Human-readable messages should be generated at the edge, not used as the underlying validation representation.

### 5.4 Helper-Method Coupling Is Real

If versioned packages in `meshery/schemas/models/...` import `meshkit/schema`, the dependency graph needs to be checked carefully. The design is attractive, but it should not be adopted blindly.

Two concrete questions need verification in the target repo:

- Does importing `meshkit/schema` from versioned schema model packages create any package-cycle issue?
- Is it acceptable for every Go consumer of those model packages to inherit the validation package's transitive dependencies even if they never call `Validate()`?

If the answer to either is no, then MeshKit should remain the consumer-side validation engine without adding helper methods to `meshery/schemas`.

### 5.5 Meshery Adoption Requires a Dependency Upgrade

This repository currently depends on an older MeshKit release and does not yet have `meshkit/schema` available. So the design is not just about code shape; it also includes a rollout dependency:

- upgrade MeshKit to a release that includes `meshery/meshkit#932`
- validate any knock-on dependency changes
- update CLI tests and output expectations accordingly

The good news is that MeshKit already ships package tests for the new validator, which reduces the amount of schema-validation logic Meshery itself needs to prove from scratch.

---

## 6. Implementation Plan

### Phase 1: Adopt MeshKit Schema Validation in Meshery Consumers

**Deliverable:** `mesheryctl` validation paths use `github.com/meshery/meshkit/schema` instead of inline field checks.

**Steps:**

1. Upgrade MeshKit in `meshery/meshery` to a version that includes `meshery/meshkit#932`.
2. Replace inline validation logic in `mesheryctl/internal/cli/root/relationships/validate.go` with MeshKit schema validation.
3. Replace inline validation logic in `mesheryctl/internal/cli/root/model/validate.go` for models, components, and relationships with MeshKit schema validation.
4. Add an adapter that converts `ValidationDetails` violations into the existing `validation.Result` output model.
5. Decide whether CLI validation is strict schema validation or a draft-friendly mode for server-generated fields such as `id`.
6. Add tests covering both success and failure paths using structured violations, not just string matching on hand-written errors.

**Estimated scope:** dependency bump plus focused CLI refactor; less custom validation code than the original proposal.

### Phase 2: Optional `Validate()` Methods in `meshery/schemas`

**Deliverable:** Typed helper methods exist only as thin delegators to MeshKit's schema package.

**Steps:**

1. Confirm package-graph viability in `meshery/schemas`.
2. Add `Validate() error` to high-value types such as `RelationshipDefinition`, `ModelDefinition`, and `ComponentDefinition`.
3. Ensure those methods return the original MeshKit error unchanged.
4. Add helper-method tests that assert `ValidationDetailsFromError(err)` still works on returned errors.
5. Update `meshery/schemas` contributor guidance to say that helper methods must delegate to the canonical validator rather than manually duplicate schema rules.

### Phase 3: Optional Codegen or Interface Refinement

If helper methods prove broadly useful, investigate generating thin wrapper methods from the schema/model inventory. If validation becomes a ubiquitous contract across entity types, only then consider whether an interface-level `Validate()` method belongs in MeshKit's entity abstractions.

That interface change remains separate because it is a broader API commitment and not required to adopt the new validator.

---

## 7. Acceptance Tests

The updated design should be considered complete only when the following are true:

1. A valid relationship document passes through `meshkit/schema.Validate(...)` without error.
2. An invalid relationship document with `kind: invalid` returns structured validation details that include `instancePath: /kind` and `keyword: enum`.
3. A document missing required fields returns schema-derived violations rather than consumer-specific hand-written checks.
4. `mesheryctl relationship validate` delegates to MeshKit schema validation and still produces clear user-facing output.
5. `mesheryctl model validate` delegates to MeshKit schema validation for models, components, and relationships.
6. The strict-versus-lenient behavior for server-generated fields such as `id` is explicitly tested according to the product decision taken.
7. If `Validate()` helper methods are added in `meshery/schemas`, those methods return errors from which `ValidationDetailsFromError` can still recover structured violations.
8. Meshery builds and tests pass with the upgraded MeshKit dependency on the touched paths.

---

## 8. Open Questions and Decisions Required

The following questions remain open after accounting for `meshery/meshkit#932`.

### 8.1 Should CLI Validation Be Strict or Draft-Friendly? **[Decision needed]**

With MeshKit as the canonical engine, the validator will enforce whatever the schema says, including server-generated required fields such as `id`.

**Proposed resolution:** Treat strict schema validation as the canonical behavior. If Meshery wants a draft-authoring workflow, add it explicitly as a separate mode (`lint`, `--lenient`, or a normalization step) rather than weakening the canonical validator.

### 8.2 Should `meshery/schemas` Expose `Validate()` at All? **[Decision needed]**

The new MeshKit package may make consumer-side adoption sufficient. `obj.Validate()` is ergonomically appealing, but it is no longer required to achieve shared validation.

**Proposed resolution:** Adopt MeshKit in consumers first. Add helper methods only if there is repeated adapter code or clear external demand for typed ergonomics.

### 8.3 If Helper Methods Exist, Should They Preserve MeshKit Errors or Wrap Them? **[Decision needed]**

**Proposed resolution:** Preserve the original MeshKit error. Helper methods should not flatten or repackage validation details unless there is a compelling downstream compatibility requirement.

### 8.4 Is the Package Graph Acceptable for Helper Delegation? **[Decision needed]**

Before adding helper wrappers in `meshery/schemas`, the project should explicitly verify dependency direction and import-cycle safety.

**Proposed resolution:** Validate the package graph in the target repository as a Phase 2 prerequisite. If the coupling is undesirable, keep validation in MeshKit consumers only.

### 8.5 Multi-Version Strategy **[Decision needed]**

Schema-backed validation is version-sensitive. `RelationshipDefinition` v1alpha2 and v1alpha3 are different schemas and should not share a silent fallback.

**Proposed resolution:** Consumer-side validation should rely on explicit `schemaVersion` or explicit refs. Optional helper methods should remain version-specific and use explicit refs defined by their package.

### 8.6 Semver and Rollout **[Advisory]**

Using MeshKit's new schema package in Meshery is a consumer-side dependency upgrade. Adding helper `Validate()` methods in `meshery/schemas` would be additive and semver-compatible. Changing common interfaces to require `Validate()` would still be a separate, potentially breaking step.

### 8.7 Long-Term Generation Path **[Advisory]**

If the project eventually wants `Validate()` on many schema types, generation should target thin delegate methods or schema/ref registration helpers, not hand-written re-encodings of `required:` and `enum:` blocks.

---

## References

- [meshery/meshery#17923](https://github.com/meshery/meshery/pull/17923) - The Meshery PR that surfaced the immediate validation drift
- [meshery/meshkit#932](https://github.com/meshery/meshkit/pull/932) - Adds embedded Meshery schema validation and structured validation errors
- [meshkit `schema/validator.go`](https://github.com/meshery/meshkit/blob/master/schema/validator.go) - Core validator API and document-validation flow
- [meshkit `schema/error.go`](https://github.com/meshery/meshkit/blob/master/schema/error.go) - `ValidationDetails` and structured error extraction
- [meshery/schemas README - Helper Files](https://github.com/meshery/schemas/blob/master/README.md#-go-helper-files) - Established pattern for hand-written helper methods
- [meshery/schemas CONTRIBUTING](https://github.com/meshery/schemas/blob/master/CONTRIBUTING.md) - Schema development workflow
- [Contributing to Schemas - docs.meshery.io](https://docs.meshery.io/project/contributing/contributing-schemas) - Official schema contributor guide
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - Go code generation tool used by `meshery/schemas`
- [protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate) - Prior art for generated validation APIs
