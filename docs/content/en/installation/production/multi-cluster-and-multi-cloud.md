---
title: Multi-Cluster & Multi-Cloud Operations
linkTitle: Multi-Cluster & Multi-Cloud
description: >-
  Managed vs. unmanaged cluster connections, one Operator per cluster,
  kubeconfig and context management, MeshSync deployment modes, and
  cloud-specific guidance for operating Meshery across many clusters and clouds.
categories: [installation]
weight: 70
aliases:
- /installation/production/multi-cluster
- /installation/production/multi-cloud
---

A core reason to run Meshery as a central management plane is to operate many
clusters—often across multiple clouds—from one place. This page covers how
Meshery connects to clusters, the difference between managed and unmanaged
connections, how discovery is deployed per cluster, and the cloud-specific
details that decide whether a fleet "just works."

## How Meshery connects to a cluster

A Meshery–cluster connection is established by providing Meshery Server access to
the cluster's Kubernetes API:

- **In-cluster**, Meshery uses its in-cluster ServiceAccount for the cluster it
  runs in.
- **Out-of-cluster**, Meshery uses a **kubeconfig context** per cluster.
  `mesheryctl system start` or uploading a kubeconfig in the UI creates the
  connection; each context becomes a managed connection.

On connection, Meshery deploys one [Operator]({{< ref "concepts/architecture/operator/index.md" >}})
per cluster (unless you use embedded MeshSync—see below), which manages that
cluster's [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) and
[Broker]({{< ref "concepts/architecture/broker/index.md" >}}). On disconnect, those components are
removed. The Operator is also controllable from the UI's on/off switch
independently of the connection.

## Managed vs. unmanaged connections

"Managed" and "unmanaged" can mean two different—and both relevant—things in a
multi-cluster context. Be clear about which you mean:

### Meshery-managed discovery vs. library (embedded) discovery

| Connection style | What's deployed into the cluster | Trade-offs |
| :--- | :--- | :--- |
| **Operator-managed** | Operator + MeshSync + Broker run in the managed cluster. | Full, event-driven discovery with in-cluster components; the Operator self-heals them. Requires permission to deploy into the cluster. |
| **Embedded (library, default)** | Nothing—MeshSync runs as a library inside Meshery Server for that connection. | No in-cluster footprint; useful where you can't or won't deploy the Operator. Shifts discovery work into the Server process. |

The mode for new connections is governed by
`MESHSYNC_DEFAULT_DEPLOYMENT_MODE` (`operator` or `embedded`), which **defaults
to `embedded`**, and the mode can be switched per connection on the connections
page. Switching from operator to
embedded undeploys the in-cluster components and starts the in-Server routine;
switching back redeploys them. Choose per cluster based on whether in-cluster
deployment is acceptable and on the Server's capacity to absorb embedded
discovery (see
[Infrastructure, Sizing & Performance]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}})).

### Cloud-managed vs. self-managed Kubernetes

Separately, the clusters themselves may be **cloud-managed** (EKS, GKE, AKS,
and similar) or **self-managed**. Meshery connects to both the same way, but
cloud-managed clusters differ in how they surface node permissions and
load-balanced endpoints, covered under [cloud-specific notes](#cloud-specific-guidance).

You can also disable Operator deployment entirely for a deployment with
`DISABLE_OPERATOR=true`, which prevents Meshery from automatically deploying the
Operator into connected clusters—useful when discovery is handled in embedded
mode or by policy.

## Managing kubeconfig and contexts at scale

For a fleet, kubeconfig/context management is the operational backbone:

- **One context per cluster**, each scoped to a **least-privilege credential**
  for that cluster (see
  [Security Hardening]({{< ref "installation/production/security-hardening.md" >}})). Avoid a
  single all-powerful credential spanning the fleet.
- Mount kubeconfig from a Secret and point Meshery at it with
  `KUBECONFIG_FOLDER` (default `~/.kube`). Keep context names stable and
  meaningful.
- Treat the **set of connections as version-controlled configuration** so the
  fleet can be reconstructed during recovery (see
  [High Availability & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}})).
