---
layout: default
title: mesheryctl system channel set
permalink: reference/mesheryctl/system/channel/set
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

# mesheryctl system channel set

## Description

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.set.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system channel set [flags]
  </div>
</pre>

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.set.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.set.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>
