# Installation

LLMS index: [llms.txt](/llms.txt)

---

Section pages:

- [Deployment Models & Reference Architecture](/installation/production/deployment-models/): In-cluster vs. out-of-cluster, Docker vs. Kubernetes, the Meshery component inventory and its statefulness, and the topology patterns for single-cluster, multi-cluster, and multi-cloud production deployments.
- [Infrastructure, Sizing & Performance](/installation/production/infrastructure-sizing-and-performance/): Resource requirements per Meshery component, capacity planning, MeshSync tiered discovery, Broker throughput, scalability levers, and the known performance bounds to plan around in production.
- [Production Deployment](/installation/production/): Considerations, best practices, security hardening, performance bounds, and operational readiness guidance for deploying Meshery in production.
- [High Availability & Resiliency](/installation/production/high-availability-and-resiliency/): Replication, health probes, failure modes and recovery, the ephemeral database, Remote Provider persistence, and backup & disaster-recovery posture for resilient Meshery operation.
- [Upgrading Meshery](/installation/upgrades/): How Meshery's components version together, and the step-by-step procedure for upgrading the CLI, Server, Operator, MeshSync, and Broker on Docker and Kubernetes deployments.
- [Networking & Connectivity](/installation/production/networking-and-connectivity/): Network ports and directional flows, ingress and Emissary configuration, secure WebSocket support, CDN and edge caching of the UI, Broker exposure, egress to Remote Providers, the OAuth callback URL, and network policies for production Meshery.
- [Security Hardening](/installation/production/security-hardening/): RBAC and least privilege, pod and container security contexts, secret and kubeconfig handling, TLS, supply-chain integrity, Broker exposure risk, namespace isolation, and extension trust for hardening Meshery in production.
- [Authentication, Authorization & Identity](/installation/production/authentication-and-identity/): Why to preselect a Remote Provider over the Local Provider, OAuth callback configuration, identity providers, capabilities, and keys/permissions for production Meshery.
- [Multi-Cluster & Multi-Cloud Operations](/installation/production/multi-cluster-and-multi-cloud/): Managed vs. unmanaged cluster connections, one Operator per cluster, kubeconfig and context management, MeshSync deployment modes, and cloud-specific guidance for operating Meshery across many clusters and clouds.
- [Monitoring, Observability & Health KPIs](/installation/production/monitoring-observability-and-kpis/): Health endpoints, the key performance indicators of Meshery's health, metrics, distributed tracing, centralized logging, alerting, and synthetic checks for observing Meshery in production.
- [Operational Readiness Checklist & Known Caveats](/installation/production/operational-readiness-checklist/): A consolidated, actionable production-readiness checklist across every dimension, the upgrade strategy for Meshery, and the known caveats and limitations to plan around.
- [Using Multiple Adapters](/installation/advanced/multiple-adapters/): Choose which Meshery Adapters to run, run more than one instance of an adapter, or run none at all.
