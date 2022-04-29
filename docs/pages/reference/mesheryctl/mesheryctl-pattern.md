---
layout: default
title: mesheryctl-pattern
permalink: reference/mesheryctl/pattern
redirect_from: reference/mesheryctl/pattern/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: nil
---

# mesheryctl pattern

Service Mesh Patterns Management

## Synopsis

Manage service meshes using predefined patterns

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>

// Apply pattern file
mesheryctl pattern apply --file [path to pattern file | URL of the file]

// Deprovision pattern file
mesheryctl pattern delete --file [path to pattern file]

// View pattern file
mesheryctl pattern view [pattern name | ID]

// List all patterns
mesheryctl pattern list
	

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for pattern
  -t, --token string   Path to token file default from current context

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/admin-pc/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
