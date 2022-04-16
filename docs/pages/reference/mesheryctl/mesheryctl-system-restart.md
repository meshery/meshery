---
layout: default
title: mesheryctl-system-restart
permalink: /reference/mesheryctl/system/restart/
redirect_from: /reference/mesheryctl/system/restart/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system restart

Stop, then start Meshery

## Synopsis

Restart all Meshery containers / pods.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system restart [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// Restart all Meshery containers, their instances and their connected volumes
	mesheryctl system restart

	// (optional) skip checking for new updates available in Meshery.
	mesheryctl system restart --skip-update
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help          help for restart
      --skip-update   (optional) skip checking for new Meshery's container images.

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

