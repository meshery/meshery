---
layout: default
title: mesheryctl app view
permalink: reference/mesheryctl/apps/view
type: reference
display-title: "false"
language: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl app view [application name/id] [flags]

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for subcommand_hash in site.data.mesheryctlcommands.apps.view.commands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl app view [app-name|app-id]
  mesheryctl app view [app-name|app-id] -o json
  mesheryctl app view --all [View all applications] 
  </div>
</pre>

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.apps.view.commands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.apps.view.flags %}{% assign flag = flag_hash[1] %}
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

{% for flag_hash in site.data.mesheryctlcommands.apps.view.flags %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.name }}
  </div>
</pre>
{% endfor %}
<br/>
