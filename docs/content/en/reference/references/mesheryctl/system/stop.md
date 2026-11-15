---
title: mesheryctl-system-stop
display_title: false
command: system
subcommand: stop
categories: [mesheryctl-system]
---

# mesheryctl system stop

Stop Meshery

## Synopsis

Stop all Meshery containers / remove all Meshery resources.
	
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system stop [flags]

</div>
</div>
</pre> 

## Examples

Stop Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system stop

</div>
</div>
</pre> 

Reset Meshery's configuration file to default settings.
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system stop --reset

</div>
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

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
