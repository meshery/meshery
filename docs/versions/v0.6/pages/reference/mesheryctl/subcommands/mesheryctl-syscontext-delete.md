---
layout: default
title: mesheryctl system context delete
permalink: reference/mesheryctl/system/context/delete
type: reference
display-title: "false"
language: en
command: system-context
subcommand: delete
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system context delete

## Description 

{% assign name = site.data.mesheryctlcommands.cmds[page.command].subcommands[page.subcommand] %}
{{ name.description }}


<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system context delete [context-name] [flags]
  </div>
</pre> 

<!-- All possible example use cases of the command -->
## Examples

{{ name.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.usage }}
  </div>
</pre>
<br />
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.example }}
  </div>
</pre>
{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>

<!-- Options/Flags available in this command -->
<!-- ## Options & Flags

{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.flag }}
  </div>
</pre>
{% endfor %}
<br/>
-->