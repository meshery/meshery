---
layout: default
title: Configuring your Deployment with Patterns, Filters and Applications
description: Learn how to configure your deployments and service meshes
permalink: guides/configuration-management
type: Guides
language: en
---

Meshery now comes with built in support for service mesh patterns, WASM filters and application workloads. Users can use Meshery to manage these artifacts and use them with their service mesh deployment.

## Service Mesh Patterns

Users can use Meshery UI or mesheryctl to manage [service mesh patterns](../tasks/patterns.md).

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

See [mesheryctl Command Reference](../reference/mesheryctl) for more details on the `pattern` subcommand.

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
