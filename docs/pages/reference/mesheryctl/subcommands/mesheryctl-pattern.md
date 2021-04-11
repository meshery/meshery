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
  mesheryctl pattern [flags] 
  </div>
</pre>

## Examples

<pre class="codeblock-pre">
  <div class="codeblock">
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.subcommands %}{% assign subcommand = subcommand_hash[1] %}
  {{ subcommand.usage }}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.list.commands %}{% assign subcommand = subcommand_hash[1] %}
  {{ subcommand.usage }}
  {% endfor %}
  {% for flag_hash in site.data.mesheryctlcommands.pattern.list.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.usage }}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.view.commands %}{% assign subcommand = subcommand_hash[1] %}
  {{ subcommand.usage }}
  {% endfor %}
  {% for flag_hash in site.data.mesheryctlcommands.pattern.view.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.usage }}
  {% endfor %}
  </div>
 </pre>
 <br/>

<!-- Options/Flags available in this command -->
## Options & Flags

<pre class="codeblock-pre">
  <div class="codeblock">
  {% for flag_hash in site.data.mesheryctlcommands.pattern.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.name }} # {{ flag.description }}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.subcommands %}{% assign subcommand = subcommand_hash[1] %}
  {{ subcommand.name }} # {{ subcommand.description }}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.list.commands %}{% assign subcommand = subcommand_hash[1] %}
  {{ subcommand.name }} # {{ subcommand.description }}
  {% endfor %}
  {% for flag_hash in site.data.mesheryctlcommands.pattern.list.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.name }} # {{ flag.description }}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.view.commands %}{% assign subcommand = subcommand_hash[1] %}
  {{ subcommand.name }} # {{ subcommand.description }}
  {% endfor %}
  {% for flag_hash in site.data.mesheryctlcommands.pattern.view.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.name }} # {{ flag.description }}
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