---
layout: default
title: mesheryctl-system-context
permalink: /reference/mesheryctl/system/context/
redirect_from: /reference/mesheryctl/system/context/
type: reference
display-title: "false"
language: en
command: system
---

# mesheryctl system context

Configure your Meshery deployment(s)

## Synopsis

Configure and switch between different named Meshery server and component versions and deployments.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context [command] [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// Base command
	mesheryctl system context
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for context

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
* [mesheryctl system context create](context/create/)	 - Create a new context (a named Meshery deployment)
* [mesheryctl system context delete](context/delete/)	 - delete context
* [mesheryctl system context list](context/list/)	 - list contexts
* [mesheryctl system context switch](context/switch/)	 - switch context
* [mesheryctl system context view](context/view/)	 - view current context

