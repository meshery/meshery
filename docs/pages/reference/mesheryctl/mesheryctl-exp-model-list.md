---
layout: default
title: mesheryctl-exp-model-list
permalink: reference/mesheryctl/exp/model/list
redirect_from: reference/mesheryctl/exp/model/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: model
---

# mesheryctl exp model list

list models

## Synopsis

list name of all registered models

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model list [flags]

</div>
</pre> 

## Examples

View list of models
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system model list

</div>
</pre> 

View list of models with specified page number (25 models per page)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system model list --page 2

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help       help for list
  -p, --page int   (optional) List next set of models with --page (default = 1) (default 1)

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
