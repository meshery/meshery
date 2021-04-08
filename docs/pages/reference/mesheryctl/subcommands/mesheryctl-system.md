---
layout: default
title: mesheryctl system
permalink: reference/mesheryctl/commands/subcommands/mesheryctl-system
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
# mesheryctl system

<!-- Description of the command. Preferably a paragraph -->
## Description

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.commands %}{% assign command = command_hash[1] %}
{{ command.description }}
{% endfor %}
<br/>

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  mesheryctl system [flags] 
    </div>
  </div>
</pre>


<!-- Options/Flags available in this command -->
## Options & Flags

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.flag }} # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.start.command %}{% assign command = command_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ command.name }} # {{ command.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.start.flag %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.flag }} # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.stop.command %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ subcommand.name }} # {{ subcommand.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.stop.flag %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.flag }} # {{ flag.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.subcommands %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ subcommand.name }} # {{ subcommand.description }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

<!-- All possible example use cases of the command -->
## Examples

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.flags %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for command_hash in site.data.mesheryctlcommands.lifecycle.system.start.command %}{% assign command = command_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ command.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.start.flag %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.stop.command %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ subcommand.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.stop.flag %}{% assign flag = flag_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ flag.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.subcommands %}{% assign subcommand = subcommand_hash[1] %}
<pre class="codeblock-pre">
  <div class="codeblock">
    <div class="clipboardjs">
  {{ subcommand.usage }}
    </div>
  </div>
</pre>
<br/>
{% endfor %}