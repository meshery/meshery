---
layout: default
title: mesheryctl perf apply
permalink: reference/mesheryctl/perf/apply
type: reference
display-title: "false"
language: en
command: perf
subcommand: apply
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl perf apply

## Description

{% assign name = site.data.mesheryctlcommands.cmds[page.command].subcommands[page.subcommand] %}
{{ name.description }}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl perf apply [flags]
  </div>
</pre>

## Examples

{{ name.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ name.usage }}
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