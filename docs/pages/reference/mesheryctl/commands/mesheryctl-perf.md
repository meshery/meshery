---
layout: default
title: mesheryctl perf
permalink: reference/mesheryctl/perf
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl perf

## Description

<!-- Description of the command. Preferably a paragraph -->
{% for command_hash in site.data.mesheryctlcommands.performance.commands %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl perf [flags] 
  </div>
</pre>

## Examples

{% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>

## Options & Flags

{% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.name }}
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