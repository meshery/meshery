---
title: Envoy WASM Filter Management
aliases: 
- /tasks/filter-management
- /configuration-management/filter-management
categories: [infrastructure]
description: Import, inspect, publish, and remove Envoy WebAssembly (WASM) filters with Meshery UI and mesheryctl.
---

Meshery stores Envoy WebAssembly (WASM) filters alongside your designs, so you can keep the filter binaries your data plane uses in one place, share them, and remove them when you are done. Use this guide to import a filter, inspect it, publish it to the Catalog, and delete it.

## Before you begin

- A running Meshery Server, and `mesheryctl` signed in to it if you want to use the CLI. See [Authenticate with Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}).
- A compiled filter: a `.wasm` file on your machine or a URL that serves one. The community maintains sample filters in [layer5io/wasm-filters](https://github.com/layer5io/wasm-filters).

## Import a filter

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
1. Open `/configuration/filters` on your Meshery Server (for example `http://localhost:9081/configuration/filters`).
2. Click **Import Filters**.
3. Enter a **Name** for the filter.
4. Choose an **Upload method**:
   - **File Upload**: select the `.wasm` file on your machine.
   - **URL Upload**: paste the URL that serves the file.
5. Click **Import**.
{{% /tab %}}
{{% tab header="mesheryctl" %}}
Import from a local file or a URL:

```bash
mesheryctl filter import /path/to/filter.wasm
mesheryctl filter import https://example.com/myfilter.wasm
```

Set a name, and pass the filter's configuration either as a file path or as an inline string:

```bash
mesheryctl filter import /path/to/filter.wasm --name "rate-limit" --wasm-config ./config.json
```
{{% /tab %}}
{{< /tabpane >}}

## Find and inspect filters

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
The Filters page lists every filter you own, in grid or table view. Filter the list by **Visibility** (Public, Private, and, when the Catalog is enabled, Published). Each filter offers these actions:

- **Config**: edit the filter's configuration.
- **Download**: download the filter binary.
- **Filter Information**: view the filter's details.
- **Clone**: copy a published filter into your own filters.
{{% /tab %}}
{{% tab header="mesheryctl" %}}
List filters, 25 per page, optionally searching by name:

```bash
mesheryctl filter list
mesheryctl filter list "rate"
mesheryctl filter list --page 2
mesheryctl filter list --count
```

View one filter by name, ID, or a unique prefix of either, or view them all:

```bash
mesheryctl filter view "rate-limit"
mesheryctl filter view --all --output-format json --save
```

`--output-format` accepts `yaml` (the default) or `json`. `--save` writes the output to a file.
{{% /tab %}}
{{< /tabpane >}}

## Publish a filter to the Catalog

If your provider enables the [Meshery Catalog](https://meshery.io/catalog) and you have permission to publish WASM filters, a **Publish** action appears on each of your filters. Click it, fill in the catalog details, and submit. Published filters go through the provider's approval flow before others can see and clone them. Use **Unpublish** to withdraw one.

## Delete a filter

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
Select one or more filters in table view and click the delete icon, then confirm **Delete**.
{{% /tab %}}
{{% tab header="mesheryctl" %}}
Delete by name, ID, or a unique prefix of either. If the prefix matches more than one filter, Meshery deletes the first match, so use the full name or ID when in doubt.

```bash
mesheryctl filter delete "rate-limit"
```
{{% /tab %}}
{{< /tabpane >}}

## Frequently Asked Questions

### Q: Why does Meshery provide and manage WASM filters?

Meshery manages WASM filters for data plane proxies for:

- <strong>Agility</strong>: in partnership with service mesh control planes, Meshery can dynamically load, configure, and unload filter(s) in running Envoy proxies.
- <strong>Maintainability</strong>: To extend the functionality you don't have to change the entire codebase.
- <strong>Diversity</strong>: WASM filters can be easily compiled with many popular programming languages like C/C++ and Rust.
- <strong>Reliability and isolation</strong>: WASM filters are deployed into VM and hence they are isolated from the Envoy process. For instance if the WASM filter crashes it will not impact the envoy process.
- <strong>Security</strong>: The WASM filter uses a well defined API to connect with the host i.e. Envoy proxy which can be modified only a limited number of connections.
- <strong>Speed</strong>: WASM runs at near-native speed.

### Q: What is Envoy proxy?

An <strong>Envoy</strong> is a high performance, programmable L3/L4 and L7 proxy that a number of service mesh use as their data plane. See the service mesh landscape for a complete list. Connection and request handling are at the core of Envoy's focus. Envoy has a pluggable filter module that allows the incorporation of WebAssembly filters that provide additional traffic filtering intelligence. Often this includes higher order functions such as access control, transformation, data enrichment, auditing and more.

### Q: What is WASM?

WebAssembly is a low-level assembly-like language with a compact binary format that runs with near-native performance and provides languages such as C/C++, C# and Rust with a compilation target so that they can run in many different environments, including behind Envoy's Application Binary Interface to extend Envoy with new capabilities. WASM programs are distributed in an efficient portable binary instruction format and embed into isolated execution environments.

### Q: What are the disadvantages of WASM filters to be taken into consideration?

Execution of WASM filters requires system resources. Just as executing the same functionality that your filter provides outside of the WASM virtual machine, resources can include overhead of latency, CPU, and memory. See the [Cloud Native Performance](https://smp-spec.io) project for more detail.

### Q: What is the Envoy Proxy WASM SDK?

Envoy proxy runs the WASM filter in a stack-based virtual machine so that the filter memory is isolated from the host environment. The interaction between the embedded host i.e. Envoy proxy, and the WASM filter are done via call backs and functions by the Envoy Proxy WASM SDK. The Envoy Proxy WASM SDK is implemented in various programming languages like C++, RUST, AssemblyScript and Go (still experimental)

### Q: What is Image Hub?

Image Hub is a sample application which was first officially demonstrated at DockerCon 2020. This sample application is written to run on consul to explore the WebAssembly (WASM) modules in the Envoy filter. This application is written in the RUST programming language. These modules are used to implement the multi-tenancy or to implement the per user rate limiting in the applications endpoint.

