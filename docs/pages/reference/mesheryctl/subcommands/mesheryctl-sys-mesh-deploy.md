---
layout: default
title: mesheryctl system mesh deploy
permalink: reference/mesheryctl/mesh/deploy
type: reference
display-title: "false"
language: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl mesh deploy

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for command_hash in site.data.mesheryctlcommands.meshes.deploy.command %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh deploy [flags]
  </div>
</pre>

## Examples

{% for flag_hash in site.data.mesheryctlcommands.meshes.deploy.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.usage }}
  </div>
</pre>
{% endfor %}

<!-- Options/Flags available in this command -->

## Options & Flags

{% for flag_hash in site.data.mesheryctlcommands.meshes.deploy.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.name }}
  </div>
</pre>
{% endfor %}

## Options inherited from parent commands
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>
