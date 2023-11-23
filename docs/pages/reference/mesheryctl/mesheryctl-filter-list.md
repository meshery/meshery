---
layout: default
title: mesheryctl-filter-list
permalink: reference/mesheryctl/filter/list
redirect_from: reference/mesheryctl/filter/list/
type: reference
display-title: "false"
language: en
command: filter
subcommand: list
---

# mesheryctl filter list

List filters

## Synopsis

Display list of all available filter files.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter list [flags]

</div>
</pre> 

## Examples

List all WASM filter files present
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter list	(maximum 25 filters)

</div>
</pre> 

Search for filter
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter list Test (maximum 25 filters)

</div>
</pre> 

Search for filter with space
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl filter list 'Test Filter' (maximum 25 filters)

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help       help for list
  -p, --page int   (optional) List next set of filters with --page (default = 1) (default 1)
  -v, --verbose    Display full length user and filter file identifiers

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -t, --token string    Path to token file default from current context

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
