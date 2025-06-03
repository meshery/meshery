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
Documentation for models search can be found at https://docs.meshery.io/reference/mesheryctl/model/search
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