- Prefer short-lived or provider-issued credentials where the cloud supports
  them, and rotate per-cluster credentials independently.

## Reaching each cluster's Broker

In multi-cluster/out-of-cluster operation, Meshery Server must reach each
cluster's Broker on `4222/tcp`. The Operator publishes a reachable `external`
endpoint into the Broker custom resource `status`, selecting (in order of
preference) the LoadBalancer hostname, the LoadBalancer IP, the kubeconfig host
with NodePort, the ClusterIP with cluster port, or a worker node IP with
NodePort.

The practical implications across a fleet:

- Ensure each cluster's Broker Service type (`LoadBalancer` or `NodePort`) is
  supported and that the resulting endpoint is reachable from the Server.
- Restrict that exposure to the Server's network origin (security groups,
  load-balancer source ranges, private connectivity). See
  [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}).

{{% alert title="The #1 multi-cloud pitfall" color="warning" %}}
The most common cross-cloud failure is an unreachable Broker endpoint: the
cluster publishes a LoadBalancer hostname or NodePort the central Server cannot
reach (blocked by security groups, private subnets, or missing routes).
Validate Broker reachability from the Server for **every** cluster you add.
{{% /alert %}}

## Cloud-specific guidance

The architecture is identical across clouds; these are the per-provider details
that matter:

- **Node-watch RBAC.** Full discovery on AKS, AWS, and GCP may require
  permission to watch nodes. Enable `rbac.nodes: true` on those clusters (it
  defaults to `false`). Grant it only where needed.
- **Load-balanced Broker endpoints.** Clouds differ in whether a `LoadBalancer`
  Service surfaces a **hostname** (commonly AWS ELB) or an **IP** (commonly GCP/
  Azure). The Operator's endpoint selection handles both, but your firewall and
  reachability checks must account for the form your cloud uses.
- **Private clusters.** For private API servers or nodes (private EKS/GKE/AKS),
  the central Server needs network reachability to both the API server and the
  Broker—via VPC/VNet peering, private load balancers, transit gateways, or a
  VPN. Public exposure is discouraged (see
  [Security Hardening]({{< ref "installation/production/security-hardening.md" >}})).
- **Cross-region latency.** A central Server managing distant clusters incurs
  latency on discovery streaming and API calls. Keep it acceptable, and consider
  regional management planes if a single central plane spans high-latency links.

## Operating a fleet

- **Per-cluster connection health.** Each connection's chip in the UI reflects
  live connectivity; Broker/Operator/MeshSync follow the connection lifecycle.
  Watch these as fleet KPIs (see
  [Monitoring, Observability & Health KPIs]({{< ref "installation/production/monitoring-observability-and-kpis.md" >}})).
- **Per-cluster discovery scope.** Tune MeshSync's `informer_config` blacklist
  per cluster to control footprint on large clusters
  ([sizing]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}})).
- **Blast-radius isolation.** Distinct per-cluster credentials and network
  policies mean a problem in one cluster's connection does not cascade across
  the fleet.
- **Consistent lifecycle.** Keep the Operator/MeshSync mode consistent with your
  policy across clusters, and codify it via `MESHSYNC_DEFAULT_DEPLOYMENT_MODE`
  and your connection configuration.

## Multi-cluster checklist

- [ ] One kubeconfig context per cluster, each with a least-privilege
      credential.
- [ ] Broker endpoint reachable from the Server and locked to the Server's
      origin for every cluster.
- [ ] `rbac.nodes` enabled only on clusters that require node watching.
- [ ] MeshSync mode (operator vs. embedded) chosen deliberately per cluster.
- [ ] Private connectivity (peering/VPN) for private clusters; no broad public
      Broker exposure.
- [ ] Cross-region latency assessed; regional planes considered if needed.
- [ ] Connection set and per-cluster config under version control for recovery.

{{< related-discussions tag="meshery" >}}
