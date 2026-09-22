---
title: Infrastructure, Sizing & Performance
linkTitle: Infrastructure & Sizing
description: >-
  Resource requirements per Meshery component, capacity planning, MeshSync
  tiered discovery, Broker throughput, scalability levers, and the known
  performance bounds to plan around in production.
categories: [installation]
weight: 20
aliases:
- /installation/production/sizing
- /installation/production/performance
---

Meshery's resource footprint is driven by two largely independent forces: the
volume of user-facing API and UI traffic, and the scope of infrastructure it
continuously discovers. Sizing a production deployment means accounting for
both. This page provides per-component resource guidance, the levers you can
pull to scale, and the performance bounds to design around.

{{% alert title="No one-size-fits-all numbers" color="info" %}}
The Helm chart intentionally ships with **no default CPU/memory requests or
limits** so Meshery can start on small clusters. That default is not a
production recommendation. Set requests and limits explicitly, measure under
your own load, and iterate. The figures below are starting points, not
guarantees.
{{% /alert %}}

## What drives Meshery's footprint

| Driver | Primary effect | Components affected |
| :--- | :--- | :--- |
| Number and size of managed clusters | Larger MeshSync snapshot; more events; larger database | MeshSync, Broker, Meshery Server, Database |
| Rate of change in managed clusters | Higher event throughput | Broker, MeshSync, Meshery Server |
| Concurrent users and API/GraphQL load | More request handling; capability lookups | Meshery Server |
| Designs, models, and registry size | Larger in-memory registry and database | Meshery Server |
| Relationship/policy evaluation | CPU during evaluation | Meshery Server |

The discovery-driven forces are easy to underestimate. A handful of large,
rapidly changing clusters can place more sustained load on Meshery than a busy
but small user base. Size for your **discovery scope** first.

## Per-component sizing guidance

Use these as conservative starting points for a production deployment, then
right-size from observed utilization. Allocate generously where state is held
(Meshery Server) and modestly where components are stateless and bursty
(Operator, Adapters).

| Component | Starting requests | Starting limits | Notes |
| :--- | :--- | :--- | :--- |
| **Meshery Server** | 500m CPU / 512Mi | 2 CPU / 2Gi | Scale memory with discovery scope, registry size, and database growth. The largest and most important allocation. |
| **Meshery Operator** | 50m CPU / 64Mi | 200m CPU / 256Mi | One per managed cluster; lightweight reconciliation. |
| **MeshSync** | 100m CPU / 128Mi | 500m CPU / 512Mi+ | Scale with cluster size and rate of change; the heaviest discovery component. |
| **Meshery Broker (NATS)** | 100m CPU / 128Mi | 500m CPU / 512Mi+ | Memory tracks in-flight, unconsumed messages. Bursts during reconnects and large resyncs. |
| **Meshery Adapter** (each, optional) | 50m CPU / 64Mi | 200m CPU / 256Mi | Only if you deploy adapters; stateless and transactional. |

Storage:

- **Meshery Server / Database.** Provision disk for the on-disk cache under the
  Server's data folder. The database grows with discovery scope and registry
  size. Because it is a cache (not a system of record), size for working-set
  performance rather than long-term retention—durable data lives with the
  [Remote Provider]({{< ref "installation/production/authentication-and-identity.md" >}}).
- **Meshery Broker.** No persistent volume is required; the Broker holds
  messages in memory until consumed. Size **memory**, not disk, for the Broker.

{{% alert title="Set Server memory limits with headroom" color="warning" %}}
Because MeshSync's working snapshot is held in the Server's database and partly
in memory, an undersized memory limit on Meshery Server can lead to restarts
under large or churny clusters. Give the Server generous memory headroom and
alert on its utilization—see
[Monitoring, Observability & Health KPIs]({{< ref "installation/production/monitoring-observability-and-kpis.md" >}}).
{{% /alert %}}

## MeshSync: tiered discovery and scoping

[MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) is the component most sensitive to
the size of managed clusters. It uses **tiered discovery** to progressively
refine identification of resources, balancing granularity against speed and
scalability. Two controls let you bound its cost:

1. **Blacklist uninteresting resources.** MeshSync's discovery is governed by an
   `informer_config` on the `meshsyncs.meshery.io` CRD. Editing the CRD to
   blacklist resource types you do not need to track reduces event volume,
   database growth, and Server memory pressure. Retrieve the CRD, edit
   `informer_config`, and re-apply it:

   ```bash
   kubectl get crd meshsyncs.meshery.io -o yaml > meshsync.yaml
   # edit informer_config to blacklist unwanted resource types
   kubectl apply -f meshsync.yaml
   ```

