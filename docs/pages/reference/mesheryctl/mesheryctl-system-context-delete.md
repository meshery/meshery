---
layout: default
title: mesheryctl-system-context-delete
permalink: /reference/mesheryctl/system/context/delete/
redirect_from: /reference/mesheryctl/system/context/delete/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system context delete

delete context

## Synopsis

Delete an existing context (a named Meshery deployment) from Meshery config file

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context delete context-name [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// Delete context
	mesheryctl system context delete [context name]
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help         help for delete
  -s, --set string   New context to deploy Meshery

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

