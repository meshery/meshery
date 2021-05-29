---
layout: default
title: mesheryctl mesh | Meshery
permalink: reference/mesheryctl/mesh
type: reference
display-title: "false"
language: en
lang: en
categories: en
list: exclude
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

<pre class="codeblock-pre">
  <div class="codeblock">
    {% for subcommand_hash in site.data.mesheryctlcommands.meshes.validate.commands %}{% assign subcommand = subcommand_hash[1] %}
    {{ subcommand.name }}  # {{ subcommand.description }}
    {% endfor %}
    {% for flag_hash in site.data.mesheryctlcommands.meshes.validate.flags %}{% assign flag = flag_hash[1] %}
    {{ flag.name }}  # {{ flag.description }}
    {% endfor %}
  </div>
</pre>
<br/>

## Options inherited from parent commands
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>