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
mesheryctl model export <modelname> -o json

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export <modelname> -o yaml

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export <modelname> -l /home/meshery/

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export <modelname> --include-components true --include-relationships true

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                     help for export
  -c, --include-components       whether to include components in the model definition (default = false)
  -r, --include-relationships    whether to include components in the model definition (default = false)
  -l, --output-location string   (optional) output location (default = current directory) (default "./")
  -o, --output-type string       (optional) format to display in [json|yaml] (default = yaml) (default "yaml")

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
