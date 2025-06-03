---
layout: default
title: mesheryctl-model-view
permalink: reference/mesheryctl/model/view
redirect_from: reference/mesheryctl/model/view/
type: reference
display-title: "false"
language: en
command: model
subcommand: view
---

# mesheryctl model view

View model

## Synopsis

View a model queried by its name
Documentation for models view can be found at https://docs.meshery.io/reference/mesheryctl/model/view
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model view [flags]

</div>
</pre> 

## Examples

View a specific model from current provider
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model view [model-name]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/n2/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
