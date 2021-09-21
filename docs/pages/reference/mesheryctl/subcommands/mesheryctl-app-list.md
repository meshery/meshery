---
layout: default
title: mesheryctl app list
permalink: reference/mesheryctl/apps/list
type: reference
display-title: "false"
language: en
command: app
subcommand: list
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl app list

<!-- Description of the command. Preferably a paragraph -->
## Description and Example

<!-- Basic usage of the command -->
{% assign name = site.data.mesheryctlcommands.cmds[page.command].subcommands[page.subcommand] %}
{{ name.description }}

<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.usage }}
  </div>
</pre>
<br/>

<!-- Options/Flags available in this command -->
## Options & Flags

<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl app list -h # Shows help for the command
  </div>
</pre>
<br/>