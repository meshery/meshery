---
title: mesheryctl-system-channel
display_title: false
command: system
subcommand: channel
---

# mesheryctl system channel

Switch between release channels

## Synopsis

Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system channel [flags]

</div>
</div>
</pre> 

## Examples

Subscribe to release channel or version
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system channel

</div>
</div>
</pre> 

To set the channel
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system channel set [stable|stable-version|edge|edge-version]

</div>
</div>
</pre> 

To pin/set the channel to a specific version
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system channel set stable-v0.6.0

</div>
</div>
</pre> 

To view release channel and version
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system channel view

</div>
</div>
</pre> 

To switch release channel and version
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system channel switch [stable|stable-version|edge|edge-version]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for channel

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
