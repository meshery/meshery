---
layout: default
title: mesheryctl pattern
permalink: reference/mesheryctl/commands/subcommands/mesheryctl-pattern
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
# mesheryctl pattern

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for command_hash in site.data.mesheryctlcommands.pattern.commands %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl pattern [flags] 
    </div>
  </div>
</pre>

<!-- Options/Flags available in this command -->
## Options & Flags

{% for flag_hash in site.data.mesheryctlcommands.pattern.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.name }} # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for subcommand_hash in site.data.mesheryctlcommands.pattern.subcommands %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ subcommand.name }} # {{ subcommand.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.pattern.subcommands %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
    {{ subcommand.usage }}
    </div>
  </div>
 </pre>
 <br/>
{% endfor %}