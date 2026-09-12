---
title: Grafana Dashboards
description: >-
  Browse, pin, and render your existing Grafana dashboards inside Meshery using a
  registered Grafana connection.
categories: [performance]
weight: 10
---

The **Telemetry → Charts** view lets you find your existing Grafana dashboards and render their panels directly inside Meshery, using data proxied securely through the Meshery server. Nothing is re-built or re-modeled — Meshery reads the same dashboards you have in Grafana and draws their panels natively.

## Prerequisites

- A registered **Grafana** [Connection]({{< ref "concepts/logical/connections/index.md" >}}) in the **Connected** or **Registered** state, with a [Credential]({{< ref "concepts/logical/credentials.md" >}}) holding a Grafana API key or service-account token. See [Telemetry prerequisites]({{< ref "guides/telemetry/_index.md" >}}).

## Open the Charts view

1. In the Meshery navigation, expand **Telemetry** and select **Charts**.
2. Use the **connection picker** to choose the Grafana connection you want to view.
3. Check the **connection status** indicator — it pings Grafana and reports whether it is reachable, along with the detected Grafana version. Clicking a Grafana connection's status chip on the [Connections]({{< ref "concepts/logical/connections/index.md" >}}) page also pings it and raises a notification with the result.

## Browse and pin dashboards

1. Open the **dashboard library** (the board browser) and search for a dashboard by name.
2. Select a dashboard to render it. Meshery loads the dashboard definition, resolves its template variables and datasources, and draws each panel.
3. Use **Add** (pin) to keep a dashboard handy. Pinned dashboards appear in your Charts view every time you return.
4. Use **Remove** (unpin) to take a dashboard out of the pinned list.

Pinned dashboards are stored on the Grafana connection itself, so they persist and are available wherever that connection is used.

## Render and read panels

- Panels are laid out following the dashboard's own grid, and rendered as native time-series charts (and single-stat values) inside Meshery.
- A shared **time range** picker (last 5 minutes up to last 7 days) applies to every panel; use **Refresh** to reload on demand or set an auto-refresh interval.
- **Template variables** (for example `$datasource`, `$job`, `$instance`) and Grafana time macros (such as `$__rate_interval`) are resolved automatically before each query is sent, so dashboards that rely on variables render correctly.
- If a single panel's query fails, the error is shown inline on that panel only — the rest of the dashboard keeps rendering.

{{% alert title="Rendering model" color="info" %}}
Meshery renders panels by querying Grafana's datasource proxy server-side and drawing the results itself. It does not embed Grafana's own UI, so you do not need browser-side access to Grafana — only the Meshery server needs to reach it.
{{% /alert %}}

## Related

- [Prometheus Metrics]({{< ref "guides/telemetry/prometheus-metrics.md" >}}) — for exploring raw Prometheus metrics and saving your own PromQL panels.
- [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}})
