---
layout: default
title: mesheryctl-model-search
permalink: reference/mesheryctl/model/search
redirect_from: reference/mesheryctl/model/search/
type: reference
display-title: "false"
language: en
command: model
subcommand: search
---

# mesheryctl model search

Search model(s)

## Synopsis

Search model(s) by search string

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model search [flags]

</div>
</pre> 

## Examples

Search model from current provider
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model search [query-text]

</div>
</pre> 

Search list of models for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model search [query-text] --page [page-number]

</div>
</pre> 

Search list of models for a specified pagesize
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model search [query-text] --pagesize [pagesize-number]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for search
  -p, --page int       (optional) List next set of models with --page (default = 1) (default 1)
  -s, --pagesize int   (optional) List next set of models with --pagesize (default = 10) (default 10)

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
