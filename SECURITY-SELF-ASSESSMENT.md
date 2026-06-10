# Meshery Security Self-Assessment

This self-assessment follows the [CNCF TAG-Security self-assessment
template](https://tag-security.cncf.io/community/assessments/guide/self-assessment/).
It is an internal analysis authored by the Meshery maintainers to describe the
project's security posture, document existing security practices, and identify
areas for improvement. It is **not** an independent audit or attestation.

> **Last updated:** 2026-06-10 · **Template:** CNCF TAG-Security Self-Assessment

---

## Table of contents

- [Metadata](#metadata)
- [Security links](#security-links)
- [Overview](#overview)
  - [Background](#background)
  - [Actors](#actors)
  - [Actions](#actions)
  - [Goals](#goals)
  - [Non-goals](#non-goals)
- [Self-assessment use](#self-assessment-use)
- [Security functions and features](#security-functions-and-features)
- [Project compliance](#project-compliance)
- [Secure development practices](#secure-development-practices)
  - [Development pipeline](#development-pipeline)
  - [Ecosystem](#ecosystem)
  - [Communication channels](#communication-channels)
- [Security issue resolution](#security-issue-resolution)
  - [Responsible disclosure process](#responsible-disclosure-process)
  - [Vulnerability response process](#vulnerability-response-process)
  - [Incident response](#incident-response)
- [Appendix](#appendix)
  - [Known issues over time](#known-issues-over-time)
  - [OpenSSF Best Practices](#openssf-best-practices)
  - [Case studies](#case-studies)
  - [Related projects and vendors](#related-projects-and-vendors)

---

## Metadata

|                       |                                                                                                                                                                                                                                                  |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Assessment stage**  | Complete — 2026-06-10                                                                                                                                                                                                             |
| **Software**          | [github.com/meshery/meshery](https://github.com/meshery/meshery). Supporting repositories: [meshery/meshkit](https://github.com/meshery/meshkit), [meshery/schemas](https://github.com/meshery/schemas), [meshery/meshery-operator](https://github.com/meshery/meshery-operator), [meshery/meshsync](https://github.com/meshery/meshsync) |
| **Security provider** | No. Meshery's primary function is the design and operation (management) of cloud native infrastructure; it is not primarily a security product. It can, however, configure data-plane security features (e.g., mTLS policies) on the meshes it manages. |
| **Languages**         | Go (server, `meshkit`, `schemas`, `mesheryctl`, adapters, operator), TypeScript/JavaScript (React + Next.js UI and provider UI), Rego (OPA relationship/policy evaluation), plus Helm/YAML for packaging                                          |
| **SBOM**              | A CycloneDX Software Bill of Materials is generated in CI ([`bom.yaml` workflow](https://github.com/meshery/meshery/blob/master/.github/workflows/bom.yaml)). Dependency manifests: [`go.mod`](https://github.com/meshery/meshery/blob/master/go.mod), [`ui/package.json`](https://github.com/meshery/meshery/blob/master/ui/package.json). The [`@meshery/schemas`](https://www.npmjs.com/package/@meshery/schemas) npm package is published with build provenance (`npm publish --provenance` via OIDC trusted publishing). A signed, release-attached SBOM is a tracked roadmap item (see [Security functions and features](#security-functions-and-features)). |

---

## Security links

| Document                                  | URL                                                                                                                  |
| ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Security policy                           | <https://github.com/meshery/meshery/blob/master/SECURITY.md>                                                         |
| Security vulnerabilities / disclosure     | <https://docs.meshery.io/project/security-vulnerabilities>                                                           |
| Security Insights manifest                | <https://github.com/meshery/meshery/blob/master/SECURITY-INSIGHTS.yml>                                               |
| OpenSSF Best Practices badge (project 3564) | <https://www.bestpractices.dev/projects/3564>                                                                        |
| Securing a Meshery installation (best practices) | <https://discuss.layer5.io/t/securing-meshery-installation-the-best-practices/679>                            |
| Architecture documentation                | <https://docs.meshery.io/concepts/architecture>                                                                      |
| Release cadence & supported versions      | <https://docs.meshery.io/project/contributing/build-and-release#release-cadence>                                     |

---

## Overview

Meshery is an extensible, self-service engineering platform that facilitates the
management of the full lifecycle of infrastructure. It enables the visual and
collaborative **design** and **operation** of Kubernetes clusters, service
meshes, and multi-cloud infrastructure, and is **highly extensible by design**:
much of its functionality is delivered through pluggable **extension points** —
adapters, providers, and models/integrations — rather than a fixed feature set.
Meshery is a vendor-neutral project that unifies the configuration, lifecycle
management, performance benchmarking, and visual collaboration of cloud native
infrastructure behind a single management plane.

That breadth is itself a security-relevant property. The platform's trust
boundaries and attack surface are defined less by a fixed feature set than by
*which* capabilities are enabled in a given deployment — which extension points
(adapters, remote providers, models/integrations) are active, and which
multi-tenant collaboration features (organizations, teams, workspaces) are in
use. This assessment scopes the **core platform** (`meshery/meshery` and its
supporting repositories) while explicitly acknowledging that extensions can
broaden that surface; see [Actors](#actors), [Non-goals](#non-goals), and
[Ecosystem](#ecosystem).

### Background

Meshery is a [CNCF Sandbox project, accepted on June 22,
2021](https://www.cncf.io/projects/meshery/) (the public "CNCF adopts Meshery"
announcement followed in October 2021). It originated as a service-mesh
management plane and has since broadened into a general management plane that
facilitates engineering teams with the collaboration, learning, planning, and
operational management of their infrastructure — spanning Kubernetes and 10+
service meshes and, increasingly, a broad set of other cloud native and non-mesh
infrastructure.

Much of that breadth is delivered through Meshery's **extension model**. Beyond
its in-tree adapters and core components, functionality is extended through
**adapters** (per-technology controllers), **providers** (pluggable identity,
persistence, and infrastructure/API back-ends), and **models/integrations** (the
registry of cloud native constructs Meshery can design and operate). Community
extensions are developed and maintained in a dedicated
[`meshery-extensions`](https://github.com/meshery-extensions) organization,
separate from the core [`meshery`](https://github.com/meshery) organization, so
that contributors can build and share providers and integrations without
requiring endorsement from the core maintainers. CNCF has described Meshery as
"a highly extensible, self-service management platform" able to manage "any
infrastructure via Providers, Models, Adapters, and its other extension points"
([Scaling Organizational Structure with Meshery's Expanding Ecosystem](https://www.cncf.io/blog/2026/03/04/scaling-organizational-structure-with-mesherys-expanding-ecosystem/),
CNCF, 2026). For this assessment, that distinction matters: the security
boundary separates the **core platform** (covered here) from **out-of-tree
extensions**, whose security posture is the responsibility of their respective
authors/maintainers (see [Non-goals](#non-goals)).

Meshery supports both greenfield and brownfield management of infrastructure.
Its discovery subsystem, **MeshSync**, continuously synchronizes Meshery's view
with the live state of managed clusters, enabling Meshery to operate against
existing ("brownfield") environments rather than only resources it created.
Users model infrastructure as **Meshery Designs** built from a registry of
**Meshery Models** (components and relationships), then deploy, operate, and
benchmark that infrastructure.

The Meshery datastore is intentionally treated as a **cache** of authoritative
cluster state rather than a system of record, which shapes several of the
security trade-offs described below.

### Actors

The following components ("actors") interact to deliver Meshery's functionality.
Isolation between them rests primarily on container/process boundaries and
Kubernetes RBAC; see [Non-goals](#non-goals) and [Security functions and
features](#security-functions-and-features) for the boundaries Meshery does and
does not enforce by default.

| Actor                       | Function                                                                                         | Security / isolation measures                                                                                                                                                                                                                                            |
| --------------------------- | ------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Meshery Server**          | Core control plane: serves the UI, REST API, and a GraphQL API; orchestrates adapters; consumes MeshSync data; manages provider sessions. | Every API route passes through provider → auth → session-injector middleware. With a remote provider, requests require a verified JWT session; with the local provider, the session is anonymous (single-user). Serves plaintext HTTP in-process — **TLS is expected to be terminated by an ingress/reverse proxy.** Runs with broad (cluster-admin-equivalent) Kubernetes RBAC by default. |
| **Meshery Adapters**        | Per-mesh controllers (Istio, Linkerd, Consul, Kuma, App Mesh, NGINX SM, Cilium, etc.) exposing mesh lifecycle/config/performance operations over gRPC. | **Disabled by default** (opt-in via `ADAPTER_URLS`). The server↔adapter gRPC channel uses **insecure (plaintext) transport** today — there is no mutual TLS between server and adapters. Adapters reuse the `meshery-server` ServiceAccount by default.               |
| **MeshSync**                | Kubernetes controller that discovers and continuously syncs cluster resource state to the server via the broker. | Discovery scope is governed by a configurable allow/deny **watch-list** on the `meshsyncs.meshery.io` CRD (defaults include Secrets, ConfigMaps, ClusterRoles, Pods, etc.). Holds no datastore of its own; publishes to NATS. Runs under an operator-provisioned ServiceAccount (cluster-admin-equivalent by default). |
| **Meshery Operator**        | Kubernetes operator that manages the lifecycle of MeshSync and the Meshery Broker via two namespaced CRDs (`brokers.meshery.io`, `meshsyncs.meshery.io`). | Defines a least-privilege `operator-role` scoped to its own `meshery.io` CRDs and a `kube-rbac-proxy` (TokenReview/SubjectAccessReview) that gates the metrics endpoint. Also defines a broad `controller-role` (`*/*`) used to reconcile managed resources. |
| **Meshery Broker (NATS)**   | Message bus that streams data between in-cluster components (MeshSync) and the server, which may run outside the cluster. | Supports username/password authentication and TLS, but both are **unset by default**. Typically exposed via `LoadBalancer`/`NodePort` so an out-of-cluster server can reach it — operators should enable broker auth + TLS when exposed. |
| **Providers** (Local / Remote) | Pluggable identity and persistence backends. The **Local** provider is built-in; **Remote** providers (e.g., Meshery Cloud / Layer5 Cloud) add cloud-backed identity and persistence. | **Local = anonymous, single-user, no login** (by design). **Remote = browser-based OAuth login + RS256 JWT verification** against the provider's JWKS, with token expiry checks, server-side introspection (revocation), and refresh, over HTTPS. The `PROVIDER` setting selects/enforces which provider is in effect. |
| **Registry / Database**     | SQLite datastore holding the model/component registry, designs, MeshSync-discovered resources, connections, credentials, events, and user preferences. | SQLite file (WAL mode) under the user-data folder, **treated as a cache**. No application-layer encryption at rest — protection relies on disk/file permissions (or the remote provider's controls for cloud-persisted data). |
| **mesheryctl**              | Command-line interface for installing, operating, and interacting with Meshery.                  | Authenticates via the same browser-based OAuth flow; the resulting token is written to `~/.meshery/auth.json` on the local filesystem (protect via file permissions).                                                                                                   |
| **UI / provider-ui**        | Next.js single-page applications served by the server for management and provider selection.     | Communicate with the server over REST and GraphQL (subscriptions over WebSocket). The GraphQL WebSocket origin check is currently permissive; restrict origins at the proxy when exposing Meshery.                                                                       |

### Actions

The following are representative actions Meshery performs, emphasizing the
security checks, sensitive data, and inter-actor interactions involved.

1. **User authentication & session establishment.** When configured with a
   remote provider, an unauthenticated request is redirected to the provider's
   browser login flow. On callback, the server receives a token and sets it as
   an `HttpOnly` cookie. On each subsequent request, `AuthMiddleware` calls
   `provider.GetSession()`, which (a) verifies the JWT's RS256 signature against
   the provider's published JWKS, (b) checks expiry, (c) performs server-side
   **introspection** against the provider to detect revocation, and (d)
   **refreshes** the token when needed. With the **local** provider this step is
   a no-op: the session is always treated as a single anonymous user and no
   credential is checked.

2. **Authorization for management functions.** Fine-grained permissions are
   modeled as capability "keys" grouped into keychains and mapped to roles
   (User, Team Admin, Workspace Admin, Org Admin, Provider Admin, etc.) in
   [`server/permissions/keys.csv`](https://github.com/meshery/meshery/blob/master/server/permissions/keys.csv).
   These are seeded into the datastore at startup. Multi-user role enforcement is
   performed in conjunction with a remote provider; the local provider is
   single-user and does not enforce role separation.

3. **Connecting a Kubernetes cluster & discovery.** A user registers a cluster
   by supplying a kubeconfig (stored as a `Credential`/`K8sContext`). Meshery
   deploys the Meshery Operator, which provisions MeshSync and the Broker.
   MeshSync watches the resource types permitted by its watch-list and publishes
   ADDED/MODIFIED/DELETED events to the NATS subject
   `meshery-server.meshsync.store`; the server consumes these and persists them
   as `KubernetesResource` records. Sensitive inputs here include cluster
   credentials and any discovered Secrets.

4. **Designing and deploying infrastructure.** A user composes a Meshery Design
   from registry components and requests deployment. The server validates the
   design against model/relationship policies (Rego/OPA evaluation) and then
   applies the resulting manifests to the target cluster using its Kubernetes
   client. Because the server holds broad cluster RBAC, it can create/update/
   delete arbitrary resources on the user's behalf.

5. **Service-mesh operations via adapters.** Mesh-specific actions (install a
   mesh, apply a configuration, run a performance benchmark, apply an mTLS
   policy on the managed mesh) are dispatched from the server to the relevant
   adapter over gRPC. The adapter performs the operation against the cluster/mesh
   and streams results back.

6. **Performance benchmarking.** Meshery runs load tests (the canonical
   implementation of the CNCF **Service Mesh Performance** specification),
   recording results for comparison and sharing. Benchmark results may be
   persisted locally or to the configured remote provider.

### Goals

**Product goals**

1. Provide a unified, vendor-neutral management plane for Kubernetes and 10+
   service meshes.
2. Enable visual, collaborative **infrastructure-as-design** (Meshery Designs
   and Models) shared across teams.
3. Provide standardized **performance benchmarking** via the canonical Service
   Mesh Performance (SMP) implementation.
4. Facilitate configuration of **data-plane security features** (such as mTLS
   policies) on the service meshes Meshery manages.
5. Support GitOps workflows and configuration **drift detection** against the
   live state discovered by MeshSync.
6. Serve as an **extensible, self-service** platform whose capabilities —
   adapters, providers, and models/integrations — can be extended by the
   community to manage a broad set of cloud native (and non-mesh) infrastructure.

**Security goals (guarantees the project aims to provide)**

- **Authenticated, authorized access** to management functions when deployed
  with a remote provider — including signature-verified, revocation-checked JWT
  sessions and a capability/role-based authorization model.
- **Software supply-chain integrity** — every change is sign-off (DCO) gated,
  peer-reviewed, and must pass a battery of CI security checks before it can be
  merged or released.
- **Bounded data collection** — MeshSync only collects the resource types
  permitted by its watch-list, and the datastore is treated as a cache of
  authoritative cluster state.
- **Coordinated, timely vulnerability response** — a documented private
  disclosure channel and response process (see [Security issue
  resolution](#security-issue-resolution)).

### Non-goals

The following are explicitly **out of scope** for the Meshery open-source
project. Several correct earlier informal descriptions of Meshery's posture.

- **Meshery does not provide mutual TLS or transport encryption between its own
  internal components by default.** The server↔adapter gRPC channel is plaintext
  today, and the server↔broker (NATS) channel ships without auth or TLS by
  default. Securing these channels is a deployment responsibility and an active
  hardening roadmap item; it is not an out-of-the-box "zero-trust" guarantee.
- **The Local provider intentionally performs no authentication** (single-user
  mode). It is meant for local/individual use and is **not** intended for shared,
  multi-tenant, or internet-exposed deployments — use an authenticated remote
  provider for those.
- **Meshery is not a runtime security tool.** It does not provide intrusion
  detection/prevention, admission control, runtime threat detection, or policy
  enforcement on the workloads it manages (it complements, rather than replaces,
  tools such as admission controllers and network policy).
- **Meshery is not a secrets manager / vault.** It does not encrypt credentials
  at rest at the application layer; the datastore must be protected via
  disk/file/cloud controls, and a dedicated secrets manager should hold
  long-lived secrets.
- **Meshery does not replace Kubernetes' own controls.** It relies on the
  cluster's RBAC, NetworkPolicy, and Pod security mechanisms; Meshery ships no
  NetworkPolicies and applies no container hardening (`securityContext`) by
  default, leaving those to the operator.
- **Meshery does not constrain resource consumption of authorized callers.** An
  authorized user can issue operations (large discoveries, benchmarks,
  deployments) that consume cluster or server resources.
- **The security of out-of-tree and community extensions is out of scope for the
  core project.** Meshery is deliberately extensible via adapters, providers, and
  models/integrations, and community extensions are maintained in a separate
  [`meshery-extensions`](https://github.com/meshery-extensions) organization.
  Enabling a third-party adapter, remote provider, or integration extends a
  deployment's trust boundary to that extension; vetting and securing it is the
  responsibility of its author/operator. This assessment covers the core
  `meshery/meshery` platform and its supporting repositories, not arbitrary
  third-party extensions.
- **Formal compliance attestations are out of scope for the OSS artifact.** Any
  SOC 2 / ISO 27001 / PCI-DSS / GDPR obligations pertain to an organization that
  *operates* Meshery as a hosted service (e.g., Meshery Cloud), not to the
  open-source software itself.

---

## Self-assessment use

This self-assessment is created by the Meshery maintainers to perform an
internal analysis of the project's security. It is **not** intended to provide a
security audit of Meshery, nor to function as an independent assessment or
attestation of Meshery's security health.

This document serves to provide Meshery users with an initial understanding of
Meshery's security posture, where to find existing security documentation,
Meshery's plans for security, and a general overview of Meshery security
practices — covering both the secure development of Meshery and the security of
operating Meshery.

This document also provides the CNCF TAG-Security community with an initial
understanding of Meshery to assist in any future joint assessment. Meshery is
currently a **CNCF Sandbox** project; taken together with a future joint
assessment, this document is intended to serve as a cornerstone for if and when
Meshery pursues incubation and, subsequently, a security audit.

---

## Security functions and features

**Critical** — design elements that make the product itself secure. These are
recommended primary inputs to threat modeling and should be tracked as
high-impact items for changes.

- **Provider authentication gate.** All API access flows through
  `AuthMiddleware → provider.GetSession()`. For remote providers this performs
  RS256 JWT signature verification against the provider's JWKS, expiry
  validation, server-side introspection (to honor revocation), and token
  refresh. This is the central authentication control for multi-user
  deployments.
- **Capability/role-based authorization model.** A fine-grained permission
  matrix (`server/permissions/keys.csv`) maps capabilities to roles and is
  enforced together with a remote provider, constraining which management
  actions a principal may perform.
- **Software supply-chain integrity controls.** Mandatory DCO sign-off, required
  peer review with branch protection, and a CI gate (lint/static analysis,
  CodeQL, container and dependency scanning, OpenSSF Scorecard, OSPS Baseline)
  that must pass before merge or release. See [Development
  pipeline](#development-pipeline).
- **Centralized input sanitization for query ordering.** A shared
  `SanitizeOrderInput` routine was introduced to close the class of SQL-injection
  issues found in `order`/`sort` query parameters (see [Known issues over
  time](#known-issues-over-time)).

**Security relevant** — configurable components that materially affect the
security posture of a deployment. These should also be included in threat
modeling.

- **Provider selection (`PROVIDER`).** The single most important lever: the
  Local provider disables authentication (single-user); a remote provider must
  be selected/enforced for authenticated, multi-user, or exposed deployments.
- **TLS termination.** The server speaks plaintext HTTP in-process; production
  deployments must terminate TLS at an ingress/reverse proxy.
- **Session cookie attributes.** Session/token cookies are `HttpOnly` but do not
  currently set `Secure` or `SameSite`; deploy behind TLS and consider proxy-level
  hardening. (Hardening these flags is a roadmap item.)
- **MeshSync watch-list.** Controls which resource types (including whether
  Secrets) are discovered and streamed; tighten to the minimum needed.
- **Broker (NATS) authentication and TLS.** Off by default; enable both whenever
  the broker is exposed outside the cluster.
- **Kubernetes RBAC scope.** The server and operator (and, by default, the
  adapters, MeshSync, and broker via a shared ServiceAccount) run with
  cluster-admin-equivalent RBAC to support discovery and lifecycle management.
  Operators should scope this down where the use case permits; the `rbac.nodes`
  toggle controls optional node access.
- **Network isolation.** No NetworkPolicies are shipped; the server `Service`
  and broker default to `LoadBalancer` exposure. Add NetworkPolicies and limit
  exposure in production.
- **Pod/container hardening.** `securityContext`/`podSecurityContext` are empty
  by default (no `runAsNonRoot`, `readOnlyRootFilesystem`, or capability drops);
  apply hardening in production.
- **Operator opt-out (`DISABLE_OPERATOR`).** Allows running Meshery without the
  in-cluster discovery stack (operator/MeshSync/broker) when discovery is not
  required.
- **GraphQL WebSocket origin policy.** The subscription endpoint's origin check
  is currently permissive; restrict allowed origins at the proxy to mitigate
  cross-site WebSocket hijacking when Meshery is exposed.
- **Extension surface (adapters, providers, models/integrations).** Because
  Meshery is extensible by design, the set of enabled extension points is itself
  a security-relevant lever: each enabled adapter, remote provider, or
  third-party integration becomes part of the deployment's trust boundary and
  attack surface. Prefer first-party/maintained extensions, review third-party
  extensions before enabling them, and enable only those a deployment actually
  needs.

**Roadmap / hardening items** identified during this assessment (tracked for
future work): in-process TLS option; mutual TLS for the server↔adapter channel;
broker auth + TLS enabled by default; tightened default RBAC and per-component
ServiceAccounts; `Secure`/`SameSite` cookie flags; optional at-rest encryption
for stored credentials; container `securityContext` defaults; signed container
images (Cosign/Sigstore) and SLSA provenance with a release-attached SBOM; and a
`SECURITY.md` pointer in the `meshkit` repository.

---

## Project compliance

- **OpenSSF Best Practices badge — passing.** Meshery holds a *passing*
  [OpenSSF Best Practices badge (project
  3564)](https://www.bestpractices.dev/projects/3564), first achieved
  2020-01-08. Beyond the *passing* tier it has satisfied roughly 45% of the
  additional Silver-level criteria (a 145% tiered score; Silver corresponds to
  200%).
- **License.** Apache License 2.0.
- **Developer Certificate of Origin (DCO).** All commits must be signed off per
  the [DCO](https://developercertificate.org/), enforced in CI.
- **Formal certifications (SOC 2, ISO 27001, PCI-DSS, FIPS, GDPR).** The Meshery
  open-source project does **not** claim any formal compliance certification.
  Such attestations would apply to an organization that *operates* Meshery as a
  hosted service (for example, Meshery Cloud), not to the open-source artifact
  itself.

---

## Secure development practices

### Development pipeline

See the [contribution flow](https://docs.meshery.io/project/contributing#general-contribution-flow)
for the end-to-end process. Source is hosted on GitHub across `meshery/meshery`,
`meshery/meshkit`, and `meshery/schemas` (all Go 1.26.x), plus the
`meshery/meshery-operator` and `meshery/meshsync` repositories.

**Automated checks (must pass before review/merge).** When a pull request is
opened, a series of CI workflows run:

- **Linting & static analysis** — `golangci-lint` v2 (with `staticcheck`,
  `govet`, `forbidigo`, `unused`, etc.) across the `server` and `mesheryctl`
  modules; `ESLint` + `Prettier` for the UIs; `Rego` lint and tests for policies.
- **SAST** — **CodeQL** for Go and JavaScript/TypeScript (scheduled and
  on-dispatch).
- **Container & dependency scanning** — **Trivy** scans the built image for
  CRITICAL/HIGH findings (auto-filing a tracking issue); a CycloneDX SBOM is
  generated and scanned with **Grype** (with **Meterian** as an additional
  scanner).
- **Supply-chain posture** — **OpenSSF Scorecard** (results published to code
  scanning) and an **OpenSSF OSPS Baseline** scan.
- **Tests** — Go unit and integration tests run with the race detector
  (`go test -race`), including `mesheryctl` integration tests against Docker and
  Kubernetes (Kind); UI unit tests with **Vitest** + React Testing Library; and
  end-to-end tests with **Playwright**. The E2E job runs under **StepSecurity
  Harden-Runner** (egress audit).
- **Packaging** — multi-arch Docker image builds (`linux/amd64,linux/arm64`)
  and `helm lint` for the `meshery` and `meshery-operator` charts.
- **DCO** — all commits must be signed off per the [DCO
  guidelines](https://developercertificate.org/).

> **Accuracy note (correcting prior descriptions):** Meshery's Go SAST/lint is
> provided by **`golangci-lint` v2 + CodeQL**, not `gosec` (which is not used in
> any Meshery repository). Container/dependency scanning is **Trivy + Grype +
> Meterian**. End-to-end testing uses **Playwright** (currently non-blocking —
> failing E2E does not yet gate merges).

**Manual code review.** After CI passes, designated maintainers review the PR
for adherence to coding conventions and quality standards, absence of security
issues, and correctness of any documentation/Helm/artifact updates. At least one
maintainer approval is required to merge; contributions typically receive
feedback from multiple maintainers.

**Branch protection and merging.** Repositories enforce branch protection: a PR
may merge only when required checks pass, the required approvals are obtained,
and the branch is up to date with `master`.

**Artifact generation on merge/release.** Releases tag and publish multi-arch
Docker images to Docker Hub (stable/edge channels), build and publish the
`mesheryctl` CLI via **GoReleaser** (with a `checksums.txt` of SHA-256 digests,
plus Homebrew and Scoop distribution), and publish Helm charts (advertised via
Artifact Hub). The `@meshery/schemas` npm package is published with **build
provenance** via OIDC trusted publishing.

> **Known gaps:** container images and CLI binaries are **not** currently signed
> (no Cosign/Sigstore), and there is no SLSA provenance or release-attached SBOM
> for the server/CLI artifacts — these are tracked hardening items.

**Dependency management.** Dependabot is configured for Go modules and npm in
`meshery/meshery` and `meshery/schemas` (monthly). (`meshkit` does not currently
configure Dependabot — a noted gap.)

**Testing summary**

| Test type            | Tools used                                                            |
| -------------------- | --------------------------------------------------------------------- |
| Unit / integration   | Go `testing` with `-race`; Kind-based integration tests for `mesheryctl` |
| UI unit              | Vitest + React Testing Library                                        |
| End-to-end           | Playwright (currently non-blocking)                                   |
| Static analysis/SAST | golangci-lint v2 (staticcheck/govet/…), CodeQL                        |
| Container / deps     | Trivy, Grype (CycloneDX SBOM), Meterian                               |
| Supply chain         | OpenSSF Scorecard, OpenSSF OSPS Baseline                              |

**Vulnerability disclosure.** See [Security issue
resolution](#security-issue-resolution).

### Ecosystem

Meshery occupies the **cloud native management plane** niche: a vendor-neutral,
**highly extensible, self-service** platform for the collaborative design and
operation of Kubernetes, service mesh, and broader cloud native infrastructure.
Its position in the ecosystem is reinforced by:

- **Standards stewardship.** Meshery is the **canonical/reference
  implementation** of the CNCF [Service Mesh Performance
  (SMP)](https://smp-spec.io/) specification and has served as a Service Mesh
  Interface (SMI) conformance tool (SMI has since been largely succeeded within
  CNCF by the Gateway API GAMMA initiative).
- **Breadth of integrations and extension points.** Meshery integrates with
  Kubernetes and 10+ service meshes (Istio, Linkerd, Consul, Kuma, AWS App Mesh,
  Open Service Mesh, NGINX Service Mesh, Traefik Mesh, Network Service Mesh,
  Cilium, and others) and a large and growing catalog of cloud native
  integrations via the **Meshery Models** registry (380+ integrations; see the
  live [integrations page](https://meshery.io/integrations)). This breadth is
  delivered through Meshery's extension points — **adapters**, **providers**, and
  **models/integrations** — with community-contributed extensions maintained in a
  separate [`meshery-extensions`](https://github.com/meshery-extensions)
  organization so they can evolve independently of the core.
- **Adoption pattern.** Because Meshery is frequently used as the management and
  benchmarking layer over widely-deployed meshes and Kubernetes, a relatively
  small number of Meshery operators can represent broad downstream usage across
  the meshes and clusters they manage.
- **Contributor ecosystem.** By CNCF's accounting, Meshery is "the CNCF's sixth
  highest-velocity project," supported by weekly newcomer meetings, self-paced
  training, and the **Certified Meshery Contributor (CMC)** program — described
  as CNCF's first contributor certification of its kind — whose curriculum spans
  Meshery Server, CLI, UI, Models, and **Extensibility**
  ([Announcing the Certified Meshery Contributor](https://www.cncf.io/blog/2025/10/27/announcing-the-certified-meshery-contributor-cmc/),
  CNCF, 2025). From a security standpoint, a large, well-onboarded contributor
  base is a positive supply-chain signal — more reviewers and structured
  review — even as it widens the population interacting with the codebase, which
  is why the development pipeline leans on the automated gates and mandatory
  review described above.

### Communication channels

| Category     | Channel               | Description                                                                                  | Link                                                                 |
| ------------ | --------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| **Internal** | Team meetings         | Regular meetings among maintainers to coordinate project progress.                            | N/A                                                                  |
| **Internal** | Security mailing list | Private channel for triaging reported vulnerabilities.                                        | <security@meshery.dev>                                               |
| **Inbound**  | Slack                 | Public community where users and contributors ask questions and collaborate.                  | [Meshery Slack](https://slack.meshery.io)                             |
| **Inbound**  | Discussion forum      | Asynchronous, in-depth discussion and support.                                                | [Meshery Discussion Forum](https://discuss.meshery.io/)               |
| **Inbound**  | Community meetings     | Weekly meetings open to all; recordings on YouTube.                                           | [Meshery Community](https://meshery.io/community)                    |
| **Outbound** | Social media          | Updates and announcements.                                                                    | [X/Twitter](https://twitter.com/mesheryio) · [LinkedIn](https://www.linkedin.com/showcase/meshery/) |
| **Outbound** | YouTube               | Tutorials, demos, and meeting recordings.                                                     | [Meshery YouTube](https://www.youtube.com/@Layer5io)                 |

---

## Security issue resolution

### Responsible disclosure process

- **Reporting.** Suspected vulnerabilities should be reported privately to the
  Meshery security mailing list, **<security@meshery.dev>**, as documented at
  <https://docs.meshery.io/project/security-vulnerabilities>. Reports should
  include the nature and location of the issue, its potential impact, and steps
  to reproduce.
- **Acknowledgement & evaluation.** The Meshery team acknowledges and begins
  analyzing each report within **10 working days** and assesses severity. Report
  details remain confidential and are not shared externally except as necessary
  to resolve the issue.

### Vulnerability response process

- **Responsibility.** The Meshery maintainers triage and resolve reported
  vulnerabilities.
- **Process.**
  1. **Fix development** — fixes are developed and tested privately to prevent
     premature disclosure.
  2. **Coordinated disclosure** — once a fix is ready, patched binaries/images are
     published and the fix is announced via the [Meshery blog](https://meshery.io/blog),
     social channels, and the `#announcements` channel on Slack. Where applicable,
     a GitHub Security Advisory / CVE is published. Supported versions and release
     cadence are documented at
     <https://docs.meshery.io/project/contributing/build-and-release#release-cadence>.

### Incident response

- **Triage & containment.** On identifying an incident, immediate containment
  measures are taken (e.g., disabling affected components or applying temporary
  mitigations).
- **Root-cause analysis.** Logs and telemetry are investigated to determine the
  root cause.
- **Stakeholder notification.** Confirmed incidents are communicated to users via
  Slack, email, and other channels, with actionable mitigation guidance.
- **Patch deployment.** Fixes are tested before release; updates are communicated
  broadly so users can apply patches promptly.

---

## Appendix

### Known issues over time

All publicly disclosed Meshery vulnerabilities to date fall into a single class —
**SQL injection (CWE-89)** arising from unsanitized `order`/`sort` query
parameters concatenated into ORM queries — and **all have been remediated**.
Remediation converged on a central `SanitizeOrderInput` routine. The supporting
`meshkit` and `schemas` repositories have **no** published advisories.

| CVE             | GHSA                  | Severity | CVSS v3 | Affected versions | Fixed in  | Published  |
| --------------- | --------------------- | -------- | ------- | ----------------- | --------- | ---------- |
| CVE-2021-31856  | —                     | Critical | 9.8     | 0.5.2             | 0.5.3     | 2021-04-28 |
| CVE-2023-46575  | GHSA-9jjc-grg5-67gj   | Critical | 9.1     | < 0.6.179         | 0.6.179   | 2023-11-24 |
| CVE-2024-29031  | GHSA-652r-q29p-m25h   | High     | 7.5     | < 0.7.17          | 0.7.17    | 2024-08-05 |
| CVE-2024-35181  | GHSA-9f24-jrv4-f8g5   | Medium   | 5.9     | < 0.7.22          | 0.7.22    | 2024-08-05 |
| CVE-2024-35182  | GHSA-h7cm-jvpp-69xf   | Medium   | 5.9     | < 0.7.22          | 0.7.22    | 2024-08-05 |

Notes:

- The 2024 advisories originated from GitHub Security Lab reports (GHSL-2024-013/014)
  and were fixed via PRs #9372, #10207, and #10280.
- The advisories are recorded under the historical module path
  `github.com/layer5io/meshery`; the project's current module path is
  `github.com/meshery/meshery`.
- **As of this assessment there are no known open/unpatched advisories.** The CI
  scanning described under [Development pipeline](#development-pipeline) (CodeQL,
  Trivy, Grype, Scorecard) is the project's track record for catching issues in
  review and automated testing.

### OpenSSF Best Practices

Meshery holds a **passing** [OpenSSF Best Practices badge (project
3564)](https://www.bestpractices.dev/projects/3564), first earned 2020-01-08,
and, beyond the *passing* tier, has satisfied roughly **45%** of the additional
Silver-level criteria (a 145% tiered score; Silver corresponds to 200%).
Achieving the remaining Silver criteria (and pursuing Gold thereafter) is the
project's badge roadmap. Meshery also runs the **OpenSSF Scorecard** action in CI
with results published to GitHub code scanning.

### Case studies

The following representative scenarios illustrate real-world Meshery usage. They
are framed as capability-grounded use cases.

1. **Multi-mesh operations for a platform team.** A platform engineering team
   standardizing on more than one service mesh (e.g., Istio in one environment,
   Linkerd in another) uses Meshery as a single management plane to install,
   configure, and visualize each mesh, and to apply consistent configuration
   (including mesh mTLS policies) across them — avoiding a separate, bespoke
   toolchain per mesh.

2. **Performance benchmarking and regression gating.** A team evaluating the
   latency/throughput overhead of a service mesh (or comparing mesh versions and
   configurations) uses Meshery's canonical SMP implementation to run repeatable
   load tests, capture standardized results, and compare them over time as part
   of release decisions.

3. **Collaborative, drift-aware design.** Teams model infrastructure as Meshery
   Designs from the Models registry, review them collaboratively, deploy them,
   and then rely on MeshSync's continuous discovery to detect drift between the
   intended design and the live cluster state.

### Related projects and vendors

Prospective users frequently compare Meshery with other Kubernetes management
and visualization tools. Brief, neutral differentiators:

- **Rancher** — server-side enterprise multi-cluster fleet management
  (provisioning, fleet ops, RBAC). Heavier operations platform; not focused on
  visual design or mesh performance benchmarking.
- **Lens / Headlamp** — Kubernetes dashboards/IDEs (desktop and, for Headlamp,
  web) focused on cluster navigation and extensibility. Meshery is an in-cluster/
  web management plane oriented toward multi-mesh management and design.
- **Backstage** — developer portal / service catalog for platform engineering;
  broader than Kubernetes and oriented to software cataloging rather than mesh/
  infrastructure operation.
- **Argo CD** — GitOps continuous delivery; overlaps on declarative desired-state
  operation but is CD-specific, not multi-mesh design/benchmarking.
- **Kiali** — Istio service-mesh observability dashboard; single-mesh
  observability versus Meshery's multi-mesh management plus SMP performance
  benchmarking.

Meshery's combination of visual infrastructure-as-design (Designs/Models),
multi-service-mesh support, and SMP performance benchmarking is the differentiator
that none of the comparable tools fully replicate.
