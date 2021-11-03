---
layout: default
title: mesheryctl system context
permalink: /v0.4/reference/mesheryctl/system/context
type: reference
display-title: "false"
language: en
command: system-context
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system context

<!-- Description of the command. Preferably a paragraph -->
## Description 

{% assign name = site.data.mesheryctlcommands.cmds[page.command] %}
{{ name.description }}


<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.usage }}
  </div>
</pre> 

<!-- All possible example use cases of the command -->
## Examples

{% for subcommand_hash in name.subcommands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}

<!-- Options/Flags available in this command -->
<!-- ## Options & Flags

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.subcommands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ subcommand.name }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.arg }}
  </div>
</pre>
{% endfor %}
<br/>
-->
## Options inherited from parent commands
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>