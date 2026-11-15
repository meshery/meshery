# Vision

Meshery exists so that platform engineers, site reliability engineers, and DevSecOps teams can collaboratively design and operate cloud native infrastructure.
It serves teams managing infrastructure and workloads across Kubernetes clusters, clouds, and integrations, and it turns their declared designs into evaluated, deployable, and observable operations.
It owns exactly one thing: the relationship-aware control surface between an infrastructure design and its safe operation.

## Design carries operational meaning

A Meshery design represents components and relationships that the platform can evaluate, render, and operate.
Meshery distinguishes relationships that guide lifecycle behavior from annotations that only aid human understanding.
Meshery welcomes annotation-only components and relationships for visual and collaborative representation when they are explicitly non-semantic.
Annotation-only components and relationships carry no deployment configuration and are excluded from deployment operations.
Every new visual object states whether it changes lifecycle semantics.
Meshery never makes user-created annotation relationships operational by default.
Meshery makes dependency behavior explicit so a failed component cannot make dependent work appear successful.

## Interoperability is schema-driven

Meshery grows integrations through designated models, adapters, providers, APIs, and UI extension points.
Meshery distributes signed third-party extensions through a marketplace only after maintainers review each publisher and administrator activation makes the trust boundary explicit.
Meshery automatically updates an enabled marketplace extension to its latest verified version on the next restart and gives administrators a changelog afterward.
Meshery keeps provider-owned identity, authorization, and durable persistence behind the provider boundary.
Every HTTP API endpoint sources its contract from meshery/schemas or a provider-owned private schema construct.
Meshery uses a generated schema client and type rather than a local endpoint or local contract copy.
Meshery treats an enabled extension as code inside the deployment's trust boundary.

## Intent changes visibly

Meshery honors an operator's explicit configuration or refuses it with an actionable diagnostic.
Meshery corrects only derived defaults against verified published state and reports each correction.
Meshery permits an air-gapped administrator to use an explicit chart pin with locally uploaded chart and checksum evidence or a typed acknowledgement of responsibility.
Meshery preserves legacy persisted data during a canonical form transition through tolerant reads rather than silent rewrites.
Meshery canonicalizes legacy credential data on an intentional save, drops fields absent from the canonical schema, logs the migration, and retains tolerant reads during the transition.
Meshery centralizes format resolution at shared layer boundaries instead of scattering shape assumptions across readers.

## Access and outcomes are explicit

Meshery gates protected content as well as protected controls and prevents denied content from loading.
Meshery hides entry points that a user cannot access rather than directing users to a permission-denied page.
Meshery reserves an explicit post-login onboarding exception rather than allowing unbounded ungated pages.
Meshery enforces configured provider selection at server startup rather than relying on a chooser, cookie, header, query parameter, or page visibility.
Meshery marks a provider CI project as skipped when required credentials are unavailable, including on protected branches, keeps that coverage visible, and does not let the visible skip block a merge.
Meshery makes test and lint gates reflect actual command failures.
Meshery permits a documented package-level SQL safety baseline for generated code only when its generator test proves it cannot emit an untrusted ORDER BY value.

## Proof follows real use

Meshery tests relationship policy against complete designs and real registrant shapes.
Meshery may use encrypted and redacted production-derived design fixtures from a controlled test store when they improve realistic coverage.
Meshery adds a regression test that fails against the reported defect before it claims a fix.
Meshery runs complete relationship integration scenarios on every pull request except documentation-only and provider-ui-only changes.
Meshery runs locally the checks that decide whether a change can merge.
Meshery documents a contract or safety boundary where contributors will encounter it.

## Scope

Meshery is not a replacement for a provider's identity or authorization system.
Meshery is not a security sandbox for extensions.
Meshery is not a generic annotation canvas whose every relationship changes lifecycle state.
Meshery is not a blind data-migration service or an unbounded compatibility layer.
Meshery is not a source of hand-written HTTP API contracts without a schemas-based source.
Meshery does not turn an unsupported explicit deployment choice into a different choice without telling the operator.

A change aligns when it improves how a declared design is evaluated, operated, authorized, extended through a named contract, or verified against the shape operators actually use.
A change should be resisted when it adds behavior without lifecycle meaning, creates an HTTP API without a schema source, bypasses a named trust boundary, hides skipped coverage, globally weakens a safety gate, or changes persisted data outside an intentional canonicalization path.
