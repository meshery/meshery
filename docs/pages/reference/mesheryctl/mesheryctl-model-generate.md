---
layout: default
title: mesheryctl-model-generate
permalink: reference/mesheryctl/model/generate
redirect_from: reference/mesheryctl/model/generate/
type: reference
display-title: "false"
language: en
command: model
subcommand: generate
---

# mesheryctl model generate

Generate models from a file

## Synopsis

Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name
Documentation for models generate can be found at https://docs.meshery.io/reference/mesheryctl/model/generate
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate [flags]

</div>
</pre> 

## Examples

Generate a model from a CSV file(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate --f [path-to-csv-drectory]

</div>
</pre> 

Generate a model from a Uri baesd on a JSON template
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate --f [URL] -t [path-to-template.json]

</div>
</pre> 

Generate a model from a Uri baesd on a JSON template skipping registration
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate --f [URL] -t [path-to-template.json] -r

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string       Specify path to the file or directory
  -h, --help              help for generate
  -r, --register          Skip registration of the model
  -t, --template string   Specify path to the template JSON file

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
