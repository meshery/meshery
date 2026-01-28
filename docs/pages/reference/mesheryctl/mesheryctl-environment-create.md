---
layout: default
title: mesheryctl-environment-create
permalink: reference/mesheryctl/environment/create
redirect_from: reference/mesheryctl/environment/create/
type: reference
display-title: "false"
language: en
command: environment
subcommand: create
---

# mesheryctl environment create

Create a new environment

## Synopsis

Create a new environment by providing the name and description of the environment
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/create
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment create [flags]

</div>
</pre> 

## Examples

Create a new environment
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl environment create --orgID [orgID] --name [name] --description [description]

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -d, --description string   Description of the environment
  -h, --help                 help for create
  -n, --name string          Name of the environment
  -o, --orgID string         Organization ID

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
