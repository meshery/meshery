# Using Metrics in Meshery

> How to view Prometheus and Grafana metrics in Meshery

Source: /pr-preview/pr-21670/guides/performance-management/meshery-metrics/

## View metrics in Meshery

Meshery integrates with **Grafana** and **Prometheus** so you can visualize the health and resource usage of your cloud native infrastructure without leaving Meshery.

Metrics visualization lives in the dedicated **[Telemetry](/pr-preview/pr-21670/guides/telemetry/)** section of Meshery. Rather than entering a Grafana or Prometheus URL and API key into a settings form, you register Grafana and Prometheus as Meshery [Connections](/pr-preview/pr-21670/concepts/logical/connections/) (each with an associated [Credential](/pr-preview/pr-21670/concepts/logical/credentials/)), and Telemetry queries them on your behalf — proxying every request securely through the Meshery server.

There are two ways to view metrics:

- **[Grafana Dashboards](/pr-preview/pr-21670/guides/telemetry/grafana-dashboards/)** (Telemetry → Charts) — browse, pin, and render your existing Grafana dashboards.
- **[Prometheus Metrics](/pr-preview/pr-21670/guides/telemetry/prometheus-metrics/)** (Telemetry → Metrics) — explore metrics, compose PromQL, preview the result, and save your own panels.

---

## Getting started

1. Create a [Credential](/pr-preview/pr-21670/concepts/logical/credentials/) for your Grafana API key / service-account token, or your Prometheus token (if authentication is required).
2. [Register a connection](/pr-preview/pr-21670/guides/infrastructure-management/registering-a-connection/) of kind **Grafana** or **Prometheus**, pointing at the instance URL and selecting the credential.
3. Open **Telemetry** in the Meshery navigation, pick your connection, and start viewing dashboards or metrics.

For details on authentication, shared controls, and how rendering works, see the [Telemetry guide](/pr-preview/pr-21670/guides/telemetry/).
