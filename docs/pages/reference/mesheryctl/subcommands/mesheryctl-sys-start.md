---
layout: default
title: mesheryctl system start | Meshery
permalink: reference/mesheryctl/system/start
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
# mesheryctl system start

## Description

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.start.command %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system start [flags]
  </div>
</pre>

## Examples

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.start.command %}{% assign command = command_hash[1] %}
{{ command.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ command.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.start.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>


<!-- Options/Flags available in this command -->
## Options

<pre class="codeblock-pre">
  <div class="codeblock">
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.start.flag %}{% assign flag = flag_hash[1] %}
    {{ flag.flag }} # {{ flag.description }}
    {% endfor %}
  </div>
</pre>
<br/>
