---
layout: default
title: meshery system config
permalink: reference/mesheryctl/system/config
redirect_from: reference/mesheryctl/system/config/
type: reference
display-title: "false"
language: en
command: system
subcommand: config
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl system config

## Description

{% assign name = site.data.mesheryctlcommands.cmds[page.command].subcommands[page.subcommand] %}
{{ name.description }}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system config &#60;managed kubernetes service name&#62; [flags]
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
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}
<br/>
{% for flag_hash in name.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.example }}
  </div>
</pre>
{% endfor %}
<br />

{% for ex_hash in name.examples %}{% assign ex = ex_hash[1] %}
{{ ex.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ ex.usage }}
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
