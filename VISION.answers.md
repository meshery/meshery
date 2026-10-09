# Meshery vision review answers

This file preserves the author's recorded decisions that calibrate `VISION.md`.

## Round 1 - 2026-08-26

### H-1 - A generic resource importer without a model

Verdict: In vision.

Reasoning: Users can use annotation-only components and annotation-only relationships (like those offered in the meshery-core and meshery-shapes models) to visually depict their full architecture or to visually and collaboratively convey a concept with the acknowledgement that those annotation-only components and relationships are not semantically meaningful; they do not bear configuration properties that can be orchestrated or deployed.

Reasoning: That annotation-only components are excluded from deployment operations.

### H-2 - A marketplace extension with a trust banner

Verdict: In vision.

### H-3 - A temporary local API endpoint

Verdict: Off mission.

Reasoning: We need to ratchet away from the local API endpoints.

Reasoning: Meshery is a schema-driven project.

Reasoning: The source of truth for all API endpoints shall be sourced from https://github.com/meshery/schemas.

### H-4 - Canonicalize credentials on the next save

Verdict: In vision.

### H-5 - Keep a permission-denied design entry visible

Verdict: Off mission.

### H-6 - Allow an unverified explicit chart pin

Verdict: In vision.

### H-7 - Skip missing protected-branch credentials

Verdict: In vision.

### H-8 - Run relationship integration tests nightly

Verdict: Conditional.

Reasoning: Not now.

Reasoning: Perhaps, in the future.

### H-9 - Make annotation edges operational by default

Verdict: Off mission.

### H-10 - Relax a conservative SQL lint rule globally

Verdict: Off mission.

## Round 2 - 2026-08-27

### H-11 - Curate the extension marketplace

Verdict: In vision.

### H-12 - Auto-update marketplace extensions

Verdict: In vision.

### H-13 - Consume a schemas prerelease

Verdict: In vision.

### H-14 - Drop unknown fields during credential canonicalization

Verdict: In vision.

### H-15 - Allow skipped provider coverage to merge

Verdict: In vision.

### H-16 - Scope relationship integration by changed path

Verdict: In vision.

Reasoning: Skip on documentation-only changes and provider-ui-only changes, but execute on all other PRs.

### H-17 - Verify an air-gapped pin from local evidence

Verdict: In vision.

### H-18 - Allow a generated-code safety baseline

Verdict: In vision.

### H-19 - Let remote providers publish private schemas

Verdict: In vision.

### H-20 - Use encrypted production-derived test fixtures

Verdict: In vision.

### A-1 - Approve the revised acceptance policy

Verdict: In vision.

## Verbatim queued verdicts

```text
H-1 (A generic resource importer without a model): In vision. Reasoning: Users can use annotation-only components and annotation-only relationships (like those offered in the meshery-core and meshery-shapes models) to visually depict their full architecture or to visually and collaboratively convey a concept with the acknowledgement that those annotation-only components and relationships are not semantically meaningful; they do not bear configuration properties that can be orchestrated or deployed. That annotation-only components are excluded from deployment operations.
H-2 (A marketplace extension with a trust banner): In vision.
H-3 (A temporary local API endpoint): Off mission. Reasoning: We need to ratchet away from the local API endpoints. Meshery is a schema-driven project. The source of truth for all API endpoints shall be sourced from https://github.com/meshery/schemas.
H-4 (Canonicalize credentials on the next save): In vision.
H-5 (Keep a permission-denied design entry visible): Off mission.
H-6 (Allow an unverified explicit chart pin): In vision.
H-7 (Skip missing protected-branch credentials): In vision.
H-8 (Run relationship integration tests nightly): Conditional. Reasoning: Not now. Perhaps, in the future.
H-9 (Make annotation edges operational by default): Off mission.
H-10 (Relax a conservative SQL lint rule globally): Off mission.
All 10 hypotheticals are answered. Fold every verdict into the draft and send back the next round.
```

## Round 2 verbatim queued verdicts

```text
H-11 (Curate the extension marketplace): In vision.
H-12 (Auto-update marketplace extensions): In vision.
H-13 (Consume a schemas prerelease): In vision.
H-14 (Drop unknown fields during credential canonicalization): In vision.
H-15 (Allow skipped provider coverage to merge): In vision.
H-16 (Scope relationship integration by changed path): In vision. Reasoning: Skip on documentation-only changes and provider-ui-only changes, but execute on all other PRs.
H-17 (Verify an air-gapped pin from local evidence): In vision.
H-18 (Allow a generated-code safety baseline): In vision.
H-19 (Let remote providers publish private schemas): In vision.
H-20 (Use encrypted production-derived test fixtures): In vision.
A-1 (Approve the revised acceptance policy): In vision.
All 11 hypotheticals are answered. Fold every verdict into the draft and send back the next round.
```

## Changelog

- H-1 In vision -> `Design carries operational meaning` now permits explicitly non-semantic, annotation-only components and relationships that carry no deployment configuration and are excluded from deployment operations.
- H-2 In vision -> `Interoperability is schema-driven` now permits a signed third-party extension marketplace when activation makes the extension trust boundary explicit.
- H-3 Off mission -> `Interoperability is schema-driven` now names meshery/schemas as the source of every HTTP API contract and refuses local endpoint or contract copies.
- H-4 In vision -> `Intent changes visibly` now permits canonicalization of legacy credentials on an intentional save while keeping tolerant reads during transition.
- H-5 Off mission -> `Access and outcomes are explicit` now requires inaccessible entry points to be hidden rather than leading to a permission-denied page.
- H-6 In vision -> `Intent changes visibly` now permits an air-gapped explicit chart pin after a typed acknowledgement of responsibility.
- H-7 In vision -> `Access and outcomes are explicit` now treats unavailable provider credentials as a visibly skipped CI project, including on protected branches.
- H-8 Conditional -> `Proof follows real use` retains pull-request relationship integration coverage today and leaves nightly-only coverage for a future decision.
- H-9 Off mission -> `Design carries operational meaning` now states that user-created annotation relationships never become operational by default.
- H-10 Off mission -> `Access and outcomes are explicit` now prefers documented narrow safety exceptions over global relaxation of a security gate.
- H-11 In vision -> `Interoperability is schema-driven` now requires maintainer review for marketplace publishers.
- H-12 In vision -> `Interoperability is schema-driven` now permits automatic update of enabled marketplace extensions to their latest verified version at restart with a post-update changelog.
- H-13 In vision -> `Interoperability is schema-driven` now permits a signed, expiring meshery/schemas prerelease for coordinated cross-repository work.
- H-14 In vision -> `Intent changes visibly` now drops fields absent from canonical credential schema during intentional canonicalization and logs the migration.
- H-15 In vision -> `Access and outcomes are explicit` now allows a visible skipped provider check to remain non-blocking for merge.
- H-16 In vision -> `Proof follows real use` now runs relationship integration scenarios on every pull request except documentation-only and provider-ui-only changes.
- H-17 In vision -> `Intent changes visibly` now accepts locally uploaded chart and checksum evidence for an air-gapped explicit pin.
- H-18 In vision -> `Access and outcomes are explicit` now permits a documented package-level generated-code SQL safety baseline only when its generator test proves the safety property.
- H-19 In vision -> `Interoperability is schema-driven` now permits provider-owned private schema constructs instead of local endpoint copies.
- H-20 In vision -> `Proof follows real use` now permits encrypted, redacted production-derived design fixtures from a controlled test store.
- A-1 In vision -> the author approved the fully folded `VISION.md` as Meshery's current acceptance policy.
