---
layout: default
title: mesheryctl-component-search
permalink: reference/mesheryctl/component/search
redirect_from: reference/mesheryctl/component/search/
type: reference
display-title: "false"
language: en
command: component
subcommand: search
---

# mesheryctl component search

Search registered components

## Synopsis

Search components registered in Meshery Server based on kind

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component search [flags]

</div>
</pre> 

## Examples

Search for components using a query
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component search [query-text]

</div>
</pre> 

Search for multi-word component names (must be quoted)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component search "Component name"

</div>
</pre> 

Search list of components of specified page [int]
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl component search [query-text] [--page 1]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for search
  -p, --page int       (optional) List next set of components with --page (default = 1) (default 1)
  -s, --pagesize int   (optional) List next set of components with --pagesize (default = 10) (default 10)

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
