---
layout: default
title: mesheryctl 
permalink: reference/mesheryctl/commands/mesheryctl
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
# mesheryctl 

## Description
Global command
<br/>

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl [flags]
    </div>
  </div>
</pre>


<!-- Options/Flags available in this command -->
## Options

{% for flag_hash in site.data.mesheryctlcommands.global.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.name }}  # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

## Examples

<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl version
    </div>
  </div>
</pre>
<br/>

{% for flag_hash in site.data.mesheryctlcommands.global.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.usage }}  
    </div>
  </div>
</pre>
<br/>
{% endfor %}