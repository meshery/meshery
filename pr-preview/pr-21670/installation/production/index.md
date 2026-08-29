# Production Deployment

> Considerations, best practices, security hardening, performance bounds, and operational readiness guidance for deploying Meshery in production.

Source: /pr-preview/pr-21670/installation/production/

Meshery is a powerful, extensible engineering platform for the collaborative
design and operation of cloud and cloud native infrastructure. Running it in
production—as a shared, always-on management plane that teams depend on—calls
for deliberate planning across reliability, scalability, security, and
operability.

This documentation set collects the considerations, best practices, known
caveats, and hardening guidance you need to deploy and operate Meshery with
confidence in a production environment. It does not replace the
[installation guides](/pr-preview/pr-21670/installation/); rather, it builds on them with the
"what to think about and why" that production demands.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Who this guide is for</div>


Platform engineers, SREs, and operators who are standing up Meshery as a
durable, multi-user service—whether self-hosted in a single cluster, spread
across multiple clusters and clouds, or run out-of-cluster alongside the
infrastructure it manages.
</div>


## How Meshery is deployed

Meshery deploys as a set of containers that can run on a Docker host or inside a
Kubernetes cluster. Any given deployment is described as either _in-cluster_
(Meshery runs inside a cluster it also manages) or _out-of-cluster_ (Meshery
runs separately from the clusters it manages). A single Meshery Server can
manage many clusters concurrently, across one or more clouds.

For an authoritative description of each component and how the pieces fit
together, start with the [Meshery Architecture](/pr-preview/pr-21670/concepts/architecture/)
reference. The
[Deployment Models & Reference Architecture](/pr-preview/pr-21670/installation/production/deployment-models/)
page in this set translates that architecture into production topology
decisions.

<a href="/pr-preview/pr-21670/images/meshery-architecture.webp" class="lightbox-image">
<img src="/pr-preview/pr-21670/images/meshery-architecture.webp" width="55%" /></a>

_Figure: Meshery deploys inside or outside of a Kubernetes cluster and manages one or more clusters._

## The production considerations, by area

This set is organized so you can read it end-to-end or jump to the area you are
working on. Each page is self-contained and cross-links to the relevant
reference material.

