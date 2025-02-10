---
layout: default
title: Helm Kanvas Snapshot
permalink: extensions/helm-kanvas-snapshot
language: en
abstract: Helm CLI plugin to visually render a snapshot of your Helm charts.
display-title: "false"
list: include
type: extensions
category: kanvas
---

# <img style="height: 4rem; width: 4rem;" src="{{site.baseurl}}/assets/img/integrations/helm-controller/components/helm-chart/icons/color/helm-chart-color.svg" /> Helm Kanvas Snapshot

## Extension Overview

The **Kanvas Snapshot Helm Plugin** allows users to generate a visual snapshot of their Helm charts directly from the command line. It simplifies the process of creating Meshery Snapshots, providing a visual representation of packaged Helm charts. This plugin integrates with Meshery Cloud and GitHub Actions to automate the workflow of snapshot creation, which is especially useful for Helm users who need to quickly visualize their chart configurations.

Helm charts can be complex, especially when custom configurations are applied via `values.yaml` files. This Meshery extension bridges the gap between Helm chart configurations and their visual representation by converting Helm charts into **Kanvas Snapshots**. These snapshots can be received either via email or as a URL displayed directly in the terminal.

### Features

1. **Snapshot Generation:** Create visual snapshots of Helm charts, complete with associated resources.
2. **Synchronous/Asynchronous Delivery:** Choose between receiving snapshots via email or directly in the terminal.
3. **Seamless Integration:** Leverages Meshery Cloud and GitHub Actions to handle snapshot rendering.
4. **Support for Packaged Charts:** Works with both packaged `.tar.gz` charts and unpackaged Helm charts.

## Installation and Use

To install the Meshery Snapshot Helm Plugin, use the following steps:

### Prerequisites

- `helm` must be [installed]( helm plugin install https://github.com/meshery/helm-kanvas-snapshot) on your system.
- (Optional) A free [Layer5 Cloud](https://cloud.layer5.io) user account.

**Plugin Installation**

1. Run the following command to install the Helm Kanvas Snapshot plugin:

   ```bash
   helm plugin install https://github.com/meshery/helm-kanvas-snapshot
   ```

3. Verify the installation by running:

   ```bash
   helm plugin list
   ```

   You should see the Kanvas Snapshot listed as `snapshot`.

4. Set up the required environment variables (see the [Environment Variables](#environment-variables) section).

### Usage

Once the plugin is installed, you can generate a snapshot using either a packaged or unpackaged Helm chart.

```bash
helm snapshot -f <chart-URI> [--name <snapshot-name>] [-e <email>]
```

- **`-f`**, **`--file`**: (required) path or URL to the Helm chart (required).
- **`--name`**: (optional) name for the snapshot. If not provided, a name will be auto-generated based on the chart name.
- **`-e`**, **`--email`**: (optional) email address to notify when snapshot is ready. If not provided, a link to the snapshot will be displayed in the terminal.

**Example**

To generate a snapshot for a Helm chart located at `https://meshery.io/charts/v0.8.0-meshery.tar.gz`, you can use:

```bash
helm snapshot -f https://meshery.io/charts/v0.8.0-meshery.tar.gz --name meshery-chart
```

<hr />

## About Meshery Extensions

[Meshery Extensions](https://meshery.io/extension) are plugins or add-ons that enhance the functionality of the Meshery platform beyond its core capabilities. Meshery supports different types of extensions ([docs](https://docs.meshery.io/extensions/)):

- [Adapters](https://docs.meshery.io/concepts/architecture/adapters): Adapters allow Meshery to interface with the different cloud native infrastructure.
- [Load Generators](https://docs.meshery.io/extensibility/load-generators): for performance characterization and benchmarking
- [Integrations](https://docs.meshery.io/extensibility/integrations): model-based support for a broad variety of design and orchestration of cloud and cloud native platforms, tools, and technologies.
- [Providers](https://docs.meshery.io/extensibility/providers): for connecting to different cloud providers and infrastructure platforms
- [UI Plugins](https://docs.meshery.io/extensibility/ui): Meshery UI has a number of extension points that allow users to customize their experience with third-party plugins.