---
layout: default
title: mesheryctl-system-channel
permalink: reference/mesheryctl/system/channel
redirect_from: reference/mesheryctl/system/channel/
type: reference
display-title: "false"
language: en
command: system
subcommand: channel
---

# mesheryctl system channel

Switch between release channels

## Synopsis

Subscribe to a release channel. Choose between either 'stable' or 'edge' channels.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel [flags]

</div>
</pre> 

## Examples

Subscribe to release channel or version
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel 

</div>
</pre> 

To set the channel
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel set [stable|stable-version|edge|edge-version]

</div>
</pre> 

To pin/set the channel to a specific version
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel set stable-v0.6.0

</div>
</pre> 

To view release channel and version
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel view

</div>
</pre> 

To switch release channel and version
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel switch [stable|stable-version|edge|edge-version]

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

Go back to [command reference index](/reference/mesheryctl/) 