| Area | What it covers |
| :--- | :--- |
| **[Deployment Models & Reference Architecture](/pr-preview/pr-21670/installation/production/deployment-models/)** | In-cluster vs. out-of-cluster, Docker vs. Kubernetes, component inventory and statefulness, and the single-cluster, multi-cluster, and multi-cloud topology patterns. |
| **[Infrastructure, Sizing & Performance](/pr-preview/pr-21670/installation/production/infrastructure-sizing-and-performance/)** | Resource requirements per component, capacity planning, MeshSync tiered discovery, Broker throughput, scalability levers, and known performance bounds. |
| **[High Availability & Resiliency](/pr-preview/pr-21670/installation/production/high-availability-and-resiliency/)** | Replication, health probes, failure modes and recovery, the ephemeral database, Remote Provider persistence, and backup & disaster recovery posture. |
| **[Networking & Connectivity](/pr-preview/pr-21670/installation/production/networking-and-connectivity/)** | Network port and directional-flow matrix, ingress and Emissary configuration, secure WebSocket support, CDN and edge caching of the UI, Broker exposure, egress, and network policies. |
| **[Security Hardening](/pr-preview/pr-21670/installation/production/security-hardening/)** | RBAC and least privilege, pod and container security contexts, secret and kubeconfig handling, TLS, supply-chain integrity, namespace isolation, and [trusting an extension](/pr-preview/pr-21670/installation/production/security-hardening/#trusting-an-extension). |
| **[Authentication, Authorization & Identity](/pr-preview/pr-21670/installation/production/authentication-and-identity/)** | Why to preselect a Remote Provider over the Local Provider, OAuth callback configuration, identity providers, and keys/permissions. |
| **[Multi-Cluster & Multi-Cloud Operations](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/)** | Managed vs. unmanaged cluster connections, one Operator per cluster, kubeconfig and context management, MeshSync modes, and cloud-specific guidance. |
| **[Monitoring, Observability & Health KPIs](/pr-preview/pr-21670/installation/production/monitoring-observability-and-kpis/)** | Health endpoints, the key performance indicators of Meshery's health, metrics, tracing, centralized logging, and alerting. |
| **[Operational Readiness Checklist & Known Caveats](/pr-preview/pr-21670/installation/production/operational-readiness-checklist/)** | A consolidated, actionable checklist across every dimension, plus upgrade strategy and the caveats to plan around. |

## Production-readiness principles

A few principles recur throughout this set. Keep them in mind as you make
deployment decisions:

1. **Treat the Meshery database as a cache, not a system of record.** Meshery's
   on-disk database is ephemeral and tied to the lifetime of its Server
   instance. Durable, long-term state lives with a
   [Remote Provider](/pr-preview/pr-21670/reference/extensibility/providers/). Design accordingly.
2. **Prefer a preselected Remote Provider in production.** Pinning a Remote
   Provider avoids unauthenticated Local Provider sessions and lets you control
   which identity providers are accepted. See
   [Authentication, Authorization & Identity](/pr-preview/pr-21670/installation/production/authentication-and-identity/).
3. **Right-size for discovery, not just for traffic.** Meshery's footprint is
   driven as much by the size and number of clusters it discovers (via
   [MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/)) as by user-facing API load.
4. **Make connectivity explicit.** Know which ports flow in which direction
   between Meshery Server, the Broker, and each managed cluster—especially for
   out-of-cluster and multi-cloud topologies.
5. **Observe the management plane itself.** Meshery exposes health endpoints and
   metrics; treat them as first-class signals and alert on them.
6. **Treat every enabled extension as trusted code.** Meshery's
   [extension points](/pr-preview/pr-21670/reference/extensibility/) - adapters,
   Providers, models and integrations, and UI extensions - run inside your
   deployment's trust boundary, without a sandbox between them and Meshery. Enable
   only what you need, from publishers you trust, at versions you pin. See
   [Trusting an extension](/pr-preview/pr-21670/installation/production/security-hardening/#trusting-an-extension).

## Before you begin

If you have not yet chosen an installation method, review the
[Installation Overview](/pr-preview/pr-21670/installation/overview/) and the
[platform-specific guides](/pr-preview/pr-21670/installation/). For Kubernetes production
deployments, the [Helm chart](/pr-preview/pr-21670/installation/kubernetes/helm/) is the recommended
path and is referenced throughout this set. For runtime configuration, keep the
[Meshery Server Environment Variables](/pr-preview/pr-21670/installation/advanced/environment-variables/)
reference close at hand.

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/design-meshery-mcp-server-architecture-and-registration-interface/7954" target="_blank" rel="noopener noreferrer">
            Design: Meshery MCP Server architecture and registration interface
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/the-meshery-mcp-server-foundation-is-up-lets-agree-on-what-to-build-next/7952" target="_blank" rel="noopener noreferrer">
            The Meshery MCP Server foundation is up, let&#39;s agree on what to build next
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-intro-topic/7975" target="_blank" rel="noopener noreferrer">
            New intro topic
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-mcp-server-poc-an-ai-agent-managing-kubernetes-through-mcp/7974" target="_blank" rel="noopener noreferrer">
            Meshery MCP Server POC: an AI agent managing Kubernetes through MCP
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/approach-for-context-window-aware-retrieval-in-the-ai-adapter-issue-20994/7963" target="_blank" rel="noopener noreferrer">
            Approach for context-window-aware retrieval in the AI Adapter (Issue #20994)
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/there-is-some-error-in-running-the-localhost-of-layer5-in-my-server/7818" target="_blank" rel="noopener noreferrer">
            There is some error in running the localhost of layer5 in my server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/rfc-aligning-the-foundation-for-the-meshery-mcp-server/7913" target="_blank" rel="noopener noreferrer">
            RFC: Aligning the Foundation for the Meshery MCP Server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-august-12th-2026/7948" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | August 12th, 2026
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
