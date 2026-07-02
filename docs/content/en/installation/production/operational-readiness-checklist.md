---
title: Operational Readiness Checklist & Known Caveats
linkTitle: Readiness Checklist & Caveats
description: >-
  A consolidated, actionable production-readiness checklist across every
  dimension, the upgrade strategy for Meshery, and the known caveats and
  limitations to plan around.
categories: [installation]
weight: 90
aliases:
- /installation/production/checklist
- /installation/production/caveats
---

This page brings the whole set together: a single checklist you can work through
before declaring a Meshery deployment production-ready, the upgrade strategy to
adopt, and the known caveats to design around. Each item links back to the page
where it is explained in full.

## Upgrade strategy

Meshery is a composition of components—Server, UI, Operator, MeshSync, Broker,
Adapters, and `mesheryctl`—some of which upgrade together and some
independently. For Kubernetes production deployments:

- **Use Helm and pin versions.** Upgrade with `helm upgrade`, and pin the image
  to an **immutable version tag** rather than tracking `stable-latest`, so
  upgrades are deliberate and reproducible
  ([security hardening]({{< ref "installation/production/security-hardening.md" >}})).
- **Use upgrade-friendly probe values.** Capability reloading can make the
  Server briefly unavailable during an upgrade. The chart provides
  `values-upgrade.yaml` (startup probe, higher failure thresholds) so pods are
  not killed mid-initialization:

  ```bash
  helm upgrade meshery meshery/meshery --namespace meshery \
    -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml \
    --wait --timeout 10m
  ```

  For reproducible production upgrades, download `values-upgrade.yaml`,
  version-control it alongside your own values, and reference your local copy
  (with a pinned chart and image version) rather than fetching from `master` at
  upgrade time.

- **Choose a release channel deliberately.** `stable` for production; `edge`
  only for testing. The channel is reflected by `RELEASE_CHANNEL`.
- **Mind component groupings.** See the
  [Upgrade Guide]({{< ref "installation/upgrades/index.md" >}}) for which components move together
  (e.g., Server/UI/Load Generators/Database) versus independently (Adapters,
  `mesheryctl`). The Operator and its controllers effectively move **with the
  Server**: each Server release pins the operator chart (and operator image)
  it deploys to managed clusters.
- **Let the Server own the Operator.** On Server-managed clusters, do not
  hand-upgrade or hand-configure the `meshery-operator` Helm release — the
  Server's reconciliation re-applies its own pinned chart version and reverts
  manual changes. Upgrading the Server is the supported way to upgrade the
  Operator; the chart's CRD update Job refreshes the `Broker`/`MeshSync` CRD
  schemas on each upgrade (Helm alone never updates CRDs). See
  [How Meshery Server manages Meshery Operator]({{< ref "installation/upgrades/index.md#how-meshery-server-manages-meshery-operator" >}}).
- **Expect CRDs to persist.** The `brokers.meshery.io` and
  `meshsyncs.meshery.io` CRDs (and their objects) deliberately survive
  operator uninstalls and Helm release deletion; include the explicit
  `kubectl delete crd ...` step in decommissioning runbooks only when
  permanent removal is intended.
- **Rehearse and roll back.** Because durable state lives with the Remote
  Provider and the local database is a cache, rolling back the deployment is
  low-risk for data—validate the rollback path anyway. On rollback, managed
  clusters converge back to the older Server's pinned operator chart
  automatically.
- **Edge caches need no purge.** If a CDN or caching reverse proxy fronts
  Meshery, its UI cache busts itself on upgrade—hashed asset URLs change and the
  HTML `ETag` follows the build/release version—so a manual cache purge is not
  required in your upgrade pipeline
  ([networking]({{< ref "installation/production/networking-and-connectivity.md" >}})).

{{% alert title="Upgrades and the cache" color="info" %}}
An upgrade that recreates the Server pod discards the local cache; this is
expected and safe. The Remote Provider preserves durable state and MeshSync
rebuilds the cluster snapshot. Do not block upgrades trying to preserve the
ephemeral database.
{{% /alert %}}

## Known caveats and limitations

Plan around these known characteristics—they are by design and shape good
production architecture:

- **The database is an ephemeral, single-writer cache.** It is SQLite/Bitcask on
  local disk, not a shared multi-writer datastore. Durable state must come from
  a [Remote Provider]({{< ref "installation/production/authentication-and-identity.md" >}});
  this also shapes how Server replication behaves
  ([HA & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}})).
- **The Broker persists messages in memory only.** No persistent volume is used.
  Brief interruptions are bridged by NATS topic persistence, but a far-behind or
  unavailable Server consumer lets Broker memory climb
  ([sizing]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}})).
- **The Local Provider has no authentication.** Never run shared production on
  it; preselect a Remote Provider
  ([identity]({{< ref "installation/production/authentication-and-identity.md" >}})).
