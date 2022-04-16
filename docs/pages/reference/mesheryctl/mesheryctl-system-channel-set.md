---
layout: default
title: mesheryctl-system-channel-set
permalink: /reference/mesheryctl/system/channel/set/
redirect_from: /reference/mesheryctl/system/channel/set/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system channel set

set release channel and version

## Synopsis

Set release channel and version of context in focus

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel set [stable|stable-version|edge|edge-version] [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// Subscribe to release channel or version
	mesheryctl system channel set [stable|stable-version|edge|edge-version]
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for set

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

