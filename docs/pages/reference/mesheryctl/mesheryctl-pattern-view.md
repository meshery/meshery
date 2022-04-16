---
layout: default
title: mesheryctl-pattern-view
permalink: /reference/mesheryctl/pattern/view/
redirect_from: /reference/mesheryctl/pattern/view/
type: reference
display-title: "false"
language: en
command: pattern
---

# mesheryctl pattern view

Display pattern(s)

## Synopsis

Displays the contents of a specific pattern based on name or id

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern view <pattern name> [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// view a pattern
	mesheryctl pattern view [pattern-name/id]
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --all                    (optional) view all patterns available
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

* [mesheryctl pattern](pattern/)	 - Service Mesh Patterns Management

