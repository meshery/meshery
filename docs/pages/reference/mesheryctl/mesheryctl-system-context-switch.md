---
layout: default
title: mesheryctl-system-context-switch
permalink: /reference/mesheryctl/system/context/switch/
redirect_from: /reference/mesheryctl/system/context/switch/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system context switch

switch context

## Synopsis

Configure mesheryctl to actively use one one context vs. another context

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context switch context-name [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	// Switch to context named "context-name"
	mesheryctl system context switch context-name
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for switch

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

* [mesheryctl system context](context/)	 - Configure your Meshery deployment(s)

