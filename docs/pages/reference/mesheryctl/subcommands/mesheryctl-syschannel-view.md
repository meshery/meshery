---
layout: default
title: mesheryctl system channel view
permalink: reference/mesheryctl/commands/subcommands/mesheryctl-system-channel-view
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
# mesheryctl system channel view

## Description

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl system channel view [flags]
    </div>
  </div>
</pre>

## Examples

<pre class="codeblock-pre">
  <div class="codeblock">
  {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.view.command %}{% assign subcommand = subcommand_hash[1] %}
  # {{ subcommand.description }}
  {{ subcommand.usage }}
  {% endfor %}
  {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.view.flag %}{% assign flag = flag_hash[1] %}
  # {{ flag.description }}
  {{ flag.usage }}
  {% endfor %}
  </div>
</pre>
<br/>


<!-- Options/Flags available in this command -->
## Options

<pre class="codeblock-pre">
  <div class="codeblock">
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.view.flag %}{% assign flag = flag_hash[1] %}
    {{ flag.flag }} # {{ flag.description }}
    {% endfor %}
  </div>
</pre>
<br/>
