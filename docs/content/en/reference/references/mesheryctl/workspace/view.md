---
title: mesheryctl-workspace-view
display_title: false
command: workspace
subcommand: view
---

# mesheryctl workspace view

View a workspace

## Synopsis

View a workspace by its ID or name.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace view [workspace-name|workspace-id] [flags]

</div>
</div>
</pre> 

## Examples

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

View details of a specific workspace in JSON format
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace view [workspace-id] --orgId [orgId] --output-format json

</div>
</div>
</pre> 

View details of a specific workspace and save it to a file
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl workspace view [workspace-id] --orgId [orgId] --output-format json --save

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                   help for view
      --orgId string           (required) organization ID
  -o, --output-format string   (optional) format to display in [json|yaml] (default "yaml")
  -s, --save                   (optional) save output as a JSON/YAML file

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
