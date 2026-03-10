---
layout: default
title: mesheryctl-model-delete
permalink: reference/mesheryctl/model/delete
redirect_from: reference/mesheryctl/model/delete/
type: reference
display-title: "false"
language: en
command: model
subcommand: delete
---

# mesheryctl model delete

Delete a model

## Synopsis

Delete a model by ID or Name

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model delete [model-id | model-name] [flags]

</div>
</pre> 

## Examples

Delete a model by ID
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model delete [model-id]

</div>
</pre> 

Delete a model by name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model delete [model-name]

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
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
