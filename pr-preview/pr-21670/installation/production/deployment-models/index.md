# Deployment Models & Reference Architecture

> In-cluster vs. out-of-cluster, Docker vs. Kubernetes, the Meshery component inventory and its statefulness, and the topology patterns for single-cluster, multi-cluster, and multi-cloud production deployments.

Source: /pr-preview/pr-21670/installation/production/deployment-models/

The first production decision is _where_ Meshery runs relative to the
infrastructure it manages, and _how_ its components are hosted. This page frames
those choices and maps the Meshery component inventory onto production
topologies. It complements the [Meshery Architecture](/pr-preview/pr-21670/concepts/architecture/)
reference with the trade-offs that matter when the deployment is shared and
long-lived.

## Where Meshery runs: in-cluster vs. out-of-cluster

Every Meshery deployment is either _in-cluster_ or _out-of-cluster_ with respect
to the clusters it manages. The distinction is not about how Meshery is
packaged—it is about the network and trust relationship between Meshery Server
and the clusters under management.

| Topology | What it means | Best suited to |
| :--- | :--- | :--- |
| **In-cluster** | Meshery Server runs as a workload inside a Kubernetes cluster that it also manages. | Teams standardizing on Kubernetes who want Meshery co-located with their primary cluster and managed by the same platform tooling (Helm, GitOps). |
| **Out-of-cluster** | Meshery Server runs separately—on a Docker host or in a dedicated cluster—and manages one or more remote clusters over the network. | Central management planes that manage many clusters, manage clusters across clouds, or must remain available independently of any single managed cluster. |

