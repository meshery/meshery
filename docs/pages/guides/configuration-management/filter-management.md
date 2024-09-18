---
layout: default
title: Envoy WASM Filter Management
permalink: guides/configuration-management/filter-management
redirect_from: /tasks/filter-management
type: guides
category: configuration
language: en
list: include
abstract: "Meshery provides the ability to manage WASM filters for data plane proxies for agility, maintainability, diversity, reliability and isolation, security, and speed."
---

Meshery offers powerful features that simplify the management of Envoy WASM filters in Istio.

### Meshery UI and CLI Integration

- **Import WASM Envoy Filters**: Easily import your WebAssembly Envoy filters into Meshery using the intuitive UI or the command-line interface ([CLI](https://docs.meshery.io/reference/mesheryctl#data-plane-intelligence)).
- **Publish and Clone Filters**: Share your filters with the community by publishing them in the [Meshery Catalog](https://meshery.io/catalog) and make it effortless for others to clone them.
- **Download WASM Binaries**: Access and download WebAssembly binaries directly from your Remote Provider like [Layer5 Cloud](https://meshery.layer5.io/).
- **Efficient Data Plane Design**: Seamlessly design and deploy Istio and Envoy data planes using extensions like [Kanvas](https://layer5.io/cloud-native-management/kanvas).
- **Contribute to Open Source**: Get involved with the cloud native community by contributing to any of the open source [wasm-filters](https://github.com/layer5io/wasm-filters) developed by the Meshery community.

<img src="https://mcusercontent.com/6b50be5aea3dfe1fd4c041d80/images/a07ef29a-4cf8-986e-9bd3-78db7dc00ce1.png" width="100px" style="float:left;position:relative;margin: 1rem 1rem 1rem 0rem" />

## Envoy WASM Filter Management

As of Meshery v0.7, the management of Envoy WASM filters includes:

**Meshery Server Integration**

- ✅ New WASM Filter Component: A new component is available in the Design Configurator to simplify filter management.
- ✅ Improved Component Icon: Enhancements have been made to provide a more user-friendly experience.
- ✅ Relationship Patch Policy: A new policy has been introduced to streamline the relationship between components.

**Meshery CLI**

- ✅ Import Filters: Use `mesheryctl filter import [URL | filepath]` to effortlessly import filters.
- ✅ Delete Filters: Remove unwanted filters with `mesheryctl filter delete [filter-name | ID]`.
- ✅ View Filters: Gain insights into your filters using `mesheryctl filter view [filter-name | ID]`.
- ✅ List and Search Filters: Easily navigate filters with `mesheryctl filter list [filter-name | ID]`.

**Kanvas Integration**

- ✅ Apply WASM Filters: Hierarchically apply filters to your Envoy configurations.
- ✅ Deploy/Undeploy Filters: Streamline the deployment and removal of filters with custom configurations.
- ✅ Enhanced Visibility: View filters within a dedicated panel and seamlessly drop them onto the canvas.

<img src="https://mcusercontent.com/6b50be5aea3dfe1fd4c041d80/images/1e9c2e71-1b3e-a132-4766-8cefdc9861d2.png" width="50px" style="float:left;position:relative;margin: 1rem 1rem 1rem 0rem" />

**Meshery UI Enhancements**

- ✅ Import Filters: Import filters directly via URL or filesystem.
- ✅ Improved Browsing: Browse, search, and view filters in both grid view and table view.
- ✅ Streamlined Lists: Easily list and search filters in grid view or table view.
- ✅ Download Filters: Download filters effortlessly from the UI.

**Remote Provider and Meshery Catalog Integration**

- ✅ Permanent Storage and Artifact Dispersal: Ensure your filters are securely stored and widely accessible.
- ✅ Catalog Features: Import, clone, download, publish, and unpublish filters within the Meshery Catalog.
- ✅ Enhanced Privacy Controls: Manage user permissions, team ownership, and visibility for your filters.
- ✅ Content Curation: Streamline content curation with an approval flow request queue.

## Performance Management Upgrades

In addition to Envoy WASM filter management, we've also introduced new performance management features:

**Meshery Server**

- ✅ Performance Profiles: Tailor your performance profiles to match your specific needs.
- ✅ SSL Certificate Support: Benefit from SSL certificate support for Fortio in the server.
- ✅ Performance Analysis Comparison: Compare performance with and without filters.
- ✅ GetNighthawk and Cloud Native Performance Releases: Access the latest releases of GetNighthawk and Cloud Native Performance.

**Meshery CLI**

- ✅ Performance Profile Flags: Specify additional load generator flags with performance profiles.

**Meshery UI**

- ✅ Performance Profile Flags: Customize load generator flags directly in the UI.

**Service Performance Project**

- ✅ Define Performance Profiles: Clearly define your performance profiles.
- ✅ Dashboard Integration: Display test results on the dashboard.
- ✅ Intel Integration: Incorporate Intel design into scheduled workflows for Istio.
- ✅ Consolidated Performance Profiles: Streamline performance profiles on the dashboard.

Meshery offers WebAssembly-based network traffic management filters in Envoy, bringing a new level of flexibility and control to your Istio deployments.

## Frequently Asked Questions

### Q: Why does Meshery provide and manage WASM filters?

Meshery manages WASM filters for data plane proxies for:

- <strong>Agility</strong>: in partnership with service mesh control planes, Meshery can dynamically load, configure, and unload filter(s) in running Envoy proxies.
- <strong>Maintainability</strong>: To extend the functionality we don’t have to change the entire codebase.
- <strong>Diversity</strong> : WASM filters can be easily compiled with many popular programming languages like C/C++ and Rust.
- <strong>Reliability and isolation</strong> : WASM filters are deployed into VM and hence they are isolated from the Envoy process. For instance if the WASM filter crashes it will not impact the envoy process.
- <strong>Security</strong> : The WASM filter uses a well defined API to connect with the host i.e envoy proxy which can be modified only a limited number of connections.
- <strong>Speed</strong> : It is originally used to speed up large web application

### Q: What is Envoy proxy?

An <strong>Envoy</strong> is a high performance, programmable L3/L4 and L7 proxy that a number of service mesh use as their data plane. See the [service mesh landscape](https://layer5.io/service-mesh-landscape) for a complete list. Connection and request handling are at the core of Envoy’s focus. Envoy's has a pluggable filter module that allows the incorporation of WebAssembly filters that provide additional traffic filtering intelligence. Often this includes higher order functions such as access control, transformation, data enrichment, auditing and more.

### Q: What is WASM ?

WebAssembly is a low-level assembly-like language with a compact binary format that runs with near-native performance and provides languages such as C/C++, C# and Rust with a compilation target so that they can run in many different environments, including behind Envoy's Application Binary Interface to extend Envoy with new capabilities. WASM programs are distributed in an efficient portable binary instruction format and embed into isolated execution environments.

### Q: What are the disadvantages of WASM filters to be taken into consideration?

Execution of WASM filters requires system resources. Just as executing the same functionality that your filter provides outside of the WASM virtual machine, resources can include overhead of latency, CPU, and memory. See the [Cloud Native Performance](https://smp-spec.io) project for more detail.

### Q: What is the Envoy Proxy WASM SDK?

Envoy proxy runs the WASM filter in a stack-based virtual machine so that the filter memory is isolated from the host environment. The interaction between the embedded host i.e. Envoy proxy and the WASM filter are done via call backs and functions by the Envoy Proxy WASM SDK. The Envoy Proxy WASM SDK is implemented in various programming languages like C++, RUST, AssemblyScript and GO(still experimental)

### Q: What is Image Hub ?

Image Hub is a sample application which was first officially demonstrated at DockerCon 2020. This sample application is written to run on consul to explore the WebAssembly (WASM) modules in the Envoy filter. This application is written in the RUST programming language. These modules are used to implement the multi-tenancy or to implement the per user rate limiting in the applications endpoint.

