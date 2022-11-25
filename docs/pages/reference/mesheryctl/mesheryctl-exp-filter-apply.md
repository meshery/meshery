---
layout: default
title: mesheryctl-exp-filter-apply
permalink: reference/mesheryctl/exp/filter/apply
redirect_from: reference/mesheryctl/exp/filter/apply/
type: reference
display-title: "false"
language: en
command: exp
subcommand: filter
---

# mesheryctl exp filter apply

Apply filter file

## Synopsis

Apply filter file will trigger deploy of the filter file

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter apply [flags]

</div>
</pre> 

## Examples

Apply WASM filter file (login required)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter apply --file [GitHub Link]

</div>
</pre> 

Apply a remote filter file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter apply --file https://github.com/layer5io/wasm-filters/tree/master/http-auth

</div>
</pre> 

Apply a filter file using file name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp filter apply [File Name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Path to filter file
  -h, --help          help for apply
      --skip-save     Skip saving a filter

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
