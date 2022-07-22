---
layout: default
title: mesheryctl
permalink: reference/mesheryctl/main
redirect_from: reference/mesheryctl/main/
type: reference
display-title: "false"
language: en
command: mesheryctl
subcommand: nil
---

# mesheryctl

Meshery Command Line tool

## Synopsis

Meshery is the service mesh management plane, providing lifecycle, performance, and configuration management of service meshes and their workloads.

## Examples

Base command
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl

</div>
</pre> 

Display help about command/subcommand
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl --help

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system start --help

</div>
</pre> 

For viewing verbose output
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl -v [or] --verbose

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -h, --help            help for mesheryctl
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
