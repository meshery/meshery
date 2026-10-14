---
title: mesheryctl-workspace-update
display_title: false
command: workspace
subcommand: update
categories: [mesheryctl-workspace]
---

# mesheryctl workspace update

Update a workspace

## Synopsis

Update a workspace's name and/or description by its ID.
At least one of --name or --description must be provided; neither can be set
to an empty value (clearing a description is not currently supported by the
server).
--orgId must match the workspace's current organization. The CLI checks this
before updating the workspace; it does not move the workspace to a different
organization or grant update permission.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace update [workspace-id] [flags]

</div>
</div>
</pre> 

## Examples

Rename a workspace
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace update [workspace-id] --orgId [orgId] --name [new-name]

</div>
</div>
</pre> 

Update a workspace's description
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace update [workspace-id] --orgId [orgId] --description [new-description]

</div>
</div>
</pre> 

Update both
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace update [workspace-id] --orgId [orgId] --name [new-name] --description [new-description]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -d, --description string   New description for the workspace
  -h, --help                 help for update
  -n, --name string          New name for the workspace
      --orgId string         (required) organization ID - must match the workspace's current organization

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

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
