---
layout: default
title: mesheryctl system channel switch
permalink: reference/mesheryctl/system/channel/switch
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

# mesheryctl system channel switch

## Description

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.switch.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system channel switch [flags]
  </div>
</pre>

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.switch.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.switch.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>
