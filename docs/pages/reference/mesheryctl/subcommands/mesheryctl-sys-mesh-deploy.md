---
layout: default
title: mesheryctl system mesh deploy
permalink: reference/mesheryctl/mesh/deploy
type: reference
display-title: "false"
language: en
command: mesh
subcommand: deploy
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl mesh deploy

<!-- Description of the command. Preferably a paragraph -->
## Description

{% assign name = site.data.mesheryctlcommands.cmds[page.command].subcommands[page.subcommand] %}
{{ name.description }}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl mesh deploy [flags]
  </div>
</pre>

## Examples

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

## Options inherited from parent commands
<pre class="codeblock-pre">
  <div class="codeblock">
  --help, -h # Shows help for the command
  </div>
</pre>
