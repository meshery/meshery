---
layout: integration
title: Jaeger
subtitle: Distributed Tracing Integration with Meshery
permalink: /extensibility/integrations/jaeger
docURL: https://docs.meshery.io/extensibility/integrations/jaeger
description: Learn how Meshery integrates with Jaeger, the CNCF distributed tracing system, to enable observability and deep insights into service-to-service communication.
integrations-category: Observability
integrations-subcategory: Tracing
registrant: meshery
---

# Jaeger — Distributed Tracing Integration (Meshery)

# Overview

Jaeger is a CNCF project used for distributed tracing and monitoring of microservices-based applications. Meshery integrates with Jaeger to provide end-to-end observability, enabling users to trace requests across services managed within Meshery.

# Setup and Requirements

- Meshery running on Kubernetes.
- Jaeger deployed and accessible within the cluster.
- Meshery CLI or UI access for configuration.
- Service mesh enabled for advanced tracing features.

# Using Meshery with Jaeger

This integration allows you to:

- Visualize distributed traces for applications deployed via Meshery.
- Correlate Meshery Designs with service-level performance metrics.
- Identify bottlenecks and latency issues across microservices.

# Meshery Catalog and Design Patterns

Sample tracing configuration pattern:

```yaml
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: sample-tracing
spec:
  tracing:
    providers:
      - name: jaeger
    randomSamplingPercentage: 100.0
```

These patterns help standardize observability practices across environments.

# Tutorials and Resources

- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Meshery Docs - Jaeger Integration](https://docs.meshery.io/extensibility/integrations/jaeger)
- [Meshery Playground](https://play.meshery.io/)

# Related Integrations

- [Prometheus Integration](/extensibility/integrations/prometheus)
- [Grafana Integration](/extensibility/integrations/grafana)

# What’s Next

- Contribute additional observability design patterns to [Meshery Catalog](https://meshery.io/catalog).
- Extend Meshery Learning Paths to cover observability practices with Jaeger.
