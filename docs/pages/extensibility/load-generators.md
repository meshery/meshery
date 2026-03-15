---
layout: default
title: "Extensibility: Load Generators"
permalink: extensibility/load-generators
type: Extensibility
#redirect_from: architecture/adapters
abstract: "Meshery offers support for more types of Kubernetes and cloud native infrastructure than any other project or product in the world. Meshery uses adapters for more deeply managing specific types of infrastructure."
language: en
list: include
---

Users may prefer to use one load generator over the next given the difference of capabilities between load generators, so Meshery provides a `load generator interface` (a gRPC interface) behind which a load generator can be implemented. Meshery provides users with choice of which load generator they prefer to use for a given performance test. Users may configure their own preference for the load generator, different from the default load generator.

### What function do load generators in Meshery provide? 

Load generators will provide the capability to run load tests from Meshery. As of today the load generators are embedded as libraries in Meshery and Meshery invokes the load generators APIs with the right load test options to run the load test. At the moment, Meshery supports HTTP load generation. Support for gRPC and TCP load generation is on the roadmap. Meshery has functional integration with fortio, wrk2, and nighthawk.

### Why support multiple load generators?

Different use cases and different opinions call for different approaches to statistical analysis of the performance results. For example, wrk2 accounts for a concept called Coordinated Omission.

### Which load generators does Meshery support?

1. [fortio](https://github.com/fortio/fortio) - Fortio load testing library, command line tool, advanced echo server and web UI in go (golang). Allows to specify a set query-per-second load and record latency histograms and other useful stats.
1. [wrk2](https://github.com/giltene/wrk2) - A constant throughput, correct latency recording variant of wrk.
1. [nighthawk](https://getnighthawk.dev/) - Enables users to run distributed performance tests to better mimic real-world, distributed systems scenarios.
  - See the Nighthawk project.
