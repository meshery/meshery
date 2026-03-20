---
title: mesheryctl-model-view
display_title: false
command: model
subcommand: view
---

# mesheryctl model view

View model

## Synopsis

View a model queried by its name or ID

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model view [flags]

</div>
</pre> 

## Examples

View a specific model from current provider by using [model-name] or [model-id] in default format yaml
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model view [model-name]

</div>
</pre> 

View a specific model in specifed format
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model view [model-name] --output-format [json|yaml]

</div>
</pre> 

View a specific model in specified format and save it as a file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model view [model-name] --output-format [json|yaml] --save

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
