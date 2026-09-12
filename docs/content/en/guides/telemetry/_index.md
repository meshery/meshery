---
title: Telemetry
description: >-
  Visualize observability data in Meshery by browsing Grafana dashboards and
  exploring Prometheus metrics directly from your registered connections.
categories: [performance]
weight: 16
---

The **Telemetry** section of Meshery lets you view your observability data without leaving Meshery. It connects to the Grafana and Prometheus instances you have already registered as [Connections]({{< ref "concepts/logical/connections/index.md" >}}) and renders their data natively inside the Meshery UI.

Telemetry is organized into two modes, available as items under the **Telemetry** entry in the Meshery navigation:

- **Charts** — Browse, pin, and render your existing **Grafana** dashboards. See [Grafana Dashboards]({{< ref "guides/telemetry/grafana-dashboards.md" >}}).
- **Metrics** — Explore **Prometheus** metrics, compose PromQL queries, and save your own panels. See [Prometheus Metrics]({{< ref "guides/telemetry/prometheus-metrics.md" >}}).

## How it works

Telemetry is **connection-driven**. Instead of pasting a URL and API key into a settings form, you register Grafana and Prometheus as Meshery Connections (with an associated [Credential]({{< ref "concepts/logical/credentials.md" >}})), and Telemetry uses them. This means:

- Credentials are stored once, securely, and reused — Meshery proxies every query server-side, so your browser never needs direct access to Grafana or Prometheus.
- The same connection can be shared across an [Environment]({{< ref "concepts/logical/environments.md" >}}) and reused by other Meshery features.
- Your pinned dashboards and saved metric panels are stored on the connection itself, so they persist and travel with it.

{{% alert title="Authentication" color="info" %}}
Telemetry derives the authentication scheme from the connection's credential secret automatically: a value containing a colon (e.g. `user:password`) is treated as HTTP Basic auth, any other non-empty value is sent as a Bearer token (Grafana API key or Prometheus bearer token), and an empty secret is treated as anonymous access.
{{% /alert %}}

## Prerequisites

Before using Telemetry, register at least one Grafana or Prometheus connection:

1. Create a [Credential]({{< ref "concepts/logical/credentials.md" >}}) holding your Grafana API key / service-account token, or your Prometheus token (if your Prometheus requires authentication).
2. [Register the connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) of kind **Grafana** or **Prometheus**, pointing at the instance URL and selecting the credential.

Once a connection reaches the **Connected** or **Registered** state, it becomes selectable in the Telemetry connection picker.

## Shared controls

Both Telemetry modes share a common set of controls:

- **Connection picker** — choose which registered Grafana/Prometheus connection to view.
- **Connection status** — a live indicator pings the connection and shows whether it is reachable, along with the detected server version.
- **Time range** — a relative time-range selector (from the last 5 minutes up to the last 7 days) applied to every panel.
- **Refresh** — refresh on demand, or set an auto-refresh interval.

All panels on a board or grid are loaded in a single batched request per refresh, keeping the view responsive even with many panels.
