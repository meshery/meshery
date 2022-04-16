---
layout: default
title: mesheryctl-system-channel-view
permalink: /reference/mesheryctl/system/channel/view/
redirect_from: /reference/mesheryctl/system/channel/view/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system channel view

view release channel and version

## Synopsis

View release channel and version of context in focus

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel view [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// View current release channel
	mesheryctl system channel view
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --all    Show release channel for all contexts
  -h, --help   help for view

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

* [mesheryctl system channel](channel/)	 - Switch between release channels

