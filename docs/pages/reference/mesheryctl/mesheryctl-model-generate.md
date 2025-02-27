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

Generate models from mesheryctl command

## Synopsis

Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model generate [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl model generate -f [ URI ]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl model generate -f [ URI ] -t [ path to template file ] ( only required in case of URL )

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl model generate -f [ URI ] -t [ path to template file ] -r ( to skip registration by default registration is true)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl model generate --f /path/to/csv-drectory

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
        mesheryctl model generate --f http://example.com/model -t /path/to/template.json 

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesheryctl model generate --f http://example.com/model -t /path/to/template.json -r

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
