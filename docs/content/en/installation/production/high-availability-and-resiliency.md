---
title: High Availability & Resiliency
linkTitle: High Availability
description: >-
  Replication, health probes, failure modes and recovery, the ephemeral
  database, Remote Provider persistence, and backup & disaster-recovery posture
  for resilient Meshery operation.
categories: [installation]
weight: 30
aliases:
- /installation/production/ha
- /installation/production/resiliency
---

Resiliency for Meshery has two layers: the platform-level self-healing that
Kubernetes provides, and the application-level recovery behavior built into
Meshery's components. This page explains how Meshery behaves under failure, how
to configure it for high availability, and—critically—where durable state lives
so that a lost Server is a recoverable event rather than a data-loss event.

## The resiliency model in one idea

Meshery is designed so that its **runtime state is reconstructable**:

- The Meshery Server database is an **ephemeral cache**, tied to the Server's
  lifetime.
- The [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) holds messages
  **in-memory** until consumed; brief interruptions are bridged by NATS topic
  persistence, but the Broker is not a durable store.
- **Durable, long-term state lives with a
  [Remote Provider]({{< ref "installation/production/authentication-and-identity.md" >}})**:
  user identity, environments, saved designs, and preferences.
- [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) **re-populates** the cluster
  snapshot on reconnect through continuous discovery.

The practical consequence: if Meshery Server is restarted or replaced, users'
durable work is safe with the Remote Provider, and the live view of each cluster
is rebuilt automatically by MeshSync. Design your HA and DR around this model
rather than around protecting the local database.

{{% alert title="Use a Remote Provider for production resiliency" color="warning" %}}
Running production on the **Local Provider** means user state has no durable
home beyond the ephemeral database, and sessions are unauthenticated. Preselect
a Remote Provider—see
[Authentication, Authorization & Identity]({{< ref "installation/production/authentication-and-identity.md" >}}).
This is the single most important resiliency decision.
{{% /alert %}}

## Health probes and self-healing

Meshery Server implements Kubernetes-compliant health endpoints, and the Helm
chart wires them into liveness, readiness, and (optional) startup probes:

| Endpoint | Probe | Meaning |
| :--- | :--- | :--- |
| `/healthz/live` | Liveness | Server is running and responsive; provider capabilities are loaded. A failing liveness probe restarts the pod. |
| `/healthz/ready` | Readiness | Server is ready to accept traffic. A failing readiness probe removes the pod from the Service until it recovers. |

The chart's defaults account for Meshery's startup characteristics—notably the
time to initialize and load provider capabilities:

- **Liveness:** `initialDelaySeconds: 80` (allow for server startup and provider
  initialization), `periodSeconds: 12`, `failureThreshold: 4`.
- **Readiness:** `initialDelaySeconds: 10`, `periodSeconds: 4`,
  `failureThreshold: 4` (frequent checks for faster readiness detection).
- **Startup probe:** disabled by default; enable it for slow-starting
  environments or initial installs where provider setup takes longer.

You can verify health directly, including a verbose breakdown:

```bash
kubectl exec --namespace meshery deployment/meshery -- \
  curl -s 'http://localhost:8080/healthz/ready?verbose=1'
```

```
[+]capabilities ok
[i]extension extension package found
healthz check passed
```

`[+]` indicates a passed check, `[-]` a failed check (which marks the pod
unhealthy), and `[i]` an informational status that does not affect health.

{{% alert title="Tune probes for upgrades" color="info" %}}
During upgrades, capability reloading can make the Server temporarily
unavailable. The chart provides `values-upgrade.yaml` with a startup probe and
higher failure thresholds so pods are not killed mid-initialization. Use it for
upgrades—see
[Operational Readiness]({{< ref "installation/production/operational-readiness-checklist.md" >}})
and the [Helm guide]({{< ref "installation/kubernetes/helm.md" >}}).
{{% /alert %}}

## High availability for Meshery Server

Run Meshery Server on Kubernetes and let the platform reschedule and self-heal
failed pods. To raise availability further:

- **Spread across failure domains.** Use `nodeSelector`, `affinity`
  (anti-affinity), and `tolerations` so replicas do not share a node or zone.
  This is the highest-value HA step for a single-region deployment.
- **Keep probes correctly tuned** so traffic is only routed to ready pods and
  unhealthy pods are restarted promptly.
- **Use a Remote Provider** so a pod loss never risks durable user state.

### Replication and the database caveat

The Meshery database is a **single-writer SQLite/Bitcask file cache**, not a
shared multi-writer datastore. This is the central caveat for horizontal
scaling:

- Increasing `replicaCount` adds request-handling capacity, but replicas do not
  transparently share one consistent local database. Durable consistency comes
  from the Remote Provider, not from the local cache.
