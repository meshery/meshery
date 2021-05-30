---
layout: default
title: mesheryctl version | Meshery
permalink: reference/mesheryctl/mesheryctl/version
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
# mesheryctl version

<!-- Description of the command. Preferably a paragraph -->
## Description 

{% for subcommand_hash in site.data.mesheryctlcommands.global.subcommands %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}


<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl version
  </div>
</pre> 
