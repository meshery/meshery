---
layout: default
title: mesheryctl-filter-import
permalink: reference/mesheryctl/filter/import
redirect_from: reference/mesheryctl/filter/import/
type: reference
display-title: "false"
language: en
command: filter
subcommand: import
---

# mesheryctl filter import

Import a WASM filter

## Synopsis

Import a WASM filter from a URI (http/s) or local filesystem path

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter import [URI] [flags]

</div>
</pre> 

## Examples

Import a filter file from local filesystem
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter import /path/to/filter.wasm

</div>
</pre> 

Import a filter file from a remote URI
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter import https://example.com/myfilter.wasm

</div>
</pre> 

Add WASM configuration 
If the string is a valid file in the filesystem, the file is read and passed as a string. Otherwise, the string is passed as is.
Use quotes if the string contains spaces
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter import /path/to/filter.wasm -config [filepath|string]

</div>
</pre> 

Specify the name of the filter to be imported. Use quotes if the name contains spaces
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter import /path/to/filter.wasm -name [string]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --config string   (optional) WASM configuration filepath/string
  -h, --help            help for import
  -n, --name string     (optional) filter name

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
  -t, --token string   Path to token file default from current context
  -v, --verbose        verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/) 
