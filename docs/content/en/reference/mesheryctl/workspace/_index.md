---
title: mesheryctl-workspace
display_title: false
command: workspace
subcommand: nil
categories: [mesheryctl-workspace]
aliases:
- /reference/references/mesheryctl/workspace/
---

# mesheryctl workspace

Manage workspaces under an organization

## Synopsis

Create, list, and view workspaces under an organization

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace [flags]

</div>
</div>
</pre> 

## Examples

List workspaces under a specific organization
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace list --orgId [orgId]

</div>
</div>
</pre> 

Create a workspace
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace create --orgId [orgId] --name [name] --description [description]

</div>
</div>
</pre> 

View details of a specific workspace by ID
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace view [workspace-id] --orgId [orgId]

</div>
</div>
</pre> 

View details of a specific workspace by name
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace view [workspace-name] --orgId [orgId]

</div>
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
      --config string   path to config file (default "~/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#documenting" >}}) for guidance.
