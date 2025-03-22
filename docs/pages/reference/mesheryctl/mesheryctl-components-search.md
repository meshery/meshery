---
layout: default
title: mesheryctl-components-search
permalink: reference/mesheryctl/components/search
redirect_from: reference/mesheryctl/components/search/
type: reference
display-title: "false"
language: en
command: components
subcommand: search
---

# mesheryctl components search

search registered components

## Synopsis

search components registered in Meshery Server based on kind
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl components search [flags]

</div>
</pre> 

## Examples

Search for components using a query
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl components search [query-text]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for search

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
