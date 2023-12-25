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

### [v0.7.0](../../milestone/4)


**Sustainability**

- [Docs] Catch up with the code

**Server**

- [Server] workflow engine, policy (crud)
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

### [v0.8.0](../../milestone/5)

**Lifecycle Management**

- [?] Component generator for direct Chart or Manifest (Operatorhub)
- [?] SQLite → Postgres (spec)
- [?] GitOps: Expand Flux and ArgoCD Integrations
- [?] Relationship Coverage for All Components 
- [?] Model import/export
- [?] GCP and AWS as platforms, Components
Code coverage goal: 25%

- [Server] - Environments and Workspaces
- [MeshSync] - Configurable and Tiered Discovery
- [MeshSync] - Composite Fingerprints
- [CLI] [?] Meshconfig: Support for multi-cluster (spec)

- [CLI] [?] Support for command:
connection, credential
**Performance Management**

- [Adapter] - Distributed Performance Testing
- [Server/UI] Multiple Telemetry Providers and Custom Boards ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_5_21))
- [Server/UI] User-defined Dashboards and Metrics ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_119))

**Configuration Management**

- [Server] Component Generator: Direct Chart or Manifest (Operatorhub)
- [Server] Registry: Model import/export; OCI
- [CLI] [Aadhitya A.] Support for commands: 
- [CLI] model, component, relationship

**Extensibility**

- [UI] - Extensible Authorization

**General / Maintenance**

- [UI] Operations Center (for Workflows)
- [UI] Finalize State Management
- [UI] Support for Material UI v5
- [CLI] `system report` - diagnostics reporting
- [CLI] [?] Refactor `system config` for AKS (spec)
- [CLI] [?] Refactor `system config` for EKS (spec)
- [CLI] [?] Full migration from Apps to Designs
[?] Code coverage goal: 90%
- [CLI] Mesheryctl Code coverage goal: 50%

### [v0.9.0](../../milestone/6)

- [CLI] Mesheryctl Code coverage goal: 60%
- [Server] - Extensible Policies
- [Server] - Configuration Insights and Recommendations

**Lifecycle Management**

- [Server] Multi-cluster meshconfig support

### [v1.0.0](../../milestone/7)

**General / Maintenance**

- [System] hardening, release process, integration tests, user acceptance testing
- [Server] Direct GCP and AWS support

### [v1.1.0](../../milestone/8)

**Configuration Management**

- [Server] Expand Policy Engine and support Policy (crud)
- [Server] Generative AI Configuration Analysis
- [Server] Notification Providers: Slack, CloudEvents, Teams
- [Server] Code coverage goal: 70%

**Extensibility**

- [Server] Remote provider: gitops
- [Adapter] Adaptive load optimizer plugin

**Performance Management**

- [Adapter] Distributed performance management
- [Adapter] Adaptive load optimizers


**CLI**

- [CLI] Refactoring `system config` for AKS, EKS
- [CLI] gRPC (streaming of status and events)
- [CLI] Colorizing output
- [CLI] Mesheryctl Code coverage goal: 40%

Refer to [Meshery Roadmap](https://docs.google.com/document/d/1kvcz8jdvFwXmYBBaY2-3fHHUUoy1GJLpZZXuoxZQoOk/edit#) document for detailed info.
