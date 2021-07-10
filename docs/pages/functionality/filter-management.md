---
layout: default
title: Filter Management
permalink: functionality/filter-management
type: functionality
language: en
list: include
---

### What is an Envoy ?

An <strong>Envoy</strong> is a high performance, programmable L3/L4 and L7 proxy that many mesh services implement like istio. There are network filters at the core of the Envoy’s connection and the traffic handling which when mixed with the filter chains allows the implementation of the higher-order functionalities for access control, transformation, data enrichment, auditing and many more. 

### What is a WASM ?

<strong>WASM (<i>Web Assembly</i>)</strong> is an extremely interesting way which greatly simplifies the process of extending Envoy with new capabilities. WASM is an efficient portable binary instruction format providing an embeddable and isolated execution environment.

### Why are WASM filters used ?

The Wasm filters are used for the following reasons :

-  <strong>Agility</strong> : We can dynamically load into the running envoy without stopping or re-compiling it.
- <strong>Maintainability</strong> : To extend the functionality we don’t have to change the entire codebase.
- <strong>Diversity</strong> : WASM filters can be easily compiled with many popular programming languages like C/C++ and Rust.
- <strong>Reliability and isolation</strong> : WASM filters are deployed into VM and hence they are isolated from the Envoy process. For instance if the WASM filter crashes it will not impact the envoy process.
- <strong>Security</strong> : The WASM filter uses a well defined API to connect with the host i.e envoy proxy which can be modified only a limited number of connections.
- <strong>Speed</strong> : It is originally used to speed up large web application

### What are the disadvantages of WASM filters to be taken into consideration ?

The disadvantages of the WASM filters are :

- The performance is nearly 70% as the native C++.
- It takes a higher memory usage due to the requirement of starting one or more WASM virtual machines.

### What are Envoy Proxy WASM SDK ?

The envoy proxy runs the WASM filter in a stack-based virtual machine so that the filter memory is isolated from the host environment. The interaction between the embedded host i.e. the envoy proxy and the WASM filter are done via call backs and functions by the Envoy Proxy WASM SDK. The Envoy Proxy WASM SDK is implemented in various programming languages like C++, RUST, AssemblyScript and GO(still experimental)
 
### What is an Image Hub ?

Image Hub is a sample application which was first officially demonstrated at DockerCon 2020. This sample application is written to run on consul to explore the  WebAssembly (WASM) modules in the Envoy filter. This application is written in the RUST programming language. These modules are used to implement the multi-tenancy or to implement the per user rate limiting in the applications endpoint.
