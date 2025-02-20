---
layout: default
title: mesheryctl-model-export
permalink: reference/mesheryctl/model/export
redirect_from: reference/mesheryctl/model/export/
type: reference
display-title: "false"
language: en
command: model
subcommand: export
---

# mesheryctl model export

export registered models

## Synopsis

export the registered model to the specified output type
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [flags]

</div>
</pre> 

## Examples

Export a model by name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] -o [oci|tar]  (default is oci)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] -t json (default is yaml)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] -l /home/meshery/

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] --discard-components --discard-relationships

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] --version v0.7.3

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --discard-components       (optional) whether to discard components in the exported model definition (default = false)
  -r, --discard-relationships    (optional) whether to discard relationships in the exported model definition (default = false)
  -h, --help                     help for export
  -t, --output-format string     (optional) format to display in [json|yaml] (default = yaml) (default "yaml")
  -l, --output-location string   (optional) output location (default = current directory) (default "./")
  -o, --output-type string       (optional) format to display in [oci|tar] (default = oci) (default "oci")
      --version string           (optional) model version to export (default = "")

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
