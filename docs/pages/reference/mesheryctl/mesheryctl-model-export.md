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

Export registered models

## Synopsis

Export the registered model to the specified output type
Documentation for models export can be found at https://docs.meshery.io/reference/mesheryctl/model/export
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

Export a model by name in JSON type
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] -t [yaml|json] (default is YAML)

</div>
</pre> 

Export a model by name in YAML type in a specific location
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] -l [path-to-location]

</div>
</pre> 

Export a model by name in YAML type discarding components and relationships
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] --discard-components --discard-relationships

</div>
</pre> 

Export a model version by name in YAML type
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl model export [model-name] --version [version (ex: v0.7.3)]

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
  -p, --page int                 (optional) List next set of models with --page (default = 1) (default 1)
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
