---
layout: default
title: Configuring your Cloud Native Infrastructure and Applications
abstract: Learn how to use Meshery Designs effectively
permalink: guides/configuration-management/working-with-designs
redirect_from:
  - /tasks/patterns
category: configuration
type: guides
language: en
---

As a cloud native management platform, Meshery has built-in support for infrastructure and application configuration management. Use Meshery to configure your multi-cluster Kubernetes environments on-premises or across clouds.

## What is a Meshery Design?

Meshery's core construct of configuration management is that of a design. A design is ultimately a document that captures and describes the desired state of your environment.

[![Context Aware Design]({{ site.baseurl }}/assets/img/patterns/context-aware-design.svg
)]({{ site.baseurl }}/assets/img/patterns/context-aware-design.svg)

## Creating a Meshery Design

You can create an infrastructure design using the design configurator in the Meshery UI or extensions like Kanvas, or by hand-coding it with the mesh model spec.


Use Meshery UI or mesheryctl to create designs.

## Design Management Through Meshery UI

Meshery also comes with seed designs when users start Meshery for the first time. These designs cover common use cases and examples for you as you explore Meshery.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-patterns.png" width="50%" />

Users can also import these designs to their remote provider from this [sample repository](https://github.com/service-mesh-patterns/service-mesh-patterns/tree/master/samples).

<img src="{{ site.baseurl }}/assets/img/configuration-management/pattern-import.png" width="50%" />

Once these designs are imported, you can edit them or use the design configurator to configure them according to your requirements.

<img src="{{ site.baseurl }}/assets/img/configuration-management/pattern-configure-button.png" width="50%" />

<img src="{{ site.baseurl }}/assets/img/configuration-management/pattern-configure.png" width="50%" />

## Design Management Through Meshery CLI

You can also manage cloud native designs through Meshery's CLI, mesheryctl.

The `mesheryctl design` subcommand lets you import and apply designs to your cluster.

For example, if you have your design written in a file say, `istio-bookinfo.yaml` which deploys Istio service mesh and onboards the BookInfo app on Istio, you can use mesheryctl to apply this design as shown below:

```
mesheryctl design apply -f istio-bookinfo.yaml
```

If you already have a design imported into Meshery, you can apply the design by name.

```
mesheryctl design apply BookInfoApp
```

This will apply the design BookInfoApp, which has already been imported into Meshery.

See [mesheryctl design subcommand section]({{ site.baseurl }}/reference/mesheryctl/#cloud-native-design-configuration-and-management) for more details on the `design` subcommand.

## WASM Filters

Meshery can be used for managing WebAssembly Filters through the UI or the CLI.

### Filter Management Through Meshery UI

Like designs, Meshery also comes with some sample WebAssembly Filters for you to experiment.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-filters.png" width="50%" />

You can also import these filters manually to your provider from the [wasm-filters](https://github.com/layer5io/wasm-filters) repo.

Meshery's sample application [ImageHub]({{ site.baseurl }}/guides/infrastructure-management/sample-apps) will let you test out configuring these filters out-of-the-box.

You can onboard ImageHub to an installed service mesh as shown below.

<img src="{{ site.baseurl }}/assets/img/configuration-management/image-hub.png" width="50%" />

### Filter Management Through Meshery CLI

You can also manage WASM filters through Meshery's CLI, mesheryctl.

The `mesheryctl filter` command lets you import and configure WebAssembly filters.

For example,

```
mesheryctl exp filter apply -f metrics_collector_bg.wasm
```

If you already have a filter imported into Meshery, you can configure the filter by name.

```
mesheryctl exp filter apply metrics_collector_bg
```