---
layout: default
title: mesheryctl system
permalink: reference/mesheryctl/system
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.commands %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system [flags] 
  </div>
</pre>

<!-- All possible example use cases of the command -->
## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.subcommands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>

<!-- Options/Flags available in this command -->
## Options & Flags

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.subcommands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.name }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.arg }}
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