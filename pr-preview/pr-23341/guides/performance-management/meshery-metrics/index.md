# Using Metrics in Meshery

> How to view Prometheus and Grafana metrics in Meshery

---

LLMS index: [llms.txt](/docs-preview/pr-preview/pr-23341/llms.txt)

---

## View metrics in Meshery

Meshery integrates with **Grafana** and **Prometheus** so you can visualize the health and resource usage of your cloud native infrastructure without leaving Meshery.

Metrics visualization lives in the dedicated **[Telemetry](/docs-preview/pr-preview/pr-23341/guides/telemetry/)** section of Meshery. Rather than entering a Grafana or Prometheus URL and API key into a settings form, you register Grafana and Prometheus as Meshery [Connections](/docs-preview/pr-preview/pr-23341/concepts/logical/connections/) (each with an associated [Credential](/docs-preview/pr-preview/pr-23341/concepts/logical/credentials/)), and Telemetry queries them on your behalf - proxying every request securely through the Meshery server.

There are two ways to view metrics:

- **[Grafana Dashboards](/docs-preview/pr-preview/pr-23341/guides/telemetry/grafana-dashboards/)** (Telemetry → Charts) - browse, pin, and render your existing Grafana dashboards.
- **[Prometheus Metrics](/docs-preview/pr-preview/pr-23341/guides/telemetry/prometheus-metrics/)** (Telemetry → Metrics) - explore metrics, compose PromQL, preview the result, and save your own panels.

---

## Getting started

1. Create a [Credential](/docs-preview/pr-preview/pr-23341/concepts/logical/credentials/) for your Grafana API key / service-account token, or your Prometheus token (if authentication is required).
2. [Register a connection](/docs-preview/pr-preview/pr-23341/guides/infrastructure-management/registering-a-connection/) of kind **Grafana** or **Prometheus**, pointing at the instance URL and selecting the credential.
3. Open **Telemetry** in the Meshery navigation, pick your connection, and start viewing dashboards or metrics.

For details on authentication, shared controls, and how rendering works, see the [Telemetry guide](/docs-preview/pr-preview/pr-23341/guides/telemetry/).
