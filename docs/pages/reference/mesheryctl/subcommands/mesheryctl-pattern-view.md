---
layout: default
title: mesheryctl pattern view
permalink: reference/mesheryctl/commands/subcommands/mesheryctl-pattern-view
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
# mesheryctl pattern view

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for subcommand_hash in site.data.mesheryctlcommands.pattern.view.commands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl pattern view [flags] 
  </div>
</pre>

## Examples

<pre class="codeblock-pre">
  <div class="codeblock">
  {% for subcommand_hash in site.data.mesheryctlcommands.pattern.view.commands %}{% assign subcommand = subcommand_hash[1] %}
  # {{ subcommand.description }}
  {{ subcommand.usage }}
  {% endfor %}
  {% for flag_hash in site.data.mesheryctlcommands.pattern.view.flags %}{% assign flag = flag_hash[1] %}
  # {{ flag.description }}
  {{ flag.usage }}
  {% endfor %}
  </div>
 </pre>
 <br/>

<!-- Options/Flags available in this command -->
## Options & Flags

<pre class="codeblock-pre">
  <div class="codeblock">
  {% for flag_hash in site.data.mesheryctlcommands.pattern.view.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.name }} # {{ flag.description }}
  {% endfor %}
  </div>
</pre>
<br/>