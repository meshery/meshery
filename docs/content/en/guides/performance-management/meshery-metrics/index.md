---
title: Using Metrics in Meshery
description: How to view Prometheus and Grafana metrics in Meshery
categories: [performance]
---

## View metrics in Meshery

Meshery integrates with **Grafana** and **Prometheus** so you can visualize the health and resource usage of your cloud native infrastructure without leaving Meshery.

Metrics visualization lives in the dedicated **[Telemetry]({{< ref "guides/telemetry/_index.md" >}})** section of Meshery. Rather than entering a Grafana or Prometheus URL and API key into a settings form, you register Grafana and Prometheus as Meshery [Connections]({{< ref "concepts/logical/connections/index.md" >}}) (each with an associated [Credential]({{< ref "concepts/logical/credentials.md" >}})), and Telemetry queries them on your behalf — proxying every request securely through the Meshery server.

There are two ways to view metrics:

- **[Grafana Dashboards]({{< ref "guides/telemetry/grafana-dashboards.md" >}})** (Telemetry → Charts) — browse, pin, and render your existing Grafana dashboards.
- **[Prometheus Metrics]({{< ref "guides/telemetry/prometheus-metrics.md" >}})** (Telemetry → Metrics) — explore metrics, compose PromQL, preview the result, and save your own panels.

---

## Getting started

1. Create a [Credential]({{< ref "concepts/logical/credentials.md" >}}) for your Grafana API key / service-account token, or your Prometheus token (if authentication is required).
2. [Register a connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) of kind **Grafana** or **Prometheus**, pointing at the instance URL and selecting the credential.
3. Open **Telemetry** in the Meshery navigation, pick your connection, and start viewing dashboards or metrics.

For details on authentication, shared controls, and how rendering works, see the [Telemetry guide]({{< ref "guides/telemetry/_index.md" >}}).