2. **Choose the right deployment mode.** MeshSync runs in **operator** mode
   (deployed into the managed cluster) or **embedded** mode (a library inside
   Meshery Server, deploying nothing into the cluster). Embedded mode reduces
   in-cluster footprint but shifts discovery work into the Server process; it is
   useful for clusters where you cannot or prefer not to deploy the Operator.
   The default for new connections is set by
   `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` (`embedded` or `operator`). See
   [Multi-Cluster & Multi-Cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}).

For very large clusters, blacklisting noisy, high-cardinality resource types is
the single most effective lever on Meshery's footprint.

## Broker throughput

The [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) (NATS) streams discovery
data and events between each cluster and the Server. Production guidance:

- Run **one Broker per managed cluster**. A single Broker instance can be scaled
  vertically to absorb a cluster's data volume; this is independent of the
  number of clusters.
- Because messages are held **in memory until consumed**, sustained
  Server unavailability or a slow consumer causes Broker memory to grow. Ensure
  Meshery Server keeps up and that the Broker has memory headroom for reconnect
  bursts and large resyncs.
- The Broker recovers messages from its NATS topics across brief connectivity
  interruptions, which smooths transient Server or network blips.

## Meshery Server: API, registry, and policy

- **API/GraphQL load** scales with concurrent users and clients. Horizontal
  replication helps with read/stateless request handling; mind the database
  caveats in
  [High Availability & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}}).
- **Static UI assets** are served by the Server with release-scoped cache
  headers, so a CDN or caching reverse proxy in front of Meshery can offload
  virtually all repeat static-asset requests—worthwhile when serving many
  concurrent UI users. See
  [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}).
- **Registry and models** are held to serve design and relationship operations;
  larger registries increase baseline memory.
- **Relationship/policy evaluation** is CPU-bound and time-boxed by
  `POLICY_EVAL_TIMEOUT` (default `3m`). If you see evaluations timing out on
  large designs, raise the timeout or allocate more CPU. The engine selection
  (`USE_GO_POLICY_ENGINE`) also affects evaluation characteristics. See the
  [environment variables reference]({{< ref "installation/advanced/environment-variables.md" >}}).

## Scalability levers, at a glance

| Lever | Effect | Where configured |
| :--- | :--- | :--- |
| Meshery Server CPU/memory requests & limits | Headroom for traffic, registry, and discovery snapshot | Helm `resources` |
| Meshery Server replicas | More request-handling capacity (see HA caveats) | Helm `replicaCount` + autoscaling |
| MeshSync blacklist (`informer_config`) | Lower event volume and database growth | `meshsyncs.meshery.io` CRD |
| MeshSync mode (operator vs. embedded) | Shifts discovery cost in-cluster vs. into Server | Per connection / `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` |
| Broker memory | Absorbs in-flight message bursts | Helm values for the Broker/Operator |
| `POLICY_EVAL_TIMEOUT`, CPU | Tolerance for large policy evaluations | Env var / Helm |
| CDN / caching reverse proxy for the UI | Offloads repeat static UI asset serving from the Server | Edge in front of Meshery |

## Known performance bounds and caveats

- **The database is a single-writer SQLite/Bitcask cache.** It is excellent for
  a cached working set but is not a horizontally shared, multi-writer datastore.
  This shapes how far a single Server instance scales and how replicas behave
  (see [HA & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}})).
- **Discovery cost is proportional to cluster size and churn.** Without
  blacklisting, very large or high-churn clusters can dominate Server memory and
  database growth.
- **Broker memory is bounded by consumption.** A wedged or far-behind Server
  consumer lets Broker memory climb; alert on it.
- **Policy evaluation is time-boxed.** Large designs may hit
  `POLICY_EVAL_TIMEOUT`; tune CPU and the timeout together.

## Capacity-planning workflow

1. Start from the per-component starting points above.
2. Connect representative clusters and let MeshSync reach steady state.
3. Measure Server memory/CPU, database size, and Broker memory under real
   discovery and user load.
4. Blacklist resource types you do not need; re-measure.
5. Set requests/limits with headroom (especially Server memory) and configure
   autoscaling where appropriate.
6. Add the resulting thresholds to your alerting—see
   [Monitoring, Observability & Health KPIs]({{< ref "installation/production/monitoring-observability-and-kpis.md" >}}).

{{< related-discussions tag="meshery" >}}
