# Meshery Roadmap

Milestones on Meshery's high-level roadmap:

### [v0.5.0](../../milestone/1)

- [CLI] Release channel subscriptions with system channel
- [CLI] Deployment management with system context
- [Adapter] NGINX Service Mesh
- [Adapter] Traefik Mesh
- Meshery Operator
- [Server] Extensible GraphQL
- [Remote Provider] Dynamic Plugin Injection
- MeshKit and Meshery Adapter Library
- [Server] Performance profiles
- [Server] Patterns (crud)
- [MeshKit] Error Codes Utility

**Lifecycle Management**

- [UI] Connection Wizard

### [v0.6.0](../../milestone/3)

**Lifecycle Management**

- [UI] Kubernetes Resource Dashboard

**Lifecycle Management**

- [CLI] Create `mesh` (adapter operations)
- [CLI] pervasive kubernetes support
- [CLI] Initial `patterns` support
- [CLI] Refactoring `perf` to support SMP better
- [CLI] Confirm support for Linux, Windows, and MacOS across all current commands
- [CLI] Refactoring `system config` for GKE
- [CLI] OAuth support with Remote Providers
- [CLI] `system check` - pre and post flight prereq check
- [CLI] Refactoring `system` commands for docker-compose

**Extensibility**

- [Provider] Full Page / Navigation Menu Plugin
- [MeshSync] - Resync

**Lifecycle Management**

- [WASMFilters] - Basic support (CRUD)

**Performance Management**

**Configuration Management**

- [Server] Designs Basic Support (CRUD)
- [Server] GitOps - GitHub Actions for Meshery (performance and conformance)
- [Server] Initial Models, Components, Relationships

**Event Management**

- [UI] Notification Center

### [v0.7.0](./docs/_releases/roadmap-v0.7.md)


**Sustainability**

- [Docs] Catch up with the code

**Server**

- [Server] Environments and multi-cluster k8s support
- [Server] Initial GCP and AWS support
- [Server] messaging framework and notification center
- [Server] Policy Engine

**Lifecycle Management**

- [UI] Registration Wizard
- [CLI] Support for environment and workspaces

**Catalog**

- [Catalog] Basic Support (CRUD)

**Extensibility**

- [Provider] GitOps Snapshots

### [v0.8.0](./docs/_releases/roadmap-v0.8.md)

**Lifecycle Management**

- [Server/UI] Robust Model Generator
- [Server/UI] 25% coverage of Relationships all Models
- [Server/UI] Model Import/Export, OCI, Extensible
- [Server] Initial GCP and AWS as Models
- [Server] Initial Environments and Workspaces
Code coverage goal: 25%

**Configuration Management**

- [Server] Component Generator: Direct Chart or Manifest (Operatorhub)
- [Server] Registry: Model import/export; OCI
- [CLI] Support for commands: model, component, releationship, registry, environment, connection, credential

**Extensibility / Extensions**

- [UI] Extensible Authorization
- [CLI] Helm Snapshot
- [CLI] kubectl MeshSync Snapshot

**General / Maintenance**

- [UI] Finalize State Management
- [UI] Support for Material UI v5
- [CLI] Deprecate: Full migration from Apps to Designs
- [CLI] Mesheryctl Code coverage goal: 50%

### [v0.9.0](./docs/_releases/roadmap-v0.9.md)

- [Server] Extensible Policies
- [Server] GitOps: PR as an Action, Expand Flux and ArgoCD Integrations
- [Server] Advanced Environments and Workspaces
- [Server] Initial Azure as Model and ASO Integration
- [Server] Advanced GCP and AWS as Connections

**Extensibility / Extensions**

- [CLI] kubectl Snapshot

### [v1.0.0](../../milestone/7)

**Lifecycle Management**

- [MeshSync] Configurable and Tiered Discovery
- [MeshSync] Composite Fingerprints

**General / Maintenance**

- [System] hardening, release process, integration tests, user acceptance testing
- [System] Schema-driven development with meshery/schemas as canonical source of resource definition

**Performance Management**

- [Server/UI] Consolidation of built-in load generators

### [v1.1.0](../../milestone/8)

**UI**
- [UI] [Restructure](https://groups.google.com/a/meshery.io/g/developers/c/-xRKOsrNdg0)

**Configuration Management**

- [Server] Support for OPA and Golang-based Policy Engines

**Extensibility / Extensions**

- [Server] Remote provider: gitops

**Performance Management**

- [Adapter] - Distributed Performance Testing
- Multiple Telemetry Providers and Custom Boards ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_5_21))
- [Server/UI] User-defined Dashboards and Metrics ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_119))

**CLI**
- [CLI] Multi-cluster meshconfig support
- [CLI] `system report` - diagnostics reporting
- [CLI] Refactoring `system config` for AKS, EKS
- [CLI] gRPC (streaming of status and events)
- [CLI] Colorizing output

**Catalog**

- [Catalog] Intellectual property protections for user-produced content
- [UI/Server] Improved performance and stability in catalog interactions to facilitate payment processing

### [v1.2.0](../../milestone/9)

**CLI**

- [CLI] Meshconfig: Support for multi-cluster (spec)

**UI**

- [UI] Operations Center (for Workflows)

**Server**

- [Server] Workflow Engine, policy (crud)
- [Server] SQLite → Postgres (spec)

**Performance Management**

- [Adapter] Distributed performance management
- [Adapter] Adaptive load optimizers

**Extensibility / Extensions**

- [Adapter] Adaptive load optimizer plugin

### [v1.3.0](../../milestone/10)

**Configuration Management**

- [Server] Expand Policy Engine and support Policy (crud)
- [Server] Generative AI Configuration Analysis
- [Server] Notification Providers: Slack, CloudEvents, Teams

**Server**

- [Server] Configuration Insights and Recommendations
- [Server] Correlated Events


Refer to [Meshery Roadmap](https://docs.google.com/document/d/1kvcz8jdvFwXmYBBaY2-3fHHUUoy1GJLpZZXuoxZQoOk/edit#) document for detailed info.
