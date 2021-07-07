---
layout: default
title: mesheryctl system status
permalink: reference/mesheryctl/system/status
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system status

## Description

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.status.command %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system status 
  </div>
</pre>

<!--## Examples

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.status.command %}{% assign command = command_hash[1] %}
{{ command.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ command.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.status.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>
-->

<!-- Options/Flags available in this command -->
<!-- ## Options

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.status.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.flag }}
  </div>
</pre>
{% endfor %}
<br/>
-->
