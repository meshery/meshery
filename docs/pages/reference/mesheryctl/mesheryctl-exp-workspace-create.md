---
layout: default
title: mesheryctl-exp-workspace-create
permalink: reference/mesheryctl/exp/workspace/create
redirect_from: reference/mesheryctl/exp/workspace/create/
type: reference
display-title: "false"
language: en
command: exp
subcommand: workspace
---

# mesheryctl exp workspace create

Create a new workspace

## Synopsis

Create a new workspace by providing the name, description, and organization ID.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace create [flags]

</div>
</pre> 

## Examples

Create a new workspace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace create --orgId [orgId] --name [name] --description [description]

</div>
</pre> 

Documentation for workspace can be found at:
<pre class='codeblock-pre'>
<div class='codeblock'>
https://docs.layer5.io/cloud/spaces/workspaces/

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -d, --description string   Description of the workspace
  -h, --help                 help for create
  -n, --name string          Name of the workspace
  -o, --orgId string         Organization ID

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
