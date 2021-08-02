---
layout: default
title: meshery system config
permalink: reference/mesheryctl/system/config
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system config

## Description

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.configure.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system config &#60;managed kubernetes service name&#62; [flags]
  </div>
</pre>

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.configure.command %}{% assign subcommand =subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.configure.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>

{% for ex_hash in site.data.mesheryctlcommands.lifecycle.system.configure.example %}{% assign ex = ex_hash[1] %}
{{ ex.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ ex.example }}
  </div>
</pre>
{% endfor %}
<br/>

<!-- Options/Flags available in this command -->
## Options

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.configure.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.flag }}
  </div>
</pre>
{% endfor %}
<br/>
