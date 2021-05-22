---
layout: default
title: apply
permalink: reference/mesheryctl/pattern/apply
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
# mesheryctl pattern apply

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for subcommand_hash in site.data.mesheryctlcommands.pattern.apply.commands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl pattern apply [flags] 
  </div>
</pre>

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.pattern.apply.commands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.pattern.apply.flags %}{% assign flag = flag_hash[1] %}
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

<pre class="codeblock-pre">
  <div class="codeblock">
  {% for flag_hash in site.data.mesheryctlcommands.pattern.apply.flags %}{% assign flag = flag_hash[1] %}
  {{ flag.name }} # {{ flag.description }}
  {% endfor %}
  </div>
</pre>
<br/>
