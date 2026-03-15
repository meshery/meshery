---
layout: default
title: mesheryctl-system-context-switch
permalink: reference/mesheryctl/system/context/switch
redirect_from: reference/mesheryctl/system/context/switch/
type: reference
display-title: "false"
language: en
command: system
subcommand: context
---

# mesheryctl system context switch

switch context

## Synopsis

Configure mesheryctl to actively use one one context vs. another context
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context switch context-name [flags]

</div>
</pre> 

## Examples

Switch to context named "sample"
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system context switch sample

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for switch

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

Usage of mesheryctl context switch
![context-switch-usage](/assets/img/mesheryctl/contextswitch.png)

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
