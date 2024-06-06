---
layout: default
title: mesheryctl-exp-workspace
permalink: reference/mesheryctl/exp/workspace
redirect_from: reference/mesheryctl/exp/workspace/
type: reference
display-title: "false"
language: en
command: exp
subcommand: workspace
---

# mesheryctl exp workspace

View list of workspaces and detail of workspaces

## Synopsis

View list of workspaces and detailed information of a specific workspaces
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace [flags]

</div>
</pre> 

## Examples

To view a list workspaces
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list --orgId [orgId]

</div>
</pre> 

To create a workspace
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
  -h, --help   help for workspace

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
