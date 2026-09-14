---
layout: default
title: "RFC: Relationship Evaluator Migration from Rego to Rust/WASM"
abstract: "This document outlines why Meshery's relationship evaluator should be migrated from OPA/Rego to Rust compiled to WebAssembly, enabling both server-side and client-side evaluation."
permalink: project/contributing/relationship-evaluator-migration-rfc
type: project
category: contributing
language: en
---

# RFC: Migrating the Relationship Evaluator from Rego to Rust/WASM

| **Status** | Proposed |
| **Authors** | Meshery Maintainers |
| **Created** | 2025 |
| **Last Updated** | 2025 |

## Executive Summary

Meshery's relationship evaluator, currently implemented in OPA/Rego, suffers from fundamental performance limitations, poor developer experience, and maintenance challenges. After thorough analysis, we conclude that:

1. **Rego is the wrong tool for this job** — The evaluator implements a graph algorithm, not a policy decision system
2. **A Rust rewrite offers significant benefits** — 10-100x performance improvement, type safety, proper debugging
3. **Our use case doesn't need a policy engine** — Relationship definitions already provide the "hot-loadable" configuration we need
4. **A hybrid approach would add complexity without benefit** — Context-switching and maintenance overhead outweigh any gains
5. **Rust compiled to WASM enables universal deployment** — Same binary runs in Go server and browser, enabling client-side evaluation

This RFC proposes a complete migration to Rust, compiled to WebAssembly (WASM) for integration with both the Go server and the React/JavaScript UI.

---

## Table of Contents

