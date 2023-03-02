---
layout: default
title: Filter Management
permalink: tasks/filter-management
type: tasks
language: en
list: include
---

## Frequently Asked Questions
### Q: What is Envoy proxy?

An <strong>Envoy</strong> is a high performance, programmable L3/L4 and L7 proxy that a number of service mesh use as their data plane. See the [service mesh landscape](https://layer5.io/service-mesh-landscape) for a complete list. Connection and request handling are at the core of Envoy’s focus. Envoy's has a pluggable filter module that allows the incorporation of WebAssembly filters that provide additional traffic filtering intelligence. Often this includes higher order functions such as access control, transformation, data enrichment, auditing and more.

### Q: What is WASM ?
WebAssembly is a low-level assembly-like language with a compact binary format that runs with near-native performance and provides languages such as C/C++, C# and Rust with a compilation target so that they can run in many different environments, including behind Envoy's Application Binary Interface to extend Envoy with new capabilities. WASM programs are distributed in an efficient portable binary instruction format and embed into isolated execution environments.

<--
### Q: Why does Meshery provide and manage WASM filters?

Meshery manages WASM filters for data plane proxies for:

-  <strong>Agility</strong>: in partnership with service mesh control planes, Meshery can dynamically load, configure, and unload filter(s) in running Envoy proxies.
- <strong>Maintainability</strong>: To extend the functionality we don’t have to change the entire codebase.
- <strong>Diversity</strong> : WASM filters can be easily compiled with many popular programming languages like C/C++ and Rust.
- <strong>Reliability and isolation</strong> : WASM filters are deployed into VM and hence they are isolated from the Envoy process. For instance if the WASM filter crashes it will not impact the envoy process.
- <strong>Security</strong> : The WASM filter uses a well defined API to connect with the host i.e envoy proxy which can be modified only a limited number of connections.
- <strong>Speed</strong> : It is originally used to speed up large web application

-->

### Q: What are the disadvantages of WASM filters to be taken into consideration?

Execution of WASM filters requires system resources. Just as executing the same functionality that your filter provides outside of the WASM virtual machine, resources can include overhead of latency, CPU, and memory. See the [Service Mesh Performance](https://smp-spec.io) project for more detail.

### Q: What is the Envoy Proxy WASM SDK?

Envoy proxy runs the WASM filter in a stack-based virtual machine so that the filter memory is isolated from the host environment. The interaction between the embedded host i.e. Envoy proxy and the WASM filter are done via call backs and functions by the Envoy Proxy WASM SDK. The Envoy Proxy WASM SDK is implemented in various programming languages like C++, RUST, AssemblyScript and GO(still experimental)
 
### Q: What is Image Hub ?

Image Hub is a sample application which was first officially demonstrated at DockerCon 2020. This sample application is written to run on consul to explore the  WebAssembly (WASM) modules in the Envoy filter. This application is written in the RUST programming language. These modules are used to implement the multi-tenancy or to implement the per user rate limiting in the applications endpoint.
