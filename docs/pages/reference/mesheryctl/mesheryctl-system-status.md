---
layout: default
title: mesheryctl-system-status
permalink: /reference/mesheryctl/system/status/
redirect_from: /reference/mesheryctl/system/status/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system status

Check Meshery status

## Synopsis

Check status of Meshery and Meshery components.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system status [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	// Check status of Meshery, Meshery adapters, Meshery Operator and its controllers.
	mesheryctl system status 
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for status

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

* [mesheryctl system](system/)	 - Meshery Lifecycle Management