1. [Background](#background)
2. [Problems with the Current Rego Implementation](#problems-with-the-current-rego-implementation)
3. [Why Rego is the Wrong Tool](#why-rego-is-the-wrong-tool)
4. [Why We Don't Need a Policy Engine](#why-we-dont-need-a-policy-engine)
5. [Why Rust is the Right Choice](#why-rust-is-the-right-choice)
6. [Why WASM for Integration](#why-wasm-for-integration)
7. [Why Not a Hybrid Approach](#why-not-a-hybrid-approach)
8. [Proposed Architecture](#proposed-architecture)
9. [Migration Plan](#migration-plan)
10. [Appendix: Code Analysis](#appendix-code-analysis)

---

## Background

### What is the Relationship Evaluator?

The relationship evaluator is a core component of Meshery that:

- **Identifies relationships** between components in a design (e.g., Pod → Service bindings, Namespace → Deployment hierarchies)
- **Validates relationships** to ensure they remain valid as the design changes
- **Generates actions** to mutate the design (add/update/delete components and relationships)
- **Applies mutations** to produce an updated design

### Current Implementation

The evaluator is implemented in ~1,500 lines of Rego policy code across 15+ files in:

```
server/meshmodel/meshery-core/0.7.2/v1.0.0/policies/
├── evaluation.rego          # Main entry point (~200 lines)
├── evaluator.rego            # Generic evaluation framework
├── alias_policy.rego         # Alias relationship handling
├── matchlabels-policy.rego   # Match labels relationship handling
├── edge_non_binding.rego     # Edge network relationships
├── hierarchical_parent_child_policy.rego
├── binding_policy.rego
├── eval_rules.rego           # Common evaluation rules
├── actions.rego              # Action definitions
├── utils.rego                # Utility functions
├── helper_rules.rego
├── feasibility_evaluation.rego
└── ...
```

---

## Problems with the Current Rego Implementation

### 1. Severe Performance Issues: O(n²) to O(n⁴) Complexity

The code contains nested loops that scale quadratically or worse:

**Match Labels Identification (O(n² × f²)):**
```rego
identify_matchlabels(design_file, relationship) := all_match_labels if {
    field_pairs := {pair |
        some component in design_file.components           # O(n)
        some other_component in design_file.components     # O(n)
        ...
        some field, value in configuration_for_component_at_path(...)   # O(f)
        some field2, value2 in configuration_for_component_at_path(...) # O(f)
    }
}
```

**Relationship Identification (O(s² × n²)):**
```rego
identify_relationships_based_on_matching_mutator_and_mutated_fields(...) := {declaration |
    some selector_set in relationship.selectors        # O(s)
    some from_selector in selector_set.allow.from      # O(s)
    some to_selector in selector_set.allow.to          # O(s)
    some component_from in design_file.components      # O(n)
    some component_to in design_file.components        # O(n)
    ...
}
```

For a design with 100 components and 10 fields each, this results in **100 × 100 × 10 × 10 = 1,000,000 iterations** for a single policy.

### 2. Excessive Memory Allocation via JSON Patching

The code uses `json.patch` extensively, creating full copies of the design at each step:

```rego
design_file_to_evaluate := json.patch(final_design_file_old, [...])
design_file_with_validated_rels := actions.apply_all_actions_to_design(design_file_to_evaluate, ...)
design_file_with_identified_rels := actions.apply_all_actions_to_design(design_file_to_evaluate, ...)
design_file_to_apply_actions := actions.apply_all_actions_to_design(design_file_with_identified_rels, ...)
design_to_return := actions.apply_all_actions_to_design(design_file_to_apply_actions, ...)
```

Each `json.patch` call creates a complete copy of the design object. A single evaluation can create 10+ full copies.

### 3. No Indexing Capability

Every component lookup requires a full scan:

```rego
component_declaration_by_id(design_file, id) := component if {
    some component in design_file.components  # O(n) scan
    component.id == id
}
```

This function is called repeatedly during evaluation. With proper indexing (HashMap), this would be O(1).

### 4. Poor Developer Experience

- **No debugger**: Development relies on `print()` statements scattered throughout the code
- **No type system**: Runtime errors from missing fields, wrong types, or null values
- **No profiler**: Cannot identify performance bottlenecks
- **Limited testing**: OPA's test framework is awkward compared to standard unit testing

### 5. Redundant Computation

The main `evaluate` rule recomputes similar data multiple times and has visible bugs:

```rego
design_file_with_identified_rels := actions.apply_all_actions_to_design(
    design_file_to_evaluate,    # BUG: Should use design_file_with_validated_rels
    validation_actions,
)
```

### 6. Non-Deterministic Behavior

UUID generation uses `time.now_ns()`, defeating Rego's memoization and making tests flaky:

```rego
new_uuid(seed) := id if {
    now := format_int(time.now_ns(), 10)
    id := uuid.rfc4122(sprintf("%s%s", [seed, now]))
}
```

---

## Why Rego is the Wrong Tool

### What Rego is Designed For

OPA/Rego is a **policy decision engine** designed for:

| Use Case | Example |
|----------|---------|
| Authorization | "Can user X access resource Y?" |
| Admission Control | "Should this Kubernetes pod be allowed?" |
| Configuration Validation | "Does this config meet compliance requirements?" |
| Simple Allow/Deny | "Is this API request permitted?" |

Rego excels when:
- Decisions are **stateless** (input → decision, no mutation)
- Logic is **declarative** (what to check, not how)
- Policies are **external** (loaded at runtime, separate from application)
- Answers are **boolean** or simple values

### What the Relationship Evaluator Actually Does

| Operation | Is This a Policy Decision? | Rego Suitability |
|-----------|---------------------------|------------------|
| Find matching component pairs | No — graph traversal | Poor |
| Build relationship graph | No — data structure construction | Poor |
| Track relationship state (pending/approved/deleted) | No — state machine | Poor |
| Generate mutation actions | No — transformation | Poor |
| Apply mutations to design | No — state change | Very Poor |
| Batch and optimize changes | No — algorithm | Poor |

**The relationship evaluator is implementing a graph algorithm, not making policy decisions.**

### Fundamental Mismatch

| Requirement | Rego's Model | What We Need |
|-------------|--------------|--------------|
| Data Structures | Immutable, no indexing | Mutable, indexed (HashMap, Graph) |
| Iteration | Find all matches | Find first, early termination |
| State | Stateless | Track state across phases |
| Mutation | Returns new copy | In-place modification |
| Control Flow | Declarative search | Imperative control |
| Optimization | Engine-controlled | Developer-controlled |

---

## Why We Don't Need a Policy Engine

### The "Hot-Loadable Policies" Argument

One argument for keeping Rego is that policy engines allow hot-loading policies at runtime without recompilation. However, **this benefit doesn't apply to Meshery's use case**.

#### We Already Have Hot-Loadable Configuration: Relationship Definitions

Relationship definitions are already stored as data in the registry:

```
server/meshmodel/
├── kubernetes/
│   └── relationships/
│       ├── hierarchical_namespace_to_deployment.json
│       ├── edge_service_to_pod.json
│       └── ...
├── istio/
│   └── relationships/
│       └── ...
└── ...
```

These definitions specify:
- Which component kinds can relate
- What type of relationship (hierarchical, edge, sibling)
- Selector patterns for matching
- Mutation references

**The "policy" is in the relationship definitions, not the Rego code.**

The Rego code is the **evaluation engine** that interprets these definitions. The engine itself doesn't need to be hot-loadable — it's core application logic.

#### Analogy

| Component | Hot-Loadable? | Equivalent |
|-----------|---------------|------------|
| Relationship Definitions | Yes (data) | Database schema, config files |
| Rego Evaluator | No (but Rego allows it) | Application code |
| Rust Evaluator | No | Application code |

Swapping the Rego evaluator for Rust doesn't reduce flexibility — the flexibility comes from the relationship definitions, which remain data-driven.

### Other Policy Engine Benefits We Don't Use

| OPA/Rego Feature | Used by Meshery? | Notes |
|------------------|------------------|-------|
| External policy bundles | No | Policies versioned with meshery-core |
| Decision logging/audit | No | Not implemented |
| Policy discovery | No | Policies are static |
| Partial evaluation | No | Full evaluation every time |
| Multi-tenancy isolation | No | Single evaluator instance |

---

## Why Rust is the Right Choice

### Direct Comparison

| Aspect | Current Rego | Rust Implementation |
|--------|--------------|---------------------|
| Component lookup | O(n) scan | O(1) HashMap |
| Relationship identification | O(n²) nested loops | O(n) with indices |
| Memory allocation | Full copy per mutation | In-place or COW |
| Type safety | None (runtime errors) | Compile-time guarantees |
| Debugging | `print()` statements | Full debugger (lldb, gdb) |
| Profiling | Not possible | perf, flamegraph, criterion |
| Parallelism | Not possible | Rayon (trivial parallelism) |
| Testing | OPA test framework | Standard Rust testing + proptest |
| IDE support | Limited | Full (rust-analyzer) |
| Error handling | Silent failures | `Result<T, E>` |

### Performance Expectations

| Operation | Rego (Current) | Rust (Expected) | Improvement |
|-----------|----------------|-----------------|-------------|
| Component lookup | 1-10ms | <0.001ms | 1000-10000x |
| Identify relationships (100 components) | 100-500ms | 1-5ms | 100x |
| Full evaluation | 500ms-2s | 5-20ms | 50-100x |

### Why Rust Over Go?

While Go is already used in Meshery's server, Rust offers advantages for this specific component:

| Factor | Rust | Go |
|--------|------|-----|
| Performance | Faster (no GC pauses) | Fast (but GC overhead) |
| Memory safety | Compile-time (ownership) | Runtime (GC) |
| Parallelism | Fearless concurrency | Goroutines (good, but less control) |
| Type system | More expressive (enums, traits) | Simpler |
| Error handling | `Result<T, E>` (explicit) | `error` (can be ignored) |
| **WASM compilation** | **Excellent (first-class support)** | **Limited (large binaries, no wasm-bindgen)** |

For a performance-critical component processing potentially thousands of components, Rust's zero-cost abstractions and lack of garbage collection make it the better choice.

**Critically, Rust has first-class WASM support**, which enables running the same evaluator in both the Go server and the browser — a capability that Go cannot match effectively.

---

## Why WASM for Integration

### The Universal Runtime Opportunity

Compiling Rust to WebAssembly (WASM) enables a powerful architectural pattern: **write once, run everywhere**.

```
┌─────────────────────────────────────────────────────────────┐
│                  Rust Evaluator Core                         │
│         (pure Rust, no platform-specific code)               │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │    Cargo Build    │
                    └─────────┬─────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
   │wasm32-wasip1│     │wasm32-unknown│    │ Native      │
   │ (for Go     │     │-unknown      │    │ (optional,  │
   │  via wazero)│     │(for browser) │    │  benchmarks)│
   └──────┬──────┘     └──────┬──────┘     └─────────────┘
          │                   │
          ▼                   ▼
   ┌─────────────┐     ┌─────────────────────────────┐
   │  Go Server  │     │       Browser               │
   │  (wazero -  │     │  (native WASM support)      │
   │   pure Go,  │     │  via wasm-bindgen           │
   │   no CGO)   │     │                             │
   └─────────────┘     └─────────────────────────────┘
```

### Benefits of WASM Integration

| Benefit | Description |
|---------|-------------|
| **Single codebase** | One Rust implementation for Go server AND browser |
| **Guaranteed consistency** | Same binary = identical evaluation results everywhere |
| **Client-side evaluation** | Real-time feedback during design editing (no server roundtrip) |
| **Reduced server load** | Offload evaluation to client for interactive editing |
| **Offline capability** | Users can edit designs without server connection |
| **No CGO complexity** | wazero is pure Go, no C compiler needed |

### Client-Side Evaluation: The UX Win

This is a significant architectural opportunity that WASM unlocks:

```
Current Flow (Server-Only Rego):
┌────────────┐    HTTP     ┌────────────┐
│  Browser   │ ──────────► │  Go Server │
│  (React)   │ ◄────────── │  (Rego)    │
└────────────┘  ~100-500ms └────────────┘
       User waits for every relationship update

With WASM (Client-Side):
┌────────────────────────────────────┐
│           Browser                  │
│  ┌────────┐      ┌──────────────┐  │
│  │ React  │ ───► │ WASM         │  │
│  │ Canvas │ ◄─── │ Evaluator    │  │
│  └────────┘      └──────────────┘  │
│              ~5-15ms               │
└────────────────────────────────────┘
       Instant feedback as user designs
```

**Use cases enabled:**
- User drags Pod onto Namespace → **instantly** see hierarchical relationship
- User connects Service to Deployment → **instantly** see edge relationship
- Preview relationship changes before saving
- Work offline, sync later

### Go Integration: wazero

[wazero](https://wazero.io/) is the ideal WASM runtime for Go:

| Feature | Benefit |
|---------|---------|
| Pure Go | No CGO, no external dependencies, simple builds |
| Zero dependencies | Just Go standard library |
| WASI support | Standard interface for WASM modules |
| Production ready | Used by Envoy, Dapr, Trivy, and other major projects |
| Good performance | ~70-90% of native for compute-bound tasks |

```go
// Example Go integration with wazero
import "github.com/tetratelabs/wazero"

type Evaluator struct {
    runtime wazero.Runtime
    module  wazero.CompiledModule
}

func (e *Evaluator) Evaluate(design, relationships []byte) (*EvaluationResult, error) {
    // Call WASM function - same binary as browser uses
    resultBytes, err := e.callWasm("evaluate_design", design, relationships)
    if err != nil {
        return nil, err
    }
    
    var result EvaluationResult
    json.Unmarshal(resultBytes, &result)
    return &result, nil
}
```

### JavaScript Integration: wasm-bindgen

For browser integration, we use `wasm-bindgen` and `wasm-pack`:

```rust
// Rust side with wasm-bindgen
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub fn evaluate_design(design_json: &str, relationships_json: &str) -> Result<String, JsValue> {
    let design: Design = serde_json::from_str(design_json)
        .map_err(|e| JsValue::from_str(&e.to_string()))?;
    
    let relationships: Vec<RelationshipDef> = serde_json::from_str(relationships_json)
        .map_err(|e| JsValue::from_str(&e.to_string()))?;
    
    let result = evaluate(&design, &relationships);
    
    serde_json::to_string(&result)
        .map_err(|e| JsValue::from_str(&e.to_string()))
}
```

```typescript
// React/TypeScript integration
import init, { evaluate_design } from '@meshery/relationship-evaluator';

// Initialize once at app startup
await init();

// Use in React component for real-time evaluation
function useRelationshipEvaluation(design: Design, relationships: RelationshipDef[]) {
    const [result, setResult] = useState<EvaluationResult | null>(null);
    
    useEffect(() => {
        // Runs instantly in browser - no server call needed
        const resultJson = evaluate_design(
            JSON.stringify(design),
            JSON.stringify(relationships)
        );
        setResult(JSON.parse(resultJson));
    }, [design, relationships]);
    
    return result;
}
```

### Bundle Size Considerations

| Metric | Estimated Size | Notes |
|--------|----------------|-------|
| Uncompressed WASM | 500KB - 1.5MB | Depends on features |
| Gzipped | 150KB - 400KB | Standard HTTP compression |
| With `wasm-opt -Oz` | 20-40% smaller | Part of build pipeline |
| Brotli compressed | Even smaller | Modern browsers support |

This is acceptable for Meshery's web app. The WASM module can be:
- Lazy-loaded when user opens design editor
- Cached by browser/service worker
- Loaded in parallel with other assets

### Performance Expectations with WASM

| Environment | vs Native Rust | Expected Latency (100 components) |
|-------------|----------------|-----------------------------------|
| Native Rust | 100% | 5-10ms |
| WASM in Browser | 70-90% | 7-14ms |
| WASM in Go (wazero) | 60-80% | 8-16ms |
| **Current Rego** | **~1%** | **500-2000ms** |

All WASM environments are **50-100x faster than current Rego**.

### Client-Server Synchronization

Since evaluation can happen both client-side and server-side:

```
┌─────────────────────────────────────────────────────────────┐
│                    Design Editing Flow                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. User edits design in browser                            │
│                     │                                        │
│                     ▼                                        │
│  2. Client-side WASM evaluation (instant feedback)          │
│     - Show preview relationships                             │
│     - Highlight validation issues                            │
│     - Update canvas in real-time                            │
│                     │                                        │
│                     ▼                                        │
│  3. User saves design                                        │
│                     │                                        │
│                     ▼                                        │
│  4. Server-side WASM evaluation (authoritative)             │
│     - Validate client results (same binary = same results)  │
│     - Persist to database                                    │
│     - Sync back to client if any edge cases                 │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

Since it's the **same WASM binary**, results are guaranteed identical (deterministic evaluation). The server-side evaluation serves as:
- Authoritative validation before persistence
- Protection against tampered clients
- Single source of truth for saved designs

---

## Why Not a Hybrid Approach

### The Hybrid Proposal

A hybrid approach would keep Rego for "policy decisions" while moving heavy computation to Rust:

```
Go Server → Rust (identification) → Rego (policy decisions) → Rust (actions)
```

### Why This Doesn't Work for Meshery

#### 1. Context-Switching Cost

Developers would need to:
- Understand two programming paradigms (declarative vs. imperative)
- Maintain two codebases with different tooling
- Debug across language boundaries
- Write integration tests spanning both systems

#### 2. Serialization Overhead

Every call between Rust and Rego requires:
- Serializing the design to JSON
- Passing to OPA runtime
- Deserializing results
- For a design with 100 components, this adds 5-15ms per boundary crossing

#### 3. No Clear Boundary

Where would you draw the line?

| Boundary Option | Problem |
|-----------------|---------|
| Rego for feasibility checks | These are O(1) already; overhead > benefit |
| Rego for allow/deny | We don't have meaningful deny policies |
| Rego for custom user policies | Not a current requirement; can add later if needed |

#### 4. Rego Isn't Adding Value

As analyzed above, none of the evaluator's operations are natural policy decisions. Keeping Rego means paying the costs of a hybrid system without receiving the benefits.

#### 5. Maintenance Burden

| Aspect | Single Codebase (Rust) | Hybrid (Rust + Rego) |
|--------|------------------------|----------------------|
| Documentation | 1 system | 2 systems |
| CI/CD | 1 pipeline | 2 pipelines |
| Dependencies | Rust toolchain | Rust + OPA runtime |
| Expertise required | Rust | Rust + Rego |
| Onboarding | Learn 1 system | Learn 2 systems + integration |

### Future Extensibility

If custom user policies become a requirement in the future:
- Add a simple YAML/JSON deny-rule system
- Implement a WASM plugin interface
- These are simpler than maintaining a full Rego integration

---

## Proposed Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────────┐
│                     Rust Relationship Evaluator                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    Design Index Layer                       │ │
│  │  ┌─────────────────┐    ┌────────────────────────────────┐ │ │
│  │  │ ComponentIndex  │    │ RelationshipGraph              │ │ │
│  │  │ - by_id: HashMap│    │ - adjacency: HashMap<Id, Vec>  │ │ │
│  │  │ - by_kind: Multi│    │ - by_status: HashMap           │ │ │
│  │  │ - by_model: Multi│   └────────────────────────────────┘ │ │
│  │  └─────────────────┘                                        │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                   Policy Evaluators                        │ │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────┐ │ │
│  │  │AliasEvaluator│ │MatchLabels  │ │HierarchicalEvaluator │ │ │
│  │  │              │ │Evaluator    │ │                      │ │ │
│  │  └──────────────┘ └──────────────┘ └──────────────────────┘ │ │
│  │  ┌──────────────┐ ┌──────────────┐                          │ │
│  │  │EdgeNetwork   │ │BindingPolicy│                          │ │
│  │  │Evaluator     │ │Evaluator    │                          │ │
│  │  └──────────────┘ └──────────────┘                          │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                   Action Collector                          │ │
│  │  - Batches mutations                                        │ │
│  │  - Deduplicates actions                                     │ │
│  │  - Applies in optimal order                                 │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   Go Server     │
                    │ (via wazero     │
                    │  WASM runtime)  │
                    └─────────────────┘
```

### Core Data Structures

```rust
/// Indexed design for O(1) lookups
pub struct IndexedDesign {
    components: HashMap<ComponentId, Component>,
    components_by_kind: HashMap<String, Vec<ComponentId>>,
    components_by_model: HashMap<String, Vec<ComponentId>>,
    relationships: HashMap<RelationshipId, Relationship>,
    relationships_by_status: HashMap<Status, Vec<RelationshipId>>,
}

/// Policy evaluator trait
pub trait PolicyEvaluator {
    fn identify(&self, design: &IndexedDesign, definitions: &[RelationshipDef]) -> Vec<Relationship>;
    fn validate(&self, design: &IndexedDesign, relationship: &Relationship) -> ValidationResult;
    fn generate_actions(&self, design: &IndexedDesign, relationship: &Relationship) -> Vec<Action>;
}

/// Action to apply to design
pub enum Action {
    AddComponent(Component),
    UpdateComponent { id: ComponentId, path: JsonPath, value: Value },
    DeleteComponent(ComponentId),
    AddRelationship(Relationship),
    UpdateRelationship { id: RelationshipId, path: JsonPath, value: Value },
    DeleteRelationship(RelationshipId),
}
```

### Integration Strategy: WASM

Both Go server and browser use the same WASM binary:

**Go Server (via wazero):**
```go
import (
    "context"
    "github.com/tetratelabs/wazero"
    "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmEvaluator struct {
    runtime wazero.Runtime
    module  wazero.CompiledModule
}

func NewWasmEvaluator(wasmBytes []byte) (*WasmEvaluator, error) {
    ctx := context.Background()
    r := wazero.NewRuntime(ctx)
    wasi_snapshot_preview1.MustInstantiate(ctx, r)
    
    module, err := r.CompileModule(ctx, wasmBytes)
    if err != nil {
        return nil, err
    }
    
    return &WasmEvaluator{runtime: r, module: module}, nil
}

func (e *WasmEvaluator) Evaluate(design, relationships []byte) ([]byte, error) {
    // Call into WASM module
    return e.callExport("evaluate_design", design, relationships)
}
```

**Browser (via wasm-bindgen):**
```typescript
// NPM package: @meshery/relationship-evaluator
import init, { evaluate_design, EvaluationResult } from '@meshery/relationship-evaluator';

class RelationshipEvaluator {
    private initialized = false;
    
    async initialize(): Promise<void> {
        if (!this.initialized) {
            await init();
            this.initialized = true;
        }
    }
    
    evaluate(design: Design, relationships: RelationshipDef[]): EvaluationResult {
        const resultJson = evaluate_design(
            JSON.stringify(design),
            JSON.stringify(relationships)
        );
        return JSON.parse(resultJson);
    }
}

export const evaluator = new RelationshipEvaluator();
```

**Build Pipeline:**
```bash
# Build for Go server (WASI target)
cargo build --target wasm32-wasip1 --release
wasm-opt -Oz target/wasm32-wasip1/release/evaluator.wasm -o evaluator.wasm

# Build for browser (with JS bindings)
wasm-pack build --target web --release
# Outputs: pkg/meshery_relationship_evaluator.js + .wasm
```

---

## Migration Plan

### Phase 1: Design & Prototype (2 weeks)
- Define Rust data structures
- Implement IndexedDesign with benchmarks
- Prototype one policy evaluator (e.g., HierarchicalParentChild)
- Validate WASM compilation for both targets

### Phase 2: Core Implementation (4 weeks)
- Implement all policy evaluators
- Implement action collector and application
- Unit tests with property-based testing
- Ensure deterministic evaluation (no time-based UUIDs)

### Phase 3: WASM Integration (3 weeks)
- Set up wasm-pack build pipeline for browser target
- Set up wazero integration for Go server
- Create NPM package `@meshery/relationship-evaluator`
- Integration tests verifying identical results in Go and browser
- Performance benchmarking in both environments

### Phase 4: UI Integration (2 weeks)
- Integrate WASM evaluator into React canvas
- Implement real-time relationship preview
- Add loading states and error handling
- Lazy-load WASM module

### Phase 5: Server Migration (2 weeks)
- Feature flag for new evaluator in Go server
- Shadow mode (run both Rego and WASM, compare results)
- Gradual rollout with monitoring

### Phase 6: Cleanup (1 week)
- Remove Rego code
- Update documentation
- Archive old implementation

**Total: ~14 weeks**

---

## Appendix: Code Analysis

### Files Analyzed

| File | Lines | Primary Issues |
|------|-------|----------------|
| `evaluation.rego` | ~200 | Monolithic, redundant computation, excessive patching |
| `evaluator.rego` | ~80 | Stub functions, unclear partial rule pattern |
| `alias_policy.rego` | ~180 | Complex nested logic, O(n) lookups |
| `matchlabels-policy.rego` | ~130 | O(n² × f²) complexity |
| `eval_rules.rego` | ~180 | O(s² × n²) identification, no early termination |
| `utils.rego` | ~170 | Repeated O(n) scans, no caching |
| `actions.rego` | ~100 | Reasonable, but excessive json.patch |

### Performance Profile (Estimated)

For a design with 100 components and 50 existing relationships:

| Phase | Current Time | With Rust |
|-------|--------------|-----------|
| Build indices | N/A (none) | 1ms |
| Identify relationships | 200-400ms | 2-5ms |
| Validate relationships | 50-100ms | 1-2ms |
| Generate actions | 20-50ms | <1ms |
| Apply actions | 50-100ms | 1-2ms |
| **Total** | **320-650ms** | **5-10ms** |

---

## Conclusion

The relationship evaluator's implementation in Rego represents a fundamental mismatch between tool and task. Rego is designed for policy decisions; the evaluator implements a graph algorithm. The result is:

- **Poor performance** from O(n²)+ complexity and lack of indexing
- **Poor developer experience** from no debugger, no types, no profiler
- **Maintenance burden** from complex, interconnected rules

A Rust implementation compiled to WASM offers:
- **50-100x performance improvement**
- **Type safety and proper tooling**
- **Cleaner, more maintainable code**
- **Universal deployment** — same binary in Go server and browser
- **Client-side evaluation** — instant feedback during design editing
- **Reduced server load** — offload interactive evaluation to browser
- **Offline capability** — evaluate designs without server connection

The "hot-loadable policies" benefit of Rego doesn't apply because relationship definitions already provide the data-driven configuration layer. A hybrid Rego approach would add complexity without meaningful benefit.

**Recommendation: Proceed with a full Rust implementation of the relationship evaluator, compiled to WebAssembly for integration with both the Go server (via wazero) and the browser (via wasm-bindgen).**

---

## References

- [OPA Documentation](https://www.openpolicyagent.org/docs/latest/)
- [Rego Language Reference](https://www.openpolicyagent.org/docs/latest/policy-language/)
- [Meshery Architecture](https://docs.meshery.io/concepts/architecture)
- [Relationship Definitions](https://docs.meshery.io/concepts/logical/relationships)
- [wazero - Zero Dependency WebAssembly Runtime for Go](https://wazero.io/)
- [wasm-bindgen - Facilitating Rust and JavaScript Communication](https://rustwasm.github.io/docs/wasm-bindgen/)
- [wasm-pack - Rust to WASM Build Tool](https://rustwasm.github.io/wasm-pack/)
- [WebAssembly Specification](https://webassembly.github.io/spec/core/)