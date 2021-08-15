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

{% for subcommand_hash in site.data.mesheryctlcommands.meshes.deploy.subcommand %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh deploy
  </div>
</pre>

## Options inherited from parent commands
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>
