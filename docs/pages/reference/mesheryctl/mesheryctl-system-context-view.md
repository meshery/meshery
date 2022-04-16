---
layout: default
title: mesheryctl-system-context-view
permalink: /reference/mesheryctl/system/context/view/
redirect_from: /reference/mesheryctl/system/context/view/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system context view

view current context

## Synopsis

Display active Meshery context

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context view [context-name | --context context-name| --all] --flags [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	View default context
	mesheryctl system context view

	View specified context
	mesheryctl system context view context-name

	View specified context with context flag
	mesheryctl system context view --context context-name

	View config of all contexts
	mesheryctl system context view --all
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --all    Show configs for all of the context
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

* [mesheryctl system context](context/)	 - Configure your Meshery deployment(s)

