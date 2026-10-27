# Helm Snapshot

> Helm CLI plugin to visually render a snapshot of your Helm charts.

---

LLMS index: [llms.txt](/pr-preview/pr-22614/llms.txt)

---

# <img style="height: 4rem; width: 4rem; vertical-align: middle;" src="images/helm-chart-color.svg" data-logo-for-dark="images/helm-chart-white.svg" data-logo-for-light="images/helm-chart-color.svg" class="logo-dark-light" alt="" /> Helm Snapshot

## Extension Overview

The **Helm Snapshot Plugin** allows users to generate a visual snapshot of their Helm charts directly from the command line. It simplifies the process of creating Meshery Snapshots, providing a visual representation of packaged Helm charts. This plugin integrates with Meshery Cloud and GitHub Actions to automate the workflow of snapshot creation, which is especially useful for Helm users who need to quickly visualize their chart configurations.

Helm charts can be complex, especially when custom configurations are applied via `values.yaml` files. This Meshery extension bridges the gap between Helm chart configurations and their visual representation by converting Helm charts into snapshots. These snapshots can be received either via email or as a URL displayed directly in the terminal.

### Features

1. **Snapshot Generation:** Create visual snapshots of Helm charts, complete with associated resources.
2. **Synchronous/Asynchronous Delivery:** Choose between receiving snapshots via email or directly in the terminal.
3. **Seamless Integration:** Leverages Meshery Cloud and GitHub Actions to handle snapshot rendering.
4. **Support for Packaged Charts:** Works with both packaged `.tar.gz` charts and unpackaged Helm charts.

## Installation and Use

To install the Meshery Snapshot Helm Plugin, use the following steps:

### Prerequisites

- `helm` must be installed on your system.
- (Optional) A free [Meshery Cloud](https://cloud.meshery.io) user account.

**Plugin Installation**

1. Run the following command to install the Helm Snapshot plugin:

<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs"> helm plugin install https://github.com/meshery-extensions/helm-kanvas-snapshot </div>
   </div>
</pre>

2. Verify the installation by running:


<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">helm plugin list</div>
   </div>
</pre>

   You should see the snapshot listed as `helm-kanvas-snapshot`.

No environment variables or tokens need to be set. The plugin installs a prebuilt release binary that already includes its Meshery Cloud configuration.

### Usage

Once the plugin is installed, you can generate a snapshot using either a packaged or unpackaged Helm chart.


<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">helm helm-kanvas-snapshot -f &lt;chart-URI&gt; [--name &lt;snapshot-name&gt;] [-e &lt;email&gt;]</div>
   </div>
</pre>

- **`-f`**, **`--file`**: (required) path or URL to the Helm chart.
- **`--name`**: (optional) name for the snapshot. If not provided, a name will be auto-generated based on the chart name.
- **`-e`**, **`--email`**: (optional) email address to notify when snapshot is ready. If not provided, a link to the snapshot will be displayed in the terminal.

**Example**

To generate a snapshot for a Helm chart located at `https://meshery.io/charts/meshery-v0.8.12.tgz`, you can use:

<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">helm helm-kanvas-snapshot -f https://meshery.io/charts/meshery-v0.8.12.tgz --name meshery-chart</div>
   </div>
</pre>

<hr />

## About Meshery Extensions

[Meshery Extensions](https://meshery.io/extensions) are plugins or add-ons that enhance the functionality of the Meshery platform beyond its core capabilities. Meshery supports different types of extensions ([docs](/pr-preview/pr-22614/extensions/)):

- [Adapters](/pr-preview/pr-22614/concepts/architecture/adapters/): Adapters allow Meshery to interface with the different cloud native infrastructure.
- [Load Generators](/pr-preview/pr-22614/reference/extensibility/load-generators/): for performance characterization and benchmarking
- [Integrations](/pr-preview/pr-22614/extensions/models/): model-based support for a broad variety of design and orchestration of cloud and cloud native platforms, tools, and technologies.
- [Providers](/pr-preview/pr-22614/reference/extensibility/providers/): for connecting to different cloud providers and infrastructure platforms
- [UI Plugins](/pr-preview/pr-22614/reference/extensibility/ui/): Meshery UI has a number of extension points that allow users to customize their experience with third-party plugins.
