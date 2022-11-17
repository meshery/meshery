---
layout: default
title: mesheryctl-exp-filter-view
permalink: reference/mesheryctl/exp/filter/view
redirect_from: reference/mesheryctl/exp/filter/view/
type: reference
display-title: "false"
language: en
command: exp
subcommand: filter
---

# mesheryctl exp filter view

Display filters(s)

## Synopsis

Displays the contents of a specific filter based on name or id

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter view [filter name] [flags]

</div>
</pre> 

## Examples

View the specified WASM filter file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter view [filter-name | ID]	

</div>
</pre> 

View using filter name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter view test-wasm

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --all                    (optional) view all filters available
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")

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
