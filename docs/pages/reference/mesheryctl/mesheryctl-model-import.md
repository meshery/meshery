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

Import models

## Synopsis

Import models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name
Documentation for models import can be found at https://docs.meshery.io/reference/mesheryctl/model/import
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import [flags]

</div>
</pre> 

## Examples

Import model
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import -f [URI]

</div>
</pre> 

Import model from a URL to a meshery model
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import -f [URL]

</div>
</pre> 

Import model from an OCI artifact
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import -f [OCI]

</div>
</pre> 

Import model from a tar.gz file
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import -f [path-to-model.tar.gz]

</div>
</pre> 

Import model from a path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import -f [path-to-model]

</div>
</pre> 

Import model using CSV files
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model import -f [path-to-csv-directory]

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
      --config string   path to config file (default "/home/n2/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
