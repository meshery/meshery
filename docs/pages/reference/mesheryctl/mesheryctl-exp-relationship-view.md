---
layout: default
title: mesheryctl-exp-relationship-view
permalink: reference/mesheryctl/exp/relationship/view
redirect_from: reference/mesheryctl/exp/relationship/view/
type: reference
display-title: "false"
language: en
command: exp
subcommand: relationship
---

# mesheryctl exp relationship view

view relationships of a model by its name

## Synopsis

view a relationship queried by the model name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship view [flags]

</div>
</pre> 

## Examples

View relationships of a model in default format yaml
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship view [model-name]

</div>
</pre> 

View relationships of a model in JSON format
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship view [model-name] --output-format json

</div>
</pre> 

View relationships of a model in json format and save it to a file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp relationship view [model-name] --output-format json --save

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")
  -s, --save                   (optional) save output as a JSON/YAML file

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
