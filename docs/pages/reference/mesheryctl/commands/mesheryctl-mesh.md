---
layout: default
title: mesheryctl mesh
permalink: reference/mesheryctl/commands/mesheryctl-mesh
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
{% for command_hash in site.data.mesheryctlcommands.meshes.commands %}{% assign command = command_hash[1] %}
    {{ command.description }}
{% endfor %}
<br/>

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl mesh [flags] 
    </div>
  </div>
</pre>

<!-- Options/Flags available in this command -->
## Options

{% for subcommand_hash in site.data.mesheryctlcommands.meshes.validate.commands %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ subcommand.name }}  # {{ subcommand.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for flag_hash in site.data.mesheryctlcommands.meshes.validate.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.name }}  # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

<!-- All possible example use cases of the command -->
## Examples


<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
    mesheryctl mesh validate --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test] --namespace [namespace to be used]
    </div>
  </div>
 </pre>

<br/>

<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
    mesheryctl mesh validate --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test]
    </div>
  </div>
 </pre>