A single Meshery Server commonly does both at once: it runs _in_ one cluster
while managing that cluster _and_ additional remote clusters
out-of-cluster. The practical consequences of the choice are mostly about
[networking](/pr-preview/pr-21670/installation/production/networking-and-connectivity/) (how Meshery
Server reaches each cluster's API server and Broker) and
[security](/pr-preview/pr-21670/installation/production/security-hardening/) (where credentials live
and what is exposed).

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Out-of-cluster Broker reachability</div>


When Meshery Server is out-of-cluster, it must be able to reach the
[Meshery Broker](/pr-preview/pr-21670/concepts/architecture/broker/) running inside each managed
cluster. The Broker's externally reachable endpoint is derived from the
cluster's Service type (LoadBalancer, NodePort, or ClusterIP). Plan Broker
exposure deliberately—see
[Networking & Connectivity](/pr-preview/pr-21670/installation/production/networking-and-connectivity/)
and [Multi-Cluster & Multi-Cloud](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/).
</div>


## How Meshery is hosted: Docker vs. Kubernetes

Meshery deploys as a set of containers to either a Docker host or a Kubernetes
cluster.

- **Docker host (via `mesheryctl` / Docker Compose).** Simple to stand up; a
  good fit for single-node, single-operator out-of-cluster deployments. High
  availability is bounded by the single host, so it is generally not the choice
  for mission-critical, multi-user production.
- **Kubernetes (via Helm).** The recommended path for production. Kubernetes
  supplies the scheduling, self-healing, rolling upgrades, and scaling
  primitives that Meshery's high-availability and resiliency posture builds on.
  The [Helm chart](/pr-preview/pr-21670/installation/kubernetes/helm/) is referenced throughout this
  set.

For production, prefer Kubernetes with Helm unless you have a specific reason to
run on a single Docker host.

## Component inventory and statefulness

Knowing which components hold state—and what kind—is foundational to sizing,
high availability, backup, and upgrade planning. The following table summarizes
Meshery's components and their persistence characteristics.

| Component | Role | Persistence |
| :--- | :--- | :--- |
| **Meshery Server** | Core control plane: REST/GraphQL APIs, orchestration, capability registry. | Caches state in an on-disk database under `~/.meshery/`. Treat as a cache. |
| **Meshery Database** | SQLite/Bitcask file database backing Server state and the MeshSync snapshot. | Ephemeral; tied to the Server instance lifetime. Not a system of record. |
| **Meshery UI** | React/Next.js web interface served by the Server. | Stateless. |
| **Meshery Operator** | Kubernetes operator that manages MeshSync and Broker; one per managed cluster. | Stateless. |
| **MeshSync** | Custom controller performing continuous, tiered discovery of cluster state. | Stateless; its working snapshot lives in the Server's database. |
| **Meshery Broker** | NATS-based event bus streaming discovery and events between clusters and Server. | Messages persist **in-memory only** until consumed; no persistent volume. |
| **Meshery Adapters** | Optional, capability-specific integrations registered with the Server. | Stateless; transactional interaction with infrastructure. |
| **`mesheryctl`** | Command-line client. | Stateless; has a local config file. |
| **Remote Provider** | Pluggable extension supplying identity and **durable, long-term persistence**. | Stateful—the system of record for user data, environments, and saved work. |

Two implications dominate production planning:

1. **The Meshery Server database is a cache.** Because both the on-disk database
   and the Broker's in-memory queue are ephemeral, durable state must come from
   a [Remote Provider](/pr-preview/pr-21670/reference/extensibility/providers/). A Server can be
   destroyed and recreated; what users care about persists with the provider,
   and MeshSync re-populates the cluster snapshot on reconnect.
2. **Discovery scope drives footprint.** MeshSync continuously mirrors cluster
   resources into the Server's database. The number and size of managed clusters
   therefore influence Server memory and database size as much as user traffic
   does. See
   [Infrastructure, Sizing & Performance](/pr-preview/pr-21670/installation/production/infrastructure-sizing-and-performance/).

## One Operator (and Broker and MeshSync) per cluster

For each Kubernetes cluster under management, Meshery deploys exactly one
[Meshery Operator](/pr-preview/pr-21670/concepts/architecture/operator/), which in turn manages one
[MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/) and one
[Broker](/pr-preview/pr-21670/concepts/architecture/broker/) in that cluster. This holds whether
Meshery Server is in-cluster or out-of-cluster. The Operator is deployed when
Meshery Server connects to a cluster and removed when it disconnects (or via the
on/off control in the UI).

MeshSync can instead run in **embedded mode**—the default for new Kubernetes
connections—where it executes as a library inside Meshery Server and deploys no
resources into the managed cluster. The mode is chosen per connection and has
meaningful production trade-offs covered in
[Multi-Cluster & Multi-Cloud Operations](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/).

## Topology patterns

### Single cluster, in-cluster

Meshery runs in the same cluster it manages. The Server reaches the Broker over
in-cluster networking (ClusterIP), and credentials are typically the in-cluster
ServiceAccount. This is the simplest production topology and a strong default
for teams with one primary cluster.

### Central management plane, many clusters (out-of-cluster)

A dedicated Meshery deployment manages a fleet of clusters over the network.
Each managed cluster runs its own Operator/MeshSync/Broker, and Meshery Server
holds a kubeconfig context per cluster. This is the canonical pattern for
platform teams. It concentrates availability requirements on the management
plane and concentrates connectivity requirements on reaching each cluster's API
server and Broker endpoint.

### Multi-cloud

A special case of the central management plane in which managed clusters span
providers (for example EKS, GKE, and AKS). The differences are practical rather
than architectural: per-cloud node-watch RBAC, how each cloud surfaces a
reachable Broker endpoint (LoadBalancer hostname vs. IP vs. NodePort), and
cross-region latency. See
[Multi-Cluster & Multi-Cloud Operations](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/).

### Highly available management plane

Any of the above can be made highly available by running Meshery Server with
appropriate replication and health probes on Kubernetes, and by relying on a
Remote Provider for durable state. The mechanics—and the caveats around the
single-writer SQLite database and the in-memory Broker—are covered in
[High Availability & Resiliency](/pr-preview/pr-21670/installation/production/high-availability-and-resiliency/).

## Choosing a model

Use this quick guide as a starting point:

- **One Kubernetes cluster, one team** → in-cluster on Kubernetes via Helm.
- **Many clusters or many teams** → out-of-cluster central management plane on
  Kubernetes, with a preselected Remote Provider.
- **Clusters across clouds** → out-of-cluster, plus the multi-cloud guidance.
- **Evaluation or a single operator on a workstation** → Docker via
  `mesheryctl` (not recommended as a durable, multi-user production target).

Whatever you choose, the remaining pages in this set apply: size the
deployment, harden it, make it resilient, and observe it.

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
