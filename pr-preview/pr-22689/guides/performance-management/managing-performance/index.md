# Managing Performance

> Organize load tests into performance profiles, choose and configure the load generator, and read infrastructure metrics alongside your results.

---

LLMS index: [llms.txt](/pr-preview/pr-22689/llms.txt)

---

Measuring and managing performance is key to operating any infrastructure efficiently. Meshery organizes that work around performance profiles. For step-by-step instructions on running a test, see [Performance Management with Meshery](/pr-preview/pr-22689/guides/performance-management/performance-management/).

## Performance Profiles

A performance profile saves a load test's configuration so you can run the same test again and track how performance changes over time. With profiles you can:

1. Save a test configuration and rerun it from Meshery UI or with `mesheryctl perf apply <profile-name>`.
2. Keep every result of a test in one place, under its profile.
3. Rename a profile that started as an ad hoc test so you can reuse it.

## Load Generators

Meshery runs tests with [fortio](https://github.com/fortio/fortio), a fast, small, embeddable Go load testing library and command-line tool. Meshery's load generator interface is [extensible](/pr-preview/pr-22689/reference/extensibility/) to support others.

## Configuring a test

Each test or profile sets:

- The endpoint URL to test.
- Duration.
- Concurrent requests.
- Queries per second.
- Request headers, cookies, body, and content type.
- A CA certificate for HTTPS endpoints.
- Additional options passed straight to the load generator, as JSON.

See [Cloud Native Performance](https://smp-spec.io) for the full description of the test types that Meshery supports.

## Comparing results

Meshery runs statistical analysis when each test finishes and presents latency as a histogram with percentile buckets, alongside throughput. Select two or more results of a profile and choose **Compare selected** to compare their latency and throughput. See [Interpreting Performance Test Results](/pr-preview/pr-22689/guides/performance-management/interpreting-performance-test-results/).

## Infrastructure metrics alongside results

To understand what a test costs your cluster, view your Grafana dashboards and Prometheus metrics in Meshery while the test runs. Register Grafana or Prometheus as a [Connection](/pr-preview/pr-22689/guides/infrastructure-management/registering-a-connection/), then open **Telemetry**. See [Grafana Dashboards](/pr-preview/pr-22689/guides/telemetry/grafana-dashboards/) and [Prometheus Metrics](/pr-preview/pr-22689/guides/telemetry/prometheus-metrics/).
