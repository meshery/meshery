---
layout: default
title: "Extensibility: CLI Plugins"
permalink: extensibility/cli-plugins
type: Extensibility
abstract: "Meshery provides CLI plugins that extend existing tools like Helm and kubectl with snapshot and synchronization workflows that connect to Meshery and Kanvas."
language: en
list: include
---

Meshery offers CLI plugins as a point of extensibility. These plugins extend existing command line tools so that you can interact with Meshery and Kanvas directly from your terminal. CLI plugins are kept in separate repositories under the Meshery Extensions organization and are versioned on their own, which keeps Meshery's core binaries small and focused.

1. **Extensibility points offer clean separation of Meshery's core functionality versus plugin functionality.**
   - See a list of [Meshery's extensions](https://meshery.io/extensions).
1. **CLI plugins reuse familiar tools while connecting them to Meshery**
   - Helm and kubectl users can add Meshery specific workflows without changing their existing workflows.

### What functionality do CLI Plugins offer?

The functionality offered by a CLI plugin depends on the plugin and the host tool it extends, but all of them are focused on turning code or cluster state into Meshery designs and Kanvas snapshots, or synchronizing that state with Meshery Server.

Some examples include:

- **Visual snapshots from infrastructure as code**
  - Generate images from packaged Helm charts or collections of Kubernetes manifests.
- **Cluster state snapshots**
  - Capture a point in time view of Kubernetes resources and upload it to Meshery Server for visualization in Kanvas.
- **Automation and Git workflows**
  - Invoke snapshots from CI pipelines and GitHub Actions so reviewers can see what will be deployed before they merge.

## Types of CLI plugins

Meshery's CLI plugins are grouped by the tool they extend and the kind of snapshot they create.

- **Helm Kanvas Snapshot**
  - A Helm plugin that generates Kanvas snapshots from packaged or unpackaged Helm charts.
- **kubectl Kanvas Snapshot**
  - A kubectl plugin that generates Kanvas snapshots from Kubernetes manifest files.
- **kubectl MeshSync Snapshot**
  - A kubectl plugin that uses the MeshSync library to capture live cluster state and serialize it for upload to Meshery Server.

### Helm Kanvas Snapshot

The Helm Kanvas Snapshot plugin extends the `helm` CLI so you can create Meshery designs and Kanvas snapshots directly from Helm charts.

Typical usage:

```bash
helm helm-kanvas-snapshot -f <chart-URI> [--name <snapshot-name>] [-e <email>]
````

* `-f`, `--file` path: or URL to the Helm chart.
* `--name` optional: display name for the snapshot. If omitted, the name is derived from the chart.
* `-e`, `--email` optional: email address that receives the snapshot when it is ready.

The plugin sends the chart definition to the Kanvas Snapshot service. The service turns the chart into a Meshery design and renders a static image of the resulting topology. Snapshots can be delivered by email, added to your Meshery design catalog, or attached to pull requests through the Kanvas Snapshot GitHub Action.

Use Helm Kanvas Snapshot when:

* You already package your applications as Helm charts.
* You want a quick visual review of a chart before you deploy it.
* You want to include before and after diagrams in a CI workflow.

### kubectl Kanvas Snapshot

The kubectl Kanvas Snapshot plugin extends `kubectl` so that you can generate Kanvas snapshots from raw Kubernetes manifests instead of Helm charts.

The plugin reads one or more manifest files, sends them to the Kanvas Snapshot service, and receives a rendered image of the combined resources. This is helpful when:

* You manage workloads as plain YAML or Kustomize rather than Helm.
* You want to see how a set of manifests relate to each other as a single design.
* You want reviewers to see the shape of a change without parsing multiple YAML files.

Just like Helm Kanvas Snapshot, this plugin is implemented as a native kubectl plugin and follows the usual kubectl plugin conventions.

### kubectl MeshSync Snapshot

The kubectl MeshSync Snapshot plugin focuses on live clusters, not just manifests. It is a Krew plugin that uses the MeshSync library to discover resources in a Kubernetes cluster and capture an ad hoc snapshot of that state.

At a high level, the plugin:

1. Connects to your current Kubernetes context.
2. Uses MeshSync to walk the API and collect resource metadata.
3. Writes the snapshot out as YAML.
4. Gives you the option to upload that snapshot to Meshery Server for visualization in Kanvas.

Use kubectl MeshSync Snapshot when:

* You want to capture the current state of a cluster without installing the full Meshery stack.
* You need a shareable, read only representation of cluster state for debugging or review.
* You want to seed Meshery with an initial view of an existing environment.

## Design Principles: Meshery CLI Plugin Framework

Meshery's CLI plugin ecosystem follows a few guiding principles:

1. **Out of tree and independently versioned**

   * Plugins live in their own repositories under the Meshery Extensions organization and can release on their own cadence.
2. **Built on existing ecosystems**

   * Helm Kanvas Snapshot uses Helm's plugin system.
   * kubectl plugins are published through Krew and follow the standard kubectl plugin conventions.
3. **Terminal first experience**

   * Plugins are designed for users who prefer the terminal and want Meshery features available in the same place they already run `helm` and `kubectl`.
4. **Tight integration with Kanvas and Meshery Server**

   * All three plugins connect to Meshery based services so that snapshots and cluster state end up in the same visual designer and management plane.

## Building a CLI plugin

If you want to build a new CLI plugin that integrates with Meshery, you will typically:

1. Choose a host CLI

   * For example, implement a Helm plugin that shells out to Kanvas Snapshot, or a kubectl plugin that talks directly to Meshery Server.
2. Follow the host tool's plugin guidelines

   * Use Helm's plugin specification or kubectl Krew conventions for discovery, installation, and distribution.
3. Integrate with Meshery APIs

   * Call the Kanvas Snapshot service to render images, or use MeshSync libraries and Meshery Server APIs to push cluster data into Meshery.
4. Keep your code out of tree

   * Publish your plugin in its own repository and treat Meshery as a dependency instead of vendoring your code into the Meshery core repos.

## Managing your CLI plugin code

Like other Meshery extensions, CLI plugins are intended to be maintained as separate projects.

This has a couple of advantages:

1. Your plugin can be closed source or open source, independent of Meshery itself.
2. Bugs or breaking changes in plugin code do not affect Meshery's stability.
3. You can automate your own build, release, and distribution pipeline while still integrating cleanly with Meshery and Kanvas.

When Meshery or its Go libraries are upgraded, verify that your plugin still compiles and that its calls into Meshery services continue to behave as expected. For example, plugin authors often test against the same Go version and Kubernetes client library versions that are used in the Meshery release they target.

For more examples of Meshery CLI plugins, explore the Meshery Extensions organization and the Helm Kanvas Snapshot, kubectl Kanvas Snapshot, and kubectl MeshSync Snapshot repositories.
