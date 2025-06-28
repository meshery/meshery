---
layout: default
title: mesheryctl-exp-model-build
permalink: reference/mesheryctl/exp/model/build
redirect_from: reference/mesheryctl/exp/model/build/
type: reference
display-title: "false"
language: en
command: exp
subcommand: model
---

# mesheryctl exp model build

Create an OCI-compliant package from the model files

## Synopsis

Create an OCI-compliant package from the model files.
Model files are taken from [path]/[model-name]/[model-version] folder.
Expects input to be in the format scaffolded by the model init command.
Documentation for exp model and subcommands can be found at https://docs.meshery.io/reference/mesheryctl/exp/model
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model build [flags]

</div>
</pre> 

## Examples

Create an OCI-compliant package from the model files
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model build [model-name]/[model-version]

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help          help for build
  -p, --path string   (optional) target directory to get model from (default: current dir) (default ".")

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
