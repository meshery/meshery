---
layout: page
title: Performance Management
permalink: functionality/performance-management
type: functionality
---

# Performance Management

## Load Generators
Meshery provides users with choice of which load generator they prefer to use for a given performance test. Users may set their configure their own preference of load generator different that the default load generator. 

## Node and Service Mesh Metrics

 Meshery provides performance test results alongside environment metrics, including service mesh control and data plane metrics as well as cluster node resource metrics, so that operators may easily understand the overhead of their service mesh's control plane and data plane in context of the overhead incurred on nodes within the cluster.

## Grafana and Meshery

Connect Meshery to your existing Grafana instance and Meshery will import the boards of your choosing. 

<a href="/docs/assets/img/performance-management/meshery-and-grafana.png">
    <img src="/docs/assets/img/performance-management/meshery-and-grafana.png" style="width: 100%" />
</a>

### Connecting to Grafana
If you have an API key configured to restrict access to your Grafana boards, you will need to enter the API key when establishing Meshery's connection to Grafana.

* Importing Grafana boards
    - Importing existing Grafana boards via API
    - Importing custom Grafana board via yaml
* Configuring graph panel preferences

## Prometheus and Meshery
Meshery allows users to connect to one or more Prometheus instances in order to gather telemetric data (in the form of metrics). These metrics may pertain to service meshes, Kubernetes, applications on the mesh or really... any metric that Prometheus has collected.

Once you have connected Meshery to your Prometheus deployment(s), you may perform ad-hoc connectivity tests to verify communication between Meshery and Prometheus.

## Suggested Reading

- Guide: [Interpreting Performance Test Results](guides/interpreting-performance-test-results)