---
layout: default
title: mesheryctl system context delete
permalink: reference/mesheryctl/system/context/delete
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system context delete

## Description 

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.delete.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}


<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system context delete [context-name] [flags]
  </div>
</pre> 

<!-- All possible example use cases of the command -->
## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.delete.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.delete.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>

<!-- Options/Flags available in this command -->
## Options & Flags

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.delete.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.flag }}
  </div>
</pre>
{% endfor %}
<br/>