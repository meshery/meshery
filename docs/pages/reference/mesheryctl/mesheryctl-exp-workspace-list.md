---
layout: default
title: mesheryctl-exp-workspace-list
permalink: reference/mesheryctl/exp/workspace/list
redirect_from: reference/mesheryctl/exp/workspace/list/
type: reference
display-title: "false"
language: en
command: exp
subcommand: workspace
---

# mesheryctl exp workspace list

List registered workspaces

## Synopsis

List name of all registered workspaces
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list [flags]

</div>
</pre> 

## Examples

List all registered workspace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list --orgId [orgId]

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
  -h, --help           help for list
  -o, --orgId string   Organization ID

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