- For most production needs, prefer a **well-resourced, probe-protected Server**
  (vertical headroom) plus Kubernetes self-healing, with a Remote Provider for
  durability—rather than assuming many Server replicas behave like a clustered
  database.
- If you scale out replicas, validate behavior under your workload and front
  them with the Service/ingress; treat the local database as a per-instance
  cache that MeshSync and the provider re-populate.
- Because each replica keeps its own local cache, enable **session affinity
  (sticky sessions)** at the ingress or load balancer so a user's requests
  consistently reach the same Server instance and the UI presents a stable
  view rather than alternating between per-replica caches.

{{% alert title="Plan replication around the cache model" color="warning" %}}
Do not treat the local database as a shared source of truth across replicas. Let
the [Remote Provider]({{< ref "installation/production/authentication-and-identity.md" >}}) be
the system of record and the local database be a per-Server cache. This keeps
scaling decisions safe.
{{% /alert %}}

## Per-component failure behavior

| Component | Failure behavior | Recovery |
| :--- | :--- | :--- |
| **Meshery Server** | Pod restarted by liveness probe / Kubernetes. Local cache rebuilt. | Durable state from Remote Provider; cluster snapshot rebuilt by MeshSync. |
| **Meshery Operator** | Stateless; rescheduled by Kubernetes. | Reconciles Broker and MeshSync back to desired state on restart. |
| **MeshSync** | Stateless; managed by the Operator. | Resumes discovery and re-syncs the snapshot on restart. |
| **Meshery Broker** | Deployed as a StatefulSet; rescheduled on failure. In-memory messages not yet consumed may be lost. | NATS topic persistence bridges brief interruptions; MeshSync re-syncs to restore the snapshot. |
| **Meshery Adapter** | Stateless; if the Server is unavailable it backs off and retries to reconnect perpetually. | Reconnects and re-registers capabilities automatically. |
| **Remote Provider** | External dependency for auth and durable state. | Provided by the Remote Provider's own availability; see below. |

### Adapter reconnection

Meshery Adapters register with Meshery Server over HTTP POST. If the Server is
unavailable, an Adapter backs off and retries to connect **perpetually**, then
re-registers its capabilities once the Server returns. No manual intervention is
required after a Server restart.

### Operator-managed recovery

The [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) continuously reconciles
the Broker and MeshSync toward their desired state. If the Broker or MeshSync is
disrupted, the Operator restores it. This is why running the Operator per
cluster is important for resilient discovery.

## Remote Provider as a dependency

When you preselect a Remote Provider, it becomes part of your availability and
recovery model:

- Meshery Server authenticates users and stores durable state through the
  provider, so the provider's availability affects login and durable reads.
- Ensure egress connectivity to the provider is reliable (see
  [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}))
  and include the provider in your monitoring and incident runbooks.
- The on-disk cache lets Meshery continue serving already-loaded data during
  brief provider blips, but treat sustained provider unavailability as a
  first-class incident.

## Backup and disaster recovery

Because durable state lives with the Remote Provider and the local database is a
cache, DR planning is refreshingly simple:

1. **Durable user data** — covered by your Remote Provider's backup/DR
   guarantees. Confirm them with your provider.
2. **Deployment configuration** — keep Helm `values.yaml`, environment-variable
   configuration, kubeconfig/context definitions, and any MeshSync
   `informer_config` customization under version control (GitOps). This _is_
   your Meshery "backup"; it lets you recreate an identical deployment.
3. **Cluster snapshot** — not backed up; it is rebuilt by MeshSync on
   reconnect.
4. **Recovery drill** — periodically validate that you can redeploy Meshery from
   source-controlled configuration, reconnect clusters, and authenticate against
   the Remote Provider.

{{% alert title="Your real backup is your configuration" color="info" %}}
Snapshotting the local database is unnecessary and misleading—it is ephemeral.
Version-control your Helm values, environment configuration, and connection
definitions instead. Recreating Meshery from that configuration restores the
deployment; the provider restores user data; MeshSync restores the live view.
{{% /alert %}}

## Resiliency checklist

- [ ] Meshery Server runs on Kubernetes with liveness/readiness probes enabled
      and tuned.
- [ ] A **Remote Provider is preselected** (no Local Provider in production).
- [ ] Replicas/scheduling use anti-affinity to avoid shared nodes/zones.
- [ ] Server memory limits have headroom for the MeshSync snapshot.
- [ ] Broker has memory headroom for reconnect/resync bursts.
- [ ] Egress to the Remote Provider is reliable and monitored.
- [ ] Helm values, env config, and connection definitions are in version
      control.
- [ ] Upgrade procedure uses startup-probe-friendly values
      (`values-upgrade.yaml`).
- [ ] A redeploy-from-config recovery drill has been validated.

{{< related-discussions tag="meshery" >}}
