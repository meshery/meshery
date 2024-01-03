---
layout: default
title: mesheryctl-system-channel-view
permalink: reference/mesheryctl/system/channel/view
redirect_from: reference/mesheryctl/system/channel/view/
type: reference
display-title: "false"
language: en
command: system
subcommand: channel
---

# mesheryctl system channel view

view release channel and version

## Synopsis

View release channel and version of context in focus
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel view [flags]

</div>
</pre> 

## Examples

View current release channel
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system channel view

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -a, --all    Show release channel for all contexts
  -h, --help   help for view

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

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
