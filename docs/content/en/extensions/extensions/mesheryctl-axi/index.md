---
title: mesheryctl-axi
description: An agent-ergonomic AXI wrapper around mesheryctl - token-efficient TOON output, always non-interactive, with help[] next-step suggestions.
display_title: false
aliases:
- /extensions/mesheryctl-axi
---

# 🤖 mesheryctl-axi

`mesheryctl-axi` is an [AXI](https://axi.md/)-compliant wrapper around [`mesheryctl`]({{< ref "reference/references/mesheryctl/_index.md" >}}), built for AI coding agents. Where `mesheryctl` is optimized for a human at a terminal, `mesheryctl-axi` is optimized for an agent consuming the output as context: it wraps the human CLI rather than replacing it, so the two stay in step.

## Features

1. **Token-efficient TOON reporting:** `list` and `view` results are rendered in [TOON](https://toonformat.dev/) rather than verbose tables or JSON, reducing the context an agent spends reading them.
2. **Always non-interactive:** no TTY prompts, so a command never wedges an unattended agent session.
3. **Structured errors:** unknown flags and failures exit non-zero with a structured TOON error rather than free-form text.
4. **`help[]` next-step suggestions:** successful commands append the commands an agent is most likely to want next.
5. **Definitive empty states:** an empty result is reported explicitly (for example, `connections: 0`) instead of as silence.
6. **Schema-faithful content:** design and model *content* is returned as raw YAML or JSON - never re-encoded as TOON.

## Installation and Use

### Prerequisites

- **Node.js** v22 or newer.
- **`mesheryctl` installed and authenticated.** `mesheryctl-axi` spawns `mesheryctl`; it does not embed Meshery. See [Installing mesheryctl]({{< ref "installation/mesheryctl/_index.md" >}}). Point at a specific binary with `MESHERYCTL_BIN=/path/to/mesheryctl`.

### Usage

Run it directly with `npx` - no global install required:

<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">npx -y mesheryctl-axi</div>
   </div>
</pre>

With no arguments, `mesheryctl-axi` prints a content-first home: what it is, the resolved `mesheryctl` path, and a best-effort system status and context.

Resource listings and views are reported in TOON:

<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">npx -y mesheryctl-axi connection list</div>
   </div>
</pre>

<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">npx -y mesheryctl-axi system status</div>
   </div>
</pre>

Design and model *content* is returned verbatim as YAML or JSON:

<pre class="codeblock-pre">
  <div class="codeblock">
     <div class="clipboardjs">npx -y mesheryctl-axi design content &lt;name&gt; --format yaml</div>
   </div>
</pre>

Source, issues, and the full command reference: [github.com/meshery-extensions/mesheryctl-axi](https://github.com/meshery-extensions/mesheryctl-axi).

<hr />

## About Meshery Extensions

[Meshery Extensions](https://meshery.io/extensions) are plugins or add-ons that enhance the functionality of the Meshery platform beyond its core capabilities. Meshery supports different types of extensions ([docs]({{< ref "extensions/_index.md" >}})):

- [Adapters]({{< ref "concepts/architecture/adapters.md" >}}): Adapters allow Meshery to interface with the different cloud native infrastructure.
- [Load Generators]({{< ref "reference/extensibility/load-generators.md" >}}): for performance characterization and benchmarking
- [Integrations]({{< ref "extensions/models/_index.md" >}}): model-based support for a broad variety of design and orchestration of cloud and cloud native platforms, tools, and technologies.
- [Providers]({{< ref "reference/extensibility/providers/index.md" >}}): for connecting to different cloud providers and infrastructure platforms
- [UI Plugins]({{< ref "reference/extensibility/ui.md" >}}): Meshery UI has a number of extension points that allow users to customize their experience with third-party plugins.
