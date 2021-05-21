---
layout: default
title: view
permalink: reference/mesheryctl/system/context/view
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
# mesheryctl system context view

<!-- Description of the command. Preferably a paragraph -->
## Description 

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system context view [flags] 
  </div>
</pre> 

<!-- All possible example use cases of the command -->
## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.flags %}{% assign flag = flag_hash[1] %}
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
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.flags %}{% assign flag = flag_hash[1] %}
    {{ flag.flag }} # {{ flag.description }}
    {% endfor %}
  </div>
</pre>
<br/>
