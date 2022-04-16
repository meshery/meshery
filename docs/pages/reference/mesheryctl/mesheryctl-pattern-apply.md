---
layout: default
title: mesheryctl-pattern-apply
permalink: /reference/mesheryctl/pattern/apply/
redirect_from: /reference/mesheryctl/pattern/apply/
type: reference
display-title: "false"
language: en
command: pattern
---

# mesheryctl pattern apply

Apply pattern file

## Synopsis

Apply pattern file will trigger deploy of the pattern file

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern apply [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

	// apply a pattern file
	mesheryctl pattern apply -f [file | URL]

	// deploy a saved pattern
	mesheryctl pattern apply [pattern-name]
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to pattern file
  -h, --help          help for apply
      --skip-save     Skip saving a pattern

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

