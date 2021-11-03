---
layout: default
title: mesheryctl app
permalink: /v0.3/reference/mesheryctl/apps
type: reference
display-title: "false"
language: en
command: app
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->

# mesheryctl app

<!-- Description of the command. Preferably a paragraph -->

## Description

{% assign name = site.data.mesheryctlcommands.cmds[page.command] %}
{{ name.description }}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl app [flags] 
  </div>
</pre>

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
<br/>

<!-- Options/Flags available in this command --> 
## Options & Flags

{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.name }}
  </div>
</pre>
{% endfor %}
<br/>

## Options inherited from parent commands

<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>