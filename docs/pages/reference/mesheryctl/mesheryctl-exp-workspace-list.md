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
Documentation for models can be found at https://docs.meshery.io/reference/mesheryctl/exp/workspace/list
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list [flags]

</div>
</pre> 

## Examples

List of workspace under a specific organization
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list --orgId [orgId]

</div>
</pre> 

List of workspace under a specific organization for a specified page
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list --orgId [orgId] --page [page-number]

</div>
</pre> 

Display number of available  workspace under a specific organization
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl exp workspace list --orgId [orgId] --count

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count          total number of registered workspaces
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
