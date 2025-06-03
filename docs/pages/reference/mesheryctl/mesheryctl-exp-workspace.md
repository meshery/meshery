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

Managge workspaces under an organization

## Synopsis

Create, list of workspaces under an organization
Documentation for models can be found at https://docs.meshery.io/reference/mesheryctl/exp/workspace
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

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count   total number of registered workspaces
  -h, --help    help for workspace

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

* [mesheryctl exp workspace create](/reference/mesheryctl/exp/workspace/create)
* [mesheryctl exp workspace list](/reference/mesheryctl/exp/workspace/list)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
