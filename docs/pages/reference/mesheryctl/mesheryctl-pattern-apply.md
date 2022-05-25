---
layout: default
title: mesheryctl-pattern-apply
permalink: reference/mesheryctl/pattern/apply
redirect_from: reference/mesheryctl/pattern/apply/
type: reference
display-title: "false"
language: en
command: pattern
subcommand: apply
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

apply a pattern file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl pattern apply -f [file | URL]

</div>
</pre> 

deploy a saved pattern
<pre class='codeblock-pre'>
<div class='codeblock'>
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
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
