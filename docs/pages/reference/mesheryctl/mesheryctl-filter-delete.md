---
layout: default
title: mesheryctl-filter-delete
permalink: reference/mesheryctl/filter/delete
redirect_from: reference/mesheryctl/filter/delete/
type: reference
display-title: "false"
language: en
command: filter
subcommand: delete
---

# mesheryctl filter delete

Delete a filter file

## Synopsis

Delete a filter file using the name or ID of a filter

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter delete [filter-name | ID] [flags]

</div>
</pre> 

## Examples

Delete the specified WASM filter file using name or ID
A unique prefix of the name or ID can also be provided. If the prefix is not unique, the first match will be deleted.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter delete [filter-name | ID]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for delete

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
