---
layout: default
title: mesheryctl perf
permalink: reference/mesheryctl/commands/mesheryctl-perf
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
# mesheryctl perf

## Description
<!-- Description of the command. Preferably a paragraph -->
{% for command_hash in site.data.mesheryctlcommands.performance.commands %}{% assign command = command_hash[1] %}
    {{ command.description }}
{% endfor %}
<br/>

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl perf [flags] 
    </div>
  </div>
</pre>

## Options

{% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.name }}  # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

## Examples

{% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}