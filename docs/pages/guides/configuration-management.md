---
layout: default
title: Configuring your Cloud Native Infrastructure and Applications with Meshery
description: Learn how to use Meshery Designs effectively
permalink: guides/configuration-management
type: Guides
language: en
---

As an cloud native management platform, Meshery has built-in support infastructure and application configuration management. Use Meshery to configure your multi-cluster Kubernetes environments on-premises or across clouds.

## What is a Meshery Design?

Meshery's core construct of configuration management is that of a design. A design is ultimately a document – and a document in which you capture and describe your desired state of improv the environment.

## Creating a Meshery Design

 you have your choice of how are you would like to create an infrastructure design. You can use the design configurator inside of measure, UI or extensions, like MeshMap, or you can hand code your design, using the mesh model spec, if you like. 
 
 Use Meshery UI or mesheryctl to manage [designs](../tasks/patterns.md).

## Pattern Management Through Meshery UI

Meshery also comes with seed patterns when users start Meshery for the first time. These patterns cover common use cases and examples for you as you explore Meshery.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-patterns.png" />

Users can also import these patterns to their remote provider from this [sample repository](https://github.com/service-mesh-patterns/service-mesh-patterns/tree/master/samples).

<img src="{{ site.baseurl }}/assets/img/configuration-management/pattern-import.png" />

Once these patterns are imported, you can then edit these patterns or use the pattern configurator to configure them according to your requirements.

<img src="{{ site.baseurl }}/assets/img/configuration-management/pattern-configure-button.png" />

<img src="{{ site.baseurl }}/assets/img/configuration-management/pattern-configure.png" />

## Pattern Management Through Meshery CLI

You can also manage service mesh patterns through Meshery's CLI, mesheryctl.

The `mesheryctl pattern` subcommand lets you import and apply patterns to your cluster.

For example, if you have your pattern written in a file say, `istio-bookinfo.yaml` which deploys Istio service mesh and onboards the BookInfo app on Istio, you can use mesheryctl to apply this pattern as shown below:

```
mesheryctl pattern apply -f istio-bookinfo.yaml
```

If you already have a pattern imported into Meshery, you can apply the pattern by name.

```
mesheryctl pattern apply BookInfoApp
```

This will apply the pattern BookInfoApp, which has already been imported into Meshery.

See [mesheryctl Command Reference](../reference/mesheryctl/#cloud-native-pattern-configuration-and-management) for more details on the `pattern` subcommand.

## WASM Filters

Meshery can be used for managing WebAssembly Filters through the UI or the CLI.

### Filter Management Through Meshery UI

Like patterns, Meshery also comes with some sample WebAssembly Filters for you to experiment.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-filters.png" />

You can also import these filters manually to your provider from the [wasm-filters](https://github.com/layer5io/wasm-filters) repo.

Meshery's sample application [ImageHub](./deploying-sample-apps.md#imagehub) will let you test out configuring these filters out-of-the-box.

You can onboard ImageHub to an installed service mesh as shown below.

<img src="{{ site.baseurl }}/assets/img/configuration-management/image-hub.png" />

### Filter Management Through Meshery CLI

You can also manage WASM filters through Meshery's CLI, mesheryctl.

The `mesheryctl exp filter` subcommand lets you import and configure filters in your service meshes.

For example,

```
mesheryctl exp filter apply -f metrics_collector_bg.wasm
```

If you already have a filter imported into Meshery, you can configure the filter by name.

```
mesheryctl exp filter apply metrics_collector_bg
```

## Applications

Meshery can also manage your Kubernetes applications and can onboard/offboard them of deployed service meshes.

### Managing Applications Through Meshery UI

Meshery has a set of [sample applications](./deploying-sample-apps.md) which you can use to quickly test out your deployment.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-applications-seeded.png" />

You can also bring in your own applications by uploading it from filesystem or importing it from a URL.

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-applications.png" />

The application files can be described in Kubernetes manifest format or pattern format.

You can also run the same application on different service meshes and compare the performance. See [Performance Management with Meshery](performance-management.md) for more details.

### Managing Applications Through Meshery CLI

The `mesheryctl app` subcommand lets you manage your custom application workloads with Meshery.

You can onboard/offboard applications from your mesh as shown in the example below.

```
mesheryctl app onboard imagehub.yaml
```

```
mesheryctl app offboard imagehub.yaml
```
