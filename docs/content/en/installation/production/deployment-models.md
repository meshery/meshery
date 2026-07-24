---
title: Deployment Models & Reference Architecture
linkTitle: Deployment Models
description: >-
  In-cluster vs. out-of-cluster, Docker vs. Kubernetes, the Meshery component
  inventory and its statefulness, and the topology patterns for single-cluster,
  multi-cluster, and multi-cloud production deployments.
categories: [installation]
weight: 10
aliases:
- /installation/production/reference-architecture
- /installation/production/topologies
---

The first production decision is _where_ Meshery runs relative to the
infrastructure it manages, and _how_ its components are hosted. This page frames
those choices and maps the Meshery component inventory onto production
topologies. It complements the [Meshery Architecture]({{< ref "concepts/architecture/_index.md" >}})
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
[networking]({{< ref "installation/production/networking-and-connectivity.md" >}}) (how Meshery
Server reaches each cluster's API server and Broker) and
[security]({{< ref "installation/production/security-hardening.md" >}}) (where credentials live
and what is exposed).

{{% alert title="Out-of-cluster Broker reachability" color="warning" %}}
When Meshery Server is out-of-cluster, it must be able to reach the
[Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) running inside each managed
cluster. The Broker's externally reachable endpoint is derived from the
cluster's Service type (LoadBalancer, NodePort, or ClusterIP). Plan Broker
exposure deliberately—see
[Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}})
and [Multi-Cluster & Multi-Cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}).
{{% /alert %}}

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
  The [Helm chart]({{< ref "installation/kubernetes/helm.md" >}}) is referenced throughout this
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
   a [Remote Provider]({{< ref "reference/extensibility/providers/index.md" >}}). A Server can be
   destroyed and recreated; what users care about persists with the provider,
   and MeshSync re-populates the cluster snapshot on reconnect.
2. **Discovery scope drives footprint.** MeshSync continuously mirrors cluster
   resources into the Server's database. The number and size of managed clusters
   therefore influence Server memory and database size as much as user traffic
   does. See
   [Infrastructure, Sizing & Performance]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}}).

## One Operator (and Broker and MeshSync) per cluster

For each Kubernetes cluster under management, Meshery deploys exactly one
[Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}), which in turn manages one
[MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) and one
[Broker]({{< ref "concepts/architecture/broker/index.md" >}}) in that cluster. This holds whether
Meshery Server is in-cluster or out-of-cluster. The Operator is deployed when
Meshery Server connects to a cluster and removed when it disconnects (or via the
on/off control in the UI).

MeshSync can instead run in **embedded mode**—the default for new Kubernetes
connections—where it executes as a library inside Meshery Server and deploys no
resources into the managed cluster. The mode is chosen per connection and has
meaningful production trade-offs covered in
[Multi-Cluster & Multi-Cloud Operations]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}).

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
[Multi-Cluster & Multi-Cloud Operations]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}).

### Highly available management plane

Any of the above can be made highly available by running Meshery Server with
appropriate replication and health probes on Kubernetes, and by relying on a
Remote Provider for durable state. The mechanics—and the caveats around the
single-writer SQLite database and the in-memory Broker—are covered in
[High Availability & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}}).

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

{{< related-discussions tag="meshery" >}}
