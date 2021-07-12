---
layout: default
title: mesheryctl
permalink: reference/mesheryctl/mesheryctl
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl 

## Description
Global command
<br/>

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl [flags]
  </div>
</pre>

## Examples

{% for flag_hash in site.data.mesheryctlcommands.global.flags %}{% assign flag = flag_hash[1] %}
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

{% for flag_hash in site.data.mesheryctlcommands.global.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.name }}
  </div>
</pre>
{% endfor %}
<br/>
