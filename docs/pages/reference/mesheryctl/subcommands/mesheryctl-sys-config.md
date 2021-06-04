---
layout: default
title: meshery system config | Meshery
permalink: reference/mesheryctl/system/config
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
# mesheryctl system config

## Description

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.configure.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system config < managed kubernetes service name > [flags]
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

To configure Meshery to use Google Kubernetes Engine-
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system config gke --token [path to token]
  </div>
</pre>

<br/>

To configure Meshery to use Azure Kubernetes Service-
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system config aks --token [path to token]
  </div>
</pre>

<br/>

To configure Meshery to use Elastic Kubernetes Service-
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system config eks --token <path-to-token>
  </div>
</pre> 

<br/>


<!-- Options/Flags available in this command -->
## Options

<pre class="codeblock-pre">
  <div class="codeblock">
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.configure.flag %}{% assign flag = flag_hash[1] %}
    {{ flag.flag }} # {{ flag.description }}
    {% endfor %}
  </div>
</pre>
<br/>
