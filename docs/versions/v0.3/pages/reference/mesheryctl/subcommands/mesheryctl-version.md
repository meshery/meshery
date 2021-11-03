---
layout: default
title: mesheryctl version
permalink: reference/mesheryctl/mesheryctl/version
type: reference
display-title: "false"
language: en
command: global
subcommand: version
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
# mesheryctl version

<!-- Description of the command. Preferably a paragraph -->
## Description 

{% assign name = site.data.mesheryctlcommands.cmds[page.command] %}
{% for subcommand_hash in name.subcommands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}


<!-- Basic usage of the command -->
{% for subcommand_hash in name.subcommands %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ subcommand.usage }}
  </div>
</pre> 
{% endfor %}
