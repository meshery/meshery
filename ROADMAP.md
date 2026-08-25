# Meshery Roadmap

Milestones on Meshery's high-level roadmap:

### [v1.3.0](../../milestone/10)

**Configuration Management**

- [Server] Expand Policy Engine and support Policy (crud)
- [Server] Generative AI Configuration Analysis
- [Server] Notification Providers: Slack, CloudEvents, Teams

**Server**

- [Server] Configuration Insights and Recommendations
- [Server] Correlated Events

### [v1.2.0](../../milestone/9)

**Infrastructure Knowledge & Registry**

- [Server] Workflow Engine, policy (CRUD)
- [UI] Operations Center (for Workflows)

- [CLI] Meshconfig: Support for multi-cluster (spec)

**Sustaining**

- [Server] Scalability: Postgres support/migration plan (spec)

**Performance Management**

- [Adapter] Distributed Performance Testing
- [Adapter] Adaptive load optimizers
- [Server/UI] Multiple Telemetry Providers and Custom Boards ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g1044af767ce_5_21))
- [Server/UI] User-defined Dashboards and Metrics ([spec](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.gcb74201a11_0_119))

### [v1.1.0](../../milestone/8)

**Infrastructure Knowledge & Registry**

- [Server/UI] Registry governance: team-level model controls, Ignore action hardening  
- [Server] Model sub-categories; Models in Catalog; Catalog tags and filters  
- [Catalog] Payment-processing readiness (see v1.1.0 **Catalog** section below for details)

**Governed Change & Lifecycle Orchestration**

- [Server] Enhanced GitHub connection with repo ingest; design import wizard  

**Multi-Cluster & Fleet Operations**

- [Server] Fine-grained Kubernetes RBAC integration; embedded MeshSync mode without CRD installation (spec)  
- [Server] Realtime MeshSync APIs

**Collaborative and Visual Infrastructure Operations**

- [UI] UI restructure is a cross-cutting effort tracked elsewhere on this roadmap; referenced here because it impacts visual infrastructure operations.  
- [UI/Server] Interactive consoles: terminal and logs over WebSockets, dockable panel (\#20622)  
- [Server/UI] User-defined dashboards and metrics; multiple telemetry providers (spec)

**Agentic Infrastructure Management**

- [Extensions] Meshery MCP Server: read-only access to Registry, designs, and cluster state (meshery-mcp-server)  
- [Server] Semantic retrieval and context assembly over the Registry (\#19645, spec)

**Performance & Reliability Intelligence**

- [Server] Prometheus and Grafana connection management

**Sustaining**

- [CLI] `system report` diagnostics; gRPC streaming of status and events; colorized output; coverage 50 percent  
- \[Server\] API stability and deprecation policy for v1 schemas  
- \[Docs/CI\] E2E coverage 50 percent; automated publication of e2e results

**Configuration Management**

- [Server] Support for OPA and Golang-based Policy Engines

**Extensibility / Extensions**

- [Server] Remote provider: gitops

**CLI**

- [CLI] Multi-cluster meshconfig support
- [CLI] `system report` - diagnostics reporting
- [CLI] Refactoring `system config` for AKS, EKS
- [CLI] gRPC (streaming of status and events)
- [CLI] Colorizing output

**Catalog**

- [Catalog] Intellectual property protections for user-produced content
- [UI/Server] Improved performance and stability in catalog interactions to facilitate payment processing

### [v1.0.0](../../milestone/7)

**Lifecycle Management**

- [MeshSync] Configurable and Tiered Discovery
- [MeshSync] Composite Fingerprints

**General / Maintenance**
- [UI] Restructure is tracked in the canonical UI roadmap entry; reference that section for status and updates.
- [UI] Sistent as sole theme provider; Material UI v5 completion  
- [System] hardening, release process, integration tests, user acceptance testing
- [System] Schema-driven development with meshery/schemas as canonical source of resource definition

**Performance Management**

- [Server/UI] Consolidation of built-in load generators

### [v0.9.0](./docs/_releases/roadmap-v0.9.md)

- [Server] Extensible Policies
- [Server] GitOps: PR as an Action, Expand Flux and ArgoCD Integrations
- [Server] Advanced Environments and Workspaces
- [Server] Initial Azure as Model and ASO Integration
- [Server] Advanced GCP and AWS as Connections

**Extensibility / Extensions**

- [CLI] kubectl Snapshot

_See [(detailed v0.9.0 roadmap)](https://discuss.meshery.io/t/meshery-v0-9-roadmap/6296)_

### [v0.8.0](./docs/_releases/roadmap-v0.8.md)

**Lifecycle Management**

- [Server/UI] Robust Model Generator
- [Server/UI] 25% coverage of Relationships across all Models
- [Server/UI] Model Import/Export, OCI, Extensible
- [Server] Initial GCP and AWS as Models
- [Server] Initial Environments and Workspaces
Code coverage goal: 25%

**Configuration Management**

- [Server] Component Generator: Direct Chart or Manifest (Operatorhub)
- [Server] Registry: Model import/export; OCI
- [CLI] Support for commands: model, component, relationship, registry, environment, connection, credential

**Extensibility / Extensions**

- [UI] Extensible Authorization
- [CLI] Helm Snapshot
- [CLI] kubectl MeshSync Snapshot

**General / Maintenance**

- [UI] Finalize State Management
- [UI] Support for Material UI v5
- [CLI] Deprecate: Full migration from Apps to Designs
- [CLI] Mesheryctl Code coverage goal: 50%

_See [(detailed v0.8.0 roadmap)](https://discuss.meshery.io/t/meshery-v0-8-0-roadmap/4336)_

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

_See [(detailed v0.7.0 roadmap)](https://discuss.meshery.io/t/meshery-v0-7-0-roadmap/232)_

### [v0.6.0](../../milestone/3)

**Lifecycle Management**

- [UI] Kubernetes Resource Dashboard
- [CLI] Create `mesh` (adapter operations)
- [CLI] pervasive Kubernetes support
- [CLI] Initial `patterns` support
- [CLI] Refactoring `perf` to support SMP better
- [CLI] Confirm support for Linux, Windows, and MacOS across all current commands
- [CLI] Refactoring `system config` for GKE
- [CLI] OAuth support with Remote Providers
- [CLI] `system check` - pre and post flight prereq check
- [CLI] Refactoring `system` commands for docker-compose
- [WASMFilters] - Basic support (CRUD)

**Extensibility**

- [Provider] Full Page / Navigation Menu Plugin
- [MeshSync] - Resync

**Configuration Management**

- [Server] Designs Basic Support (CRUD)
- [Server] GitOps - GitHub Actions for Meshery (performance and conformance)
- [Server] Initial Models, Components, Relationships

**Event Management**

- [UI] Notification Center

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
