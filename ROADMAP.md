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

**UI**

- [UI] Connection Wizard

### [v0.6.0](../../milestone/3)

**UI**

- [UI] Kubernetes Resource Dashboard

**CLI**

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

- [Provider] MeshMap (beta)
- [MeshSync] - Resync
  
**Lifecycle Management**

- [WASMFilters] - Basic support (CRUD)

**Performance Management**

**Server**

- [Server] Designs Basic Support (CRUD)
- [Server] GitOps - GitHub Actions for Meshery (performance and conformance)
- [Server] Initial Models, Components, Relatioships

**UI**

- [UI] Notification Center

### [v0.7.0](../../milestone/4)

- [Docs] Catch up with the code

**Server**

- [Server] workflow engine, policy (crud)
- [Server] Environments and multi-cluster k8s support
- [Server] Initial GCP and AWS support
- [Server] messaging framework and notification center
- [Server] Policy Engine

**UI**

- [UI] Registration Wizard

**Catalog**

- [Catalog] Basic Support (CRUD)

### [v0.8.0](../../milestone/5)

- [Server] - Environments and Workspaces
- [Adapter] - Distributed Performance Testing
- [MeshSync] - Configurable and Tiered Discovery
- [MeshSync] - Composite Fingerprints

- [Server] Component Generator: Direct Chart or Manifest (Operatorhub)
- [Server] Component Regenerator

**UI**

- [UI] - Finalize State Management
- [UI] - Extensible Authorization
- [UI] - Operations Center (for Workflows)
- [Server/UI] Multiple Telemetry Providers and Custom Boards ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_5_21))
- [Server/UI] User-defined Dashboards and Metrics ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_119))
- [UI] Support for Material UI v5
- [CLI] `system report` - diagnostics reporting
- [CLI] gRPC (streaming of status and events)
- [CLI] Colorizing output

- [CLI] Mesheryctl Code coverage goal: 50%

### [v0.9.0](../../milestone/6)

- [CLI] Mesheryctl Code coverage goal: 60%
- [Server] - Extensible Policies
- [Server] - Configuration Insights and Recommendations
- [Provider] - GitOps Snapshots

### [v1.0.0](../../milestone/7)

**CLI**
- Multi-cluster meshconfig support
- [System] hardening, release process, integration tests, user acceptance testing
- [Server] Direct GCP and AWS support

### [v1.1.0](../../milestone/8)

- [Adapter] Distributed performance management
- [Adapter] Adaptive load optimizers

Refer to [Meshery Roadmap](https://docs.google.com/document/d/1kvcz8jdvFwXmYBBaY2-3fHHUUoy1GJLpZZXuoxZQoOk/edit#) document for detailed info.

**CLI**

- 
- Support for environment and workspaces
- Refactoring `system config` for AKS, EKS
- Mesheryctl Code coverage goal: 40%



Refer to [Meshery Roadmap](https://docs.google.com/document/d/1kvcz8jdvFwXmYBBaY2-3fHHUUoy1GJLpZZXuoxZQoOk/edit#) document for detailed info.
