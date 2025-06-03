---
layout: default
title: mesheryctl-exp-model-init
permalink: reference/mesheryctl/exp/model/init
redirect_from: reference/mesheryctl/exp/model/init/
type: reference
display-title: "false"
language: en
command: exp
subcommand: model
---

# mesheryctl exp model init

Generates scaffolding for convenient model creation

## Synopsis

Generates a folder structure and guides user on model creation
Documentation for exp models init can be found at https://docs.meshery.io/reference/mesheryctl/exp/model/init
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model init [flags]

</div>
</pre> 

## Examples

generates a folder structure
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model init [model-name]

</div>
</pre> 

generates a folder structure and sets up model version
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model init [model-name] --version [version] (default is v0.1.0)

</div>
</pre> 

generates a folder structure under specified path
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model init [model-name] --path [path-to-location] (default is current folder)

</div>
</pre> 

generate a folder structure in json format
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp model init [model-name] --output-format [json|yaml|csv] (default is json)

</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
    

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for init
  -o, --output-format string   (optional) format to display in [json|yaml] (default "json")
  -p, --path string            (optional) target directory (default: current dir) (default ".")
      --version string         (optional) model version (default: v0.1.0) (default "v0.1.0")

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
