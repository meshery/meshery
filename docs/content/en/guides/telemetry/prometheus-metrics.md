---
title: Prometheus Metrics
description: >-
  Explore Prometheus metrics, compose PromQL queries with live preview, and save
  your own panels — all from a registered Prometheus connection.
categories: [performance]
weight: 20
---

The **Telemetry → Metrics** view is a Prometheus-native explorer. Where Grafana gives you ready-made dashboards to browse, Prometheus has none — so this view lets you do what Grafana cannot: discover metrics, write [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/), preview the result live, and save the panels you care about into your own grid.

## Prerequisites

- A registered **Prometheus** [Connection]({{< ref "concepts/logical/connections/index.md" >}}) in the **Connected** or **Registered** state. Add a [Credential]({{< ref "concepts/logical/credentials.md" >}}) if your Prometheus requires authentication. See [Telemetry prerequisites]({{< ref "guides/telemetry/_index.md" >}}).

## Open the Metrics view

1. In the Meshery navigation, expand **Telemetry** and select **Metrics**.
2. Use the **connection picker** to choose the Prometheus connection you want to query.
3. Check the **connection status** indicator — it pings Prometheus and reports reachability and the detected server version.

## Explore metrics and build a query

Open the **metric explorer** to compose a query:

1. **Search metrics** — type to filter the list of available metric names served by your Prometheus.
2. **Inspect metadata** — select a metric to see its type (counter, gauge, histogram, …), help text, and unit.
3. **Browse labels** — list a metric's label names and their values to build precise matchers.
4. **Write PromQL** — compose an expression in the editor. A **live preview** chart shows the result over the current time range as you refine it.

Grafana-style time macros (such as `$__rate_interval` and `$__interval`) are expanded automatically based on the selected time window and step, so expressions copied from Grafana panels work here too.

## Save panels

When a query is showing what you want:

1. Give the panel a **title**.
2. Choose a **visualization type** — time series, stat, gauge, or bar.
3. Optionally set a **unit** for axis/value formatting.
4. **Save** the panel. It is added to your metrics grid.

Saved panels can be **edited** (title, type, unit) or **deleted** at any time. They are stored on the Prometheus connection itself, so your grid persists and travels with the connection.

## Read the grid

- The grid renders all saved panels, fetching their data in a single batched request per refresh.
- The shared **time range** picker and **refresh** controls apply to the whole grid.
- A failing panel shows its error inline without affecting the others.

## Related

- [Grafana Dashboards]({{< ref "guides/telemetry/grafana-dashboards.md" >}}) — for browsing and rendering existing Grafana dashboards.
- [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}})
