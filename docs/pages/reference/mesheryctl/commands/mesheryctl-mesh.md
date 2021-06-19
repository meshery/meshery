---
layout: default
title: mesheryctl mesh
permalink: reference/mesheryctl/mesh
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl mesh

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for command_hash in site.data.mesheryctlcommands.meshes.commands %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh [flags] 
  </div>
</pre>

<!-- All possible example use cases of the command -->
## Examples


<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl mesh validate --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test] --namespace [namespace to be used]

    mesheryctl mesh validate --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test]
  </div>
 </pre>
<br/>


<!-- Options/Flags available in this command -->
## Options & Flags

{% for subcommand_hash in site.data.mesheryctlcommands.meshes.validate.subcommand %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ subcommand.name }}
  </div>
</pre>
{% endfor %}
<br/>

## Options inherited from parent commands
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>