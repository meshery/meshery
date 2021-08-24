---
layout: default
title: mesheryctl app list
permalink: reference/mesheryctl/apps/list
type: reference
display-title: "false"
language: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl app list

<!-- Description of the command. Preferably a paragraph -->
## Description and Example

<!-- Basic usage of the command -->
{% for subcommand_hash in site.data.mesheryctlcommands.apps.list %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
<br/>

<!-- Options/Flags available in this command -->
## Options & Flags

<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl app list -h # Shows help for the command
  </div>
</pre>
<br/>
