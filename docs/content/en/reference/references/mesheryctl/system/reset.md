---
title: mesheryctl-system-reset
display_title: false
command: system
subcommand: reset
---

# mesheryctl system reset

Reset Meshery's configuration

## Synopsis

Reset Meshery to it's default configuration.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system reset [flags]

</div>
</pre> 

## Examples

Resets meshery.yaml file with a copy from Meshery repo
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system reset

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for reset

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## Screenshots

Usage of mesheryctl system reset
![reset-usage](../../../images/reset.png)

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
