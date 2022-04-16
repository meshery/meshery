---
layout: default
title: mesheryctl-pattern-delete
permalink: /reference/mesheryctl/pattern/delete/
redirect_from: /reference/mesheryctl/pattern/delete/
type: reference
display-title: "false"
language: en
command: pattern
---

# mesheryctl pattern delete

Delete pattern file

## Synopsis

delete pattern file will trigger deletion of the pattern file

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern delete [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>


	// delete a pattern file
	mesheryctl pattern delete [pattern-name]
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to pattern file
  -h, --help          help for delete

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

