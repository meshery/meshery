---
title: Monitoring, Observability & Health KPIs
linkTitle: Monitoring & KPIs
description: >-
  Health endpoints, the key performance indicators of Meshery's health, metrics,
  distributed tracing, centralized logging, alerting, and synthetic checks for
  observing Meshery in production.
categories: [installation]
weight: 80
aliases:
- /installation/production/monitoring
- /installation/production/observability
- /installation/production/kpis
---

A production management plane must itself be observed. This page defines what
"healthy" looks like for Meshery, the key performance indicators (KPIs) to track
per component, and how to wire Meshery into your monitoring, tracing, logging,
and alerting stack.

## Health endpoints

Meshery Server exposes Kubernetes-compliant health endpoints that are the
foundation of both self-healing and external monitoring:

| Endpoint | Used by | Signals |
| :--- | :--- | :--- |
| `/healthz/live` | Liveness probe | Server is running and responsive; provider capabilities are loaded. |
| `/healthz/ready` | Readiness probe | Server is ready to accept traffic. |

Query readiness with a verbose breakdown for detail:

```bash
kubectl exec --namespace meshery deployment/meshery -- \
  curl -s 'http://localhost:8080/healthz/ready?verbose=1'
```

```
[+]capabilities ok
[i]extension extension package found
healthz check passed
```

`[+]` passed, `[-]` failed (marks the pod unhealthy), `[i]` informational. These
endpoints are the right target for uptime/synthetic checks from outside the
cluster as well.

## The KPIs of Meshery's health

Track these indicators as the canonical signals of a healthy Meshery
deployment. Group them by component so dashboards and alerts map to where action
is taken.

### Meshery Server

| KPI | Why it matters | Healthy signal |
| :--- | :--- | :--- |
| Liveness/readiness status | Core availability | Both passing; readiness stable, not flapping |
| API/GraphQL latency & error rate | User-facing responsiveness | Low p95 latency; low 5xx/error rate |
| CPU utilization | Saturation under load/policy eval | Comfortably below limit |
| Memory utilization | Holds the MeshSync snapshot/registry; restart risk if exhausted | Comfortable headroom below limit |
| Restart count | Crash-loop / OOM detection | Stable; no recurring restarts |
| Database/cache size growth | Discovery scope and retention | Grows then plateaus; no unbounded climb |

### Meshery Operator, MeshSync, Broker (per cluster)

| KPI | Why it matters | Healthy signal |
| :--- | :--- | :--- |
| Operator pod running & reconciling | Manages Broker/MeshSync lifecycle | Running; no reconcile errors |
| MeshSync running & syncing | Cluster snapshot freshness | Running; snapshot updates as the cluster changes |
| Broker pod running | Eventing path up | StatefulSet pod ready |
| Broker memory | In-memory message backlog | Stable; not climbing (a climb means the Server consumer is behind) |
| Connection/chip status (per cluster) | End-to-end connectivity | Connected; Broker/Operator/MeshSync following the connection |

### Remote Provider

| KPI | Why it matters | Healthy signal |
| :--- | :--- | :--- |
| Provider reachability (egress) | Login & durable state depend on it | Reachable over HTTPS; auth succeeding |
| Auth success rate | User access | High success; no spikes in login failures |

{{% alert title="Watch Broker memory and Server memory together" color="info" %}}
Two early-warning signals deserve dedicated alerts: **Meshery Server memory**
(undersizing causes restarts that drop the cache) and **Broker memory** (a climb
indicates the Server consumer is falling behind). Together they catch most
capacity problems before they become outages. See
[Infrastructure, Sizing & Performance]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}}).
{{% /alert %}}

## Metrics with Prometheus and Grafana

Meshery integrates with Prometheus and Grafana, both for managing your
infrastructure's performance and for observing Meshery itself:

- Scrape Kubernetes workload metrics (CPU, memory, restarts) for the Meshery
  namespace to drive the Server/Operator/MeshSync/Broker KPIs above.
