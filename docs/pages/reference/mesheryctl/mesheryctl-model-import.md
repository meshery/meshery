---
layout: default
title: mesheryctl-model-import
permalink: reference/mesheryctl/model/import
redirect_from: reference/mesheryctl/model/import/
type: reference
display-title: "false"
language: en
command: model
subcommand: import
---

# mesheryctl model import

Import models from mesheryctl command

## Synopsis

Import models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import [flags]

</div>
</pre> 

## Examples

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesehryctl model import -f [ URI ]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesehryctl model import -f URL 

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesehryctl model import -f OCI 

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesehryctl model import -f model.tar.gz 

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
	mesehryctl model import -f /path/to/models

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -f, --file string   Specify path to the file or directory
  -h, --help          help for import

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