- **WebSocket support is mandatory at the ingress.** Without it the UI loads but
  never updates live
  ([networking]({{< ref "installation/production/networking-and-connectivity.md" >}})).
- **OAuth requires the external callback URL.** Behind a proxy, an unset or wrong
  `MESHERY_SERVER_CALLBACK_URL` breaks login.
- **Out-of-cluster requires a reachable Broker endpoint.** Cross-cloud
  reachability of the Broker is the most common multi-cloud failure
  ([multi-cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}})).
- **Discovery cost scales with cluster size and churn.** Without MeshSync
  blacklisting, very large/churny clusters dominate Server memory and database
  growth.
- **Policy evaluation is time-boxed.** Large designs may hit
  `POLICY_EVAL_TIMEOUT` (default `3m`); tune CPU and the timeout together.
- **The chart ships permissive defaults for portability.** No resource
  requests/limits, empty security contexts, and a `stable-latest`,
  `Always`-pull image are convenient for getting started but are **not**
  production settings—set them explicitly.

## The consolidated production-readiness checklist

Work through these before going live. Each group links to its source page.

### Deployment model

- [ ] Topology chosen deliberately (in-cluster vs. out-of-cluster; Docker vs.
      Kubernetes), with Kubernetes + Helm for production.
      [Deployment Models]({{< ref "installation/production/deployment-models.md" >}})
- [ ] One Operator/MeshSync/Broker per managed cluster (or embedded mode chosen
      deliberately).

### Sizing & performance

- [ ] Explicit CPU/memory **requests and limits** set for Server, Operator,
      MeshSync, and Broker.
      [Sizing]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}})
- [ ] Server memory headroom for the MeshSync snapshot; Broker memory headroom
      for bursts.
- [ ] MeshSync `informer_config` blacklist tuned for large clusters.
- [ ] Storage provisioned for the Server cache; autoscaling configured where
      appropriate.

### High availability & resiliency

- [ ] Liveness/readiness probes enabled and tuned.
      [HA & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}})
- [ ] Anti-affinity spreads replicas across nodes/zones.
- [ ] Helm values, env config, and connection definitions in version control
      (your real backup).
- [ ] Redeploy-from-config recovery drill validated.

### Networking & connectivity

- [ ] Ingress fronts Meshery with TLS; **WebSocket upgrades** verified.
      [Networking]({{< ref "installation/production/networking-and-connectivity.md" >}})
- [ ] `MESHERY_SERVER_CALLBACK_URL` set to the external URL.
- [ ] Broker endpoint reachable from the Server and restricted to its origin.
- [ ] Egress to the Remote Provider and registries allowed; network policies
      applied.
- [ ] Any fronting CDN/caching proxy honors origin cache headers (immutable
      assets, `no-cache` HTML) and never caches `/api/*` or the WebSocket.

### Security & identity

- [ ] **Remote Provider preselected** (`PROVIDER`, `PROVIDER_BASE_URLS`); Local
      Provider not used.
      [Identity]({{< ref "installation/production/authentication-and-identity.md" >}})
- [ ] Least-privilege RBAC; `rbac.nodes` only where required; per-cluster
      credentials.
- [ ] Hardened `securityContext` (non-root, dropped caps, read-only root FS with
      writable data volume) on Server and Operator.
      [Security Hardening]({{< ref "installation/production/security-hardening.md" >}})
- [ ] Secrets (kubeconfig, provider, pull) sourced from Secrets/external manager.
- [ ] Images pinned to immutable tags from a trusted/mirrored registry.

### Multi-cluster & multi-cloud

- [ ] One least-privilege kubeconfig context per cluster.
      [Multi-Cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}})
- [ ] Per-cloud node-watch RBAC and Broker endpoint reachability validated.
- [ ] Private connectivity for private clusters; cross-region latency assessed.

### Monitoring & observability

- [ ] External uptime check on `/healthz/ready`; workload metrics scraped.
      [Monitoring & KPIs]({{< ref "installation/production/monitoring-observability-and-kpis.md" >}})
- [ ] Dedicated alerts on Server memory and Broker memory; per-cluster
      connection health and provider reachability monitored.
- [ ] OpenTelemetry tracing to a real collector; logs centralized at
      `LOG_LEVEL=4`.
- [ ] `mesheryctl system check` scheduled as a synthetic check.

### Operations

- [ ] Upgrade procedure uses `values-upgrade.yaml`; `stable` release channel.
      [Upgrade Guide]({{< ref "installation/upgrades/index.md" >}})
- [ ] Runbooks reference the
      [troubleshooting guides]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}).
- [ ] Ownership, on-call, and escalation defined for the Meshery deployment.

## Where to go next

- Revisit any area above via its linked page in this set.
- For component internals, see the
  [Meshery Architecture]({{< ref "concepts/architecture/_index.md" >}}) reference.
- For runtime configuration, see the
  [Server Environment Variables]({{< ref "installation/advanced/environment-variables.md" >}})
  reference.

{{< related-discussions tag="meshery" >}}
