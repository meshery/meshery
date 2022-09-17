---
layout: default
title: mesheryctl-pattern-delete
permalink: reference/mesheryctl/pattern/delete
redirect_from: reference/mesheryctl/pattern/delete/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: delete
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

delete a pattern file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern delete [file | URL]

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
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
