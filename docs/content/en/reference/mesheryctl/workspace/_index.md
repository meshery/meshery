---
title: mesheryctl-workspace
display_title: false
command: workspace
subcommand: nil
---

# mesheryctl workspace

Manage workspaces under an organization

## Synopsis

Create, list of workspaces under an organization

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl workspace [flags]

</div>
</pre> 

## Examples

To view a list workspaces
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl workspace list --orgId [orgId]

</div>
</pre> 

To create a workspace
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl workspace create --orgId [orgId] --name [name] --description [description]

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
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