- The Broker exposes an **HTTP monitoring endpoint on `8222/tcp`** (NATS
  monitoring) for Broker-level visibility.
- Connect Prometheus and Grafana to Meshery to correlate management-plane health
  with the performance of the infrastructure under management. See the
  [performance management guides]({{< ref "guides/performance-management/meshery-metrics/index.md" >}}).

Build dashboards that put Server availability, resource saturation, per-cluster
connection health, and provider reachability on one pane.

## Distributed tracing with OpenTelemetry

Meshery Server supports OpenTelemetry tracing, configured via `OTEL_CONFIG`
(inline YAML). When unset, tracing is disabled. Use it in production to trace
request flows and diagnose latency:

```yaml
# Example OTEL_CONFIG (set as an env var / Helm env value)
service_name: meshery-server
service_version: 1.0.0
endpoint: otel-collector.observability:4317
insecure: false
```

Point the `endpoint` at your collector and keep `insecure: false` with proper
TLS in production. Avoid the insecure local-development settings on a real
deployment.

## Centralized logging

- **Aggregate stdout/stderr.** Meshery components log to standard streams; ship
  them to your centralized logging stack (e.g. Loki, Elasticsearch/OpenSearch,
  or a cloud logging service) for retention, search, and correlation.
- **Control verbosity with `LOG_LEVEL`.** Values are `0=panic`, `1=fatal`,
  `2=error`, `3=warn`, `4=info` (default), `5=debug`, `6=trace`. Run at `info`
  in production; raise to `debug`/`trace` temporarily when investigating.
  `DEBUG=true` forces debug-level logging.
- **Correlate across components.** Tie Server, Operator, MeshSync, and Broker
  logs together (by namespace/labels) so a discovery or connectivity issue can
  be traced across the path.

## Alerting

Turn the KPIs into actionable alerts. A solid baseline:

- **Availability:** readiness failing or pod not ready for N minutes; recurring
  restarts (crash loop / OOM).
- **Saturation:** Server CPU or memory sustained near limit; database/cache size
  growth anomalous.
- **Eventing:** Broker memory climbing; Broker or MeshSync pod not running;
  Operator reconcile errors.
- **Connectivity:** a cluster connection unhealthy; Broker endpoint unreachable
  from the Server.
- **Identity:** Remote Provider unreachable; spike in authentication failures.

Route these to the team that operates Meshery, and include the relevant runbook
links (troubleshooting, sizing, networking) in the alert.

## Synthetic and connectivity checks

- **`mesheryctl system check`** runs pre- and post-deployment health checks,
  including connectivity, and is well suited to scheduled synthetic validation.
  See the [reference]({{< ref "reference/references/mesheryctl/system/check.md" >}}).
- **External uptime checks** against `/healthz/ready` validate the full ingress
  → Server path (TLS, routing, readiness) from a user's perspective.
- **Per-connection checks** in the UI provide on-demand validation of each
  cluster's connectivity.

## Troubleshooting entry points

When a KPI trips, these guides are the fastest path to resolution:

- [Operator & MeshSync troubleshooting]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}})
  — Broker, MeshSync, and Operator issues.
- [Meshery Server troubleshooting]({{< ref "guides/troubleshooting/meshery-server.md" >}})
  — Server startup, provider, and API issues.

## Monitoring checklist

- [ ] Liveness/readiness probes enabled; external uptime check on
      `/healthz/ready`.
- [ ] Workload metrics scraped for the Meshery namespace; dashboards for
      Server/Operator/MeshSync/Broker KPIs.
- [ ] Dedicated alerts on Server memory and Broker memory.
- [ ] Per-cluster connection health monitored.
- [ ] Remote Provider reachability and auth success monitored.
- [ ] OpenTelemetry tracing configured to a real collector (`insecure: false`).
- [ ] Logs centralized; `LOG_LEVEL=4` (info) in steady state.
- [ ] `mesheryctl system check` scheduled as a synthetic check.

{{< related-discussions tag="meshery" >}}
