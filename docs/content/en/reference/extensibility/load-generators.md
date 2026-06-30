---
title: "Extensibility: Load Generators"
description: Meshery offers support for more types of Kubernetes and cloud native infrastructure than any other project or product in the world. Meshery uses adapters for more deeply managing specific types of infrastructure.
categories: [extensibility]
aliases:
- /extensibility/load-generators
---

Users may prefer to use one load generator over the next given the difference of capabilities between load generators, so Meshery provides a `load generator interface` (a gRPC interface) behind which a load generator can be implemented. Meshery provides users with choice of which load generator they prefer to use for a given performance test. Users may configure their own preference for the load generator, different from the default load generator.

### What function do load generators in Meshery provide? 

Load generators will provide the capability to run load tests from Meshery. As of today the load generator is embedded as a library in Meshery and Meshery invokes the load generator APIs with the right load test options to run the load test. At the moment, Meshery supports HTTP load generation. Support for gRPC and TCP load generation is on the roadmap. Meshery has functional integration with fortio.

### Which load generators does Meshery support?

1. [fortio](https://github.com/fortio/fortio) - Fortio load testing library, command line tool, advanced echo server and web UI in go (golang). Allows to specify a set query-per-second load and record latency histograms and other useful stats.
