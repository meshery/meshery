---
layout: default
title: mesheryctl system check
permalink: reference/mesheryctl/system/check
redirect_from: reference/mesheryctl/system/check/
type: reference
display-title: "false"
language: en
command: system
subcommand: check
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system check

## Description 

{% assign name = site.data.mesheryctlcommands.cmds[page.command].subcommands[page.subcommand] %}
{{ name.description }}


<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ name.usage }} 
  </div>
</pre> 

<!-- All possible example use cases of the command -->
## Examples

{{ name.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.example }}
  </div>
</pre>
{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.example }}
  </div>
</pre>
{% endfor %}
<br/>


<!-- Options/Flags available in this command -->
## Options

{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.name }}
  </div>
</pre>
{% endfor %}
<br/>
